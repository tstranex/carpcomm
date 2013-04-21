#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

# Script to convert TechEdSat schema into our schema and generate a frame
# decoder.

# TechEdSat schema comes from: http://en.wikipedia.org/wiki/TechEdSat

from carpcomm.pb import telemetry_pb2
from carpcomm.pb import text_pb2


techedsat_schema = """
ncasst.org	10	0	Location of spacecraft information	URL
SCETTime	8	10	Spacecraft Elapsed Time	seconds
5V_min	3	18	SM 5V HK min value	12 bit raw ADC
5V_max	3		SM 5V HK max value	12 bit raw ADC
5V_avg	3	24	SM 5V HK rolling average	12 bit raw ADC
5V_now	3		SM 5V HK last value	12 bit raw ADC
3V3_min	3	30	SM 3V3 HK last value	12 bit raw ADC
3V3_max	3		SM 3V3 HK last value	12 bit raw ADC
3V3_avg	3	36	SM 3V3 HK last value	12 bit raw ADC
3V3_now	3		SM 3V3 HK last value	12 bit raw ADC
Temp_min	3	42	SM Temperature HK min value	12 bit raw ADC
Temp_max	3		SM Temperature HK max value	12 bit raw ADC
Temp_avg	3	48	SM Temperature HK rolling average	12 bit raw ADC
Temp_now	3		SM Temperature HK last value	12 bit raw ADC
Curr_min	3	54	SM Current HK min value	12 bit raw ADC
Curr_max	3		SM Current HK max value	12 bit raw ADC
Curr_avg	3	60	SM Current HK rolling average	12 bit raw ADC
Curr_now	3		SM Battery Bus Voltage last value	12 bit raw ADC
VBUS_min	3	66	SM Battery Bus Voltage max value	12 bit raw ADC
VBUS_max	3		SM Battery Bus Voltage last value	12 bit raw ADC
VBUS_avg	3	72	SM Battery Bus Voltage rolling average	12 bit raw ADC
VBUS_now	3		SM Battery Bus Voltage last value	12 bit raw ADC
IBUS_min	3	78	SM Battery Bus Current min value	12 bit raw ADC
IBUS_max	3		SM Battery Bus Current max value	12 bit raw ADC
IBUS_avg	3	84	SM Battery Bus Current rolling average	12 bit raw ADC
IBUS_now	3		SM Battery Bus Current last value	12 bit raw ADC
ICH_min	3	90	SM Battery Charge Current min value	12 bit raw ADC
ICH_max	3		Battery Charge Current max value	12 bit raw ADC
ICH_avg	3	96	SM Battery Charge Current rolling average	12 bit raw ADC
ICH_now	3		SM Battery Charge Current last value	12 bit raw ADC
SA_Status	2	102	SM Solar Array Status Byte	SA1: Bit 0 SA2: Bit 1 SA3: Bit 2 SA4: Bit 3 SA5: Bit 4
NON_minutes	4	104	How long the spacecraft has been in Nominal Mode	minutes
NOF_minutes	4	108	How long the spacecraft has been in Safe Mode	minutes
Single Errors	2	112	Safe Mode Processor Error Counter	number of errors
All Errors	2		Safe Mode Processor Error Counter	number of errors
All Errors 2	2		Safe Mode Processor Error Counter	number of errors
CRC	4	118	CRC over all previous data bytes CRC-16-IBM polynomial x16 + x15 + x2 + 1\t0xA001 Key
"""

# From http://techedsat.com/index.php?option=com_content&view=article&id=32&Itemid=187
adc_calibration = {
    '5V_min': (0.0, 6.25716),
    '5V_max': (0.0, 6.25716),
    '5V_avg': (0.0, 6.25716),
    '5V_now': (0.0, 6.25716),
    '3V3_min': (0.0, 6.25716),
    '3V3_max': (0.0, 6.25716),
    '3V3_avg': (0.0, 6.25716),
    '3V3_now': (0.0, 6.25716),
    'Temp_min': (-413.7142857, 171.2857143),
    'Temp_max': (-413.7142857, 171.2857143),
    'Temp_avg': (-413.7142857, 171.2857143),
    'Temp_now': (-413.7142857, 171.2857143),
    'Curr_min': (0.0, 0.29398005),
    'Curr_max': (0.0, 0.29398005),
    'Curr_avg': (0.0, 0.29398005),
    'Curr_now': (0.0, 0.29398005),
    'VBUS_min': (0.0, 9.801528),
    'VBUS_max': (0.0, 9.801528),
    'VBUS_avg': (0.0, 9.801528),
    'VBUS_now': (0.0, 9.801528),
    'IBUS_min': (0.0, 3.3159398),
    'IBUS_max': (0.0, 3.3159398),
    'IBUS_avg': (0.0, 3.3159398),
    'IBUS_now': (0.0, 3.3159398),
    'ICH_min': (-3.48432056, 3.48261923),
    'ICH_max': (-3.48432056, 3.48261923),
    'ICH_avg': (-3.48432056, 3.48261923),
    'ICH_now': (-3.48432056, 3.48261923),
}


