#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

# Script to convert "telemetry_types.xml" from RAX2 and CSSWE
# into our schema and generate a frame decoder.

from carpcomm.pb import telemetry_pb2
from carpcomm.pb import text_pb2

import xml.sax
import sys
import time


def GetUnitEnum(unit):
    mapping = {
        'volts': (telemetry_pb2.TelemetryDatumSchema.VOLT, (1.0, 0.0)),
        'deg C': (telemetry_pb2.TelemetryDatumSchema.KELVIN, (1.0, 273.15)),
        'milliamps': (telemetry_pb2.TelemetryDatumSchema.AMPERE, (1e-3, 0.0)),
        'days': (telemetry_pb2.TelemetryDatumSchema.SECOND, (86400.0, 0.0)),
        }
    return mapping.get(unit, (None, None))

def GetKey(satellite_id, name, unit):
    unit_suffix = {
        telemetry_pb2.TelemetryDatumSchema.VOLT: '_v',
        telemetry_pb2.TelemetryDatumSchema.KELVIN: '_t',
        telemetry_pb2.TelemetryDatumSchema.AMPERE: '_c',
        telemetry_pb2.TelemetryDatumSchema.SECOND: '_s',
        telemetry_pb2.TelemetryDatumSchema.RADIAN_PER_SECOND: '_r',
        }

    strip_suffixes = [
        'voltage',
        'current',
        'temp']

    n = name.strip().lower()
    for ss in strip_suffixes:
        if n.endswith(ss):
            n = n[:-len(ss)]
    n = n.strip()
    n = n.replace(':', '').replace('.', 'p').replace('/', '_').replace(' ', '_')

    return 'd:%s:%s%s' % (satellite_id, n, unit_suffix.get(unit, ''))

def ConvertParamsToSchema(satellite_id, params):
    s = telemetry_pb2.TelemetryDatumSchema()

    n = s.name.add()
    n.lang = 'en'
    n.text = params['name']

    unit, linear_transform = GetUnitEnum(params['unit'])
    if unit is not None:
        s.unit = unit
        s.type = telemetry_pb2.TelemetryDatumSchema.DOUBLE

    if params['datatype'] == 'bool':
        s.type = telemetry_pb2.TelemetryDatumSchema.BOOL

    if unit is None and params['datatype'] == 'int':
        s.type = telemetry_pb2.TelemetryDatumSchema.INT64

    if unit is None and params['datatype'] == 'lut':
        print 'Error: cannot handle lut'
        return None

    s.key = GetKey(satellite_id, params['name'], unit)

    assert s.HasField('type')
    return s


code_start = """// Generated from %(input_path)s using %(script_name)s on %(date)s.

package telemetry

import "fmt"
import "errors"
import "carpcomm/pb"

func DecodeFrame_%(satellite_id)s(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {

	if len(frame) < %(min_length)d {
		return nil, errors.New(fmt.Sprintf("Frame too short: %%d, expected %%d.", len(frame), %(min_length)d))
	}
"""

code_end = """
	return data, nil
}
"""