def GetUnitEnum(name):
    if (name.startswith('5V') or
        name.startswith('3V') or
        name.startswith('VBUS')):
        return telemetry_pb2.TelemetryDatumSchema.VOLT, (1.0, 0.0)
    if name.startswith('Temp'):
        return telemetry_pb2.TelemetryDatumSchema.KELVIN, (1.0, 273.15)
    if (name.startswith('Curr') or
        name.startswith('IBUS') or
        name.startswith('ICH')):
        return telemetry_pb2.TelemetryDatumSchema.AMPERE, (1.0, 0.0)
    if name == 'SCETTime':
        return telemetry_pb2.TelemetryDatumSchema.SECOND, (1.0, 0.0)
    if name.endswith('_minutes'):
        return telemetry_pb2.TelemetryDatumSchema.SECOND, (60.0, 0.0)
    return None, (None, None)

def GetKey(name, unit):
    unit_suffix = {
        telemetry_pb2.TelemetryDatumSchema.VOLT: '_v',
        telemetry_pb2.TelemetryDatumSchema.KELVIN: '_t',
        telemetry_pb2.TelemetryDatumSchema.AMPERE: '_c',
        telemetry_pb2.TelemetryDatumSchema.SECOND: '_s',
        }

    strip_prefixes = [
        'curr_',
        'temp_']

    replaces = {
        'ibus': 'bat',
        'vbus': 'bat',
        'ich': 'charge',
        'non_minutes': 'nominal_mode',
        'nof_minutes': 'safe_mode',
        'scettime': 'elapsed',
        'crc': 'crc_valid',
        }

    n = name.strip().lower()
    for ss in strip_prefixes:
        if n.startswith(ss):
            n = n[len(ss):]
    for k, v in replaces.items():
        n = n.replace(k, v)
    n = n.strip()
    n = n.replace(' ', '_')

    return 'd:techedsat:%s%s' % (n, unit_suffix.get(unit, ''))


adc_code = """d, err = techedsat_DecodeADC("%s", timestamp, payload[%d:%d], %f, %f, %f, %f)
if err != nil {
  return nil, err
}
data = append(data, d)
"""

double_code = """d, err = techedsat_DecodeDouble("%s", timestamp, payload[%d:%d], %f, %f)
if err != nil {
  return nil, err
}
data = append(data, d)
"""

int_code = """d, err = techedsat_DecodeInt("%s", timestamp, payload[%d:%d])
if err != nil {
  return nil, err
}
data = append(data, d)
"""


schema = telemetry_pb2.TelemetrySchema()
code = []

for line in techedsat_schema.split('\n'):
    if not line.strip():
        continue
    name, size, offset, desc, unit_column = line.split('\t')
    size = int(size.strip())
    offset = offset.strip()
    if offset:
        offset = int(offset)
    else:
        offset = last_offset + last_size
    last_offset = offset
    last_size = size

    s = telemetry_pb2.TelemetryDatumSchema()

    n = s.name.add()
    n.lang = 'en'
    n.text = desc

    unit, linear_transform = GetUnitEnum(name)
    if unit is not None:
        s.unit = unit
        s.type = telemetry_pb2.TelemetryDatumSchema.DOUBLE

    if 'number' in unit_column:
        s.type = telemetry_pb2.TelemetryDatumSchema.INT64

    s.key = GetKey(name, unit)
    s.source_key = name

    schema.datum.extend([s])

    if unit_column.strip() == '12 bit raw ADC':
        adc_zero, adc_full = adc_calibration[name]
        m, c = linear_transform
        code.append(adc_code % (
                s.key, offset, offset+size, adc_zero, adc_full, m, c))
    elif s.type == telemetry_pb2.TelemetryDatumSchema.DOUBLE:
        m, c = linear_transform
        code.append(double_code % (s.key, offset, offset+size, m, c))
    elif s.type == telemetry_pb2.TelemetryDatumSchema.INT64:
        code.append(int_code % (s.key, offset, offset+size))
    else:
        code.append('// %s payload[%d:%d]' % (s.key, offset, offset+size))

s = str(schema)
for line in s.split('\n'):
    print '    ' + line

print
print '\n'.join(code)