def GenerateCodeForDatum(satellite_id, params, datum_schema):
    offset = params['beacon_offset'].split('.')
    byte_offset = int(offset[0])
    if len(offset) == 2:
        bit_offset = int(offset[1])

    byte_offset += 16  # skip AX.25 header

    unit, linear_transform = GetUnitEnum(params['unit'])

    if params['datatype'] == 'int':
        length = int(params['data_length'])
        signed = bool(int(params['signed']))
        int_t = ({True: 's', False: 'u'}[signed] +
                 {1: 'int8', 2: 'int16', 4: 'int32'}[length])

        # Little endian.
        parts = []
        for i in range(length):
            parts.append('%s(frame[%d]) << %d' % (
                    int_t, byte_offset+i, 8*i))
        raw = ' | '.join(parts)

        min_length = byte_offset + length

        if datum_schema.type == telemetry_pb2.TelemetryDatumSchema.DOUBLE:
            return min_length, """
	{
		raw := %(raw)s
		X := float64(raw)
		y := %(conversion_func)s
		z := %(m)f * y + %(c)f
		data = append(data, NewDoubleDatum("%(key)s", timestamp, z))
	}""" % {
                'raw': raw,
                'conversion_func': params['conversion_func'],
                'key': datum_schema.key,
                'm': linear_transform[0],
                'c': linear_transform[1]}
        elif datum_schema.type == telemetry_pb2.TelemetryDatumSchema.INT64:
            return min_length, """
	{
		raw := %(raw)s
		data = append(data, NewInt64Datum("%(key)s", timestamp, int64(raw)))
	}""" % {
                'raw': raw,
                'key': datum_schema.key}
        else:
            assert False

    elif params['datatype'] == 'bool':
        assert datum_schema.type == telemetry_pb2.TelemetryDatumSchema.BOOL
	return byte_offset+1, """
	{
		raw := (frame[%(byte_offset)d] >> %(bit_offset)d) & 1 == 1
		data = append(data, NewBoolDatum("%(key)s", timestamp, raw))
	}""" % {
            'byte_offset': byte_offset,
            'bit_offset': bit_offset,
            'key': datum_schema.key}

    elif params['datatype'] == 'lut':
        assert datum_schema.type == telemetry_pb2.TelemetryDatumSchema.DOUBLE
        assert int(params['data_length']) == 1

        mappings = params['lut_contents'].split(',')
        items = []
        keys = []
        for m in mappings:
            k, v = m.split('=')
            k = int(k)
            v = float(v)
            keys.append(k)
            items.append('\t\t\t%d: %f,' % (k, v))

        keys.sort()
        assert keys == range(256)

        return byte_offset+1, """
	{
		m := map[byte]float64{
%(table)s
		}
		y := m[frame[%(byte_offset)d]]
		z := %(m)f * y + %(c)f
		data = append(data, NewDoubleDatum("%(key)s", timestamp, z))
	}""" % {
            'table': '\n'.join(items),
            'byte_offset': byte_offset,
            'key': datum_schema.key,
            'm': linear_transform[0],
            'c': linear_transform[1]}

    else:
        assert False


class XMLHandler(xml.sax.ContentHandler):
    def __init__(self, satellite_id):
        self.params = None
        self.key = None
        self.satellite_id = satellite_id
        self.schema = []
        self.code = []
        self.min_length = 0

    def startElement(self, name, attrs):
        if name == 'telemetry_types':
            self.params = {}
            return

        if self.params is not None:
            self.key = name
            self.params[name] = ''

    def endElement(self, name):
        if name == 'telemetry_types':
            stripped = {}
            for k, v in self.params.items():
                stripped[k] = v.strip()
            ds = ConvertParamsToSchema(self.satellite_id, stripped)
            if ds:
                self.schema.append(ds)
                min_length, code = GenerateCodeForDatum(
                    self.satellite_id, stripped, ds)
                if code:
                    self.code.append(code)
                    self.min_length = max(self.min_length, min_length)

            self.params = None
            self.key = None

    def characters(self, content):
        if self.key is not None:
            self.params[self.key] = self.params[self.key] + content


def ConvertXMLToSchema(satellite_id, input_path):
    f = file(input_path)

    handler = XMLHandler(satellite_id)
    parser = xml.sax.make_parser()
    parser.setContentHandler(handler)
    parser.parse(f)

    schema = telemetry_pb2.TelemetrySchema()
    schema.datum.extend(handler.schema)

    blocks = [code_start % {
            'satellite_id': satellite_id,
            'min_length': handler.min_length,
            'input_path': input_path,
            'script_name': sys.argv[0],
            'date': time.ctime()}]
    blocks.extend(handler.code)
    blocks.append(code_end)
    code = ''.join(blocks)

    return schema, code


satellite_id = 'csswe'

schema, code = ConvertXMLToSchema(satellite_id, '/Users/tstranex/Downloads/csswe/csswe_public_gs_dist/CSSWE_GS_Client/telemetry_types.xml')

file('src/carpcomm/telemetry/%s.go' % satellite_id, 'w').write(code)

print schema
