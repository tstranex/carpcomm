// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "encoding/hex"
import "errors"
import "fmt"
import "carpcomm/pb"

// Documentation:
// http://en.wikipedia.org/wiki/TechEdSat
// http://techedsat.com/index.php?option=com_content&view=article&id=32&Itemid=187

func techedsat_DecodeHex(data string) (int64, error) {
	if len(data) == 0 {
		return 0, errors.New("Empty data")
	}

	var value int64
	for _, c := range data {
		bytes, err := hex.DecodeString("0" + (string)(c))
		if err != nil || len(bytes) != 1 {
			return 0, errors.New(fmt.Sprintf("Hex decode error: %s",
				err.Error()))
		}
		value = (value << 4) + (int64)(bytes[0])
	}
	return value, nil
}

func techedsat_DecodeADC(key string, timestamp int64,
	data string, adc_zero, adc_full, m, c float64) (
	datum pb.TelemetryDatum, err error) {

	v, err := techedsat_DecodeHex(data)
	if err != nil {
		return datum, err
	}

	denominator := (1 << (4 * (uint)(len(data)))) - 1
	adc_val := adc_zero +
		float64(v) * (adc_full - adc_zero) / float64(denominator)
	datum_val := m*adc_val + c

	return NewDoubleDatum(key, timestamp, datum_val), nil
}

func techedsat_DecodeDouble(
	key string, timestamp int64, data string, m, c float64) (
	datum pb.TelemetryDatum, err error) {
	
	v, err := techedsat_DecodeHex(data)
	if err != nil {
		return datum, err
	}

	datum_val := float64(v) * m + c
	return NewDoubleDatum(key, timestamp, datum_val), nil
}

func techedsat_DecodeInt(key string, timestamp int64, data string) (
	datum pb.TelemetryDatum, err error) {
	
	v, err := techedsat_DecodeHex(data)
	if err != nil {
		return datum, err
	}

	return NewInt64Datum(key, timestamp, v), nil
}

func techedsat_DecodeSolarArrayStatus(timestamp int64, data string) (
	r []pb.TelemetryDatum, err error) {

	v, err := techedsat_DecodeHex(data)
	if err != nil {
		return nil, err
	}

	keys := []string{
		"d:techedsat:cell5_active",
		"d:techedsat:cell4_active",
		"d:techedsat:cell3_active",
		"d:techedsat:cell2_active",
		"d:techedsat:cell1_active",
	}
	for i, key := range keys {
		mask := 1 << (uint)(i)
		value := ((int)(v) & mask) > 0
		r = append(r, NewBoolDatum(key, timestamp, value))
	}

	return r, nil
}

func techedsat_DecodePayload(payload string, timestamp int64) (
	data []pb.TelemetryDatum, err error) {

	if len(payload) < 122 {
		return nil, errors.New("Data payload too short.")
	}

	var d pb.TelemetryDatum

	if payload[0:10] != "ncasst.org" {
		return nil, errors.New("Missing ncasst.org tag.")
	}

// d:techedsat:ncasst.org payload[0:10]
d, err = techedsat_DecodeDouble("d:techedsat:elapsed_s", timestamp, payload[10:18], 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:5v_min_v", timestamp, payload[18:21], 0.000000, 6.257160, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:5v_max_v", timestamp, payload[21:24], 0.000000, 6.257160, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:5v_avg_v", timestamp, payload[24:27], 0.000000, 6.257160, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:5v_now_v", timestamp, payload[27:30], 0.000000, 6.257160, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:3v3_min_v", timestamp, payload[30:33], 0.000000, 6.257160, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:3v3_max_v", timestamp, payload[33:36], 0.000000, 6.257160, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:3v3_avg_v", timestamp, payload[36:39], 0.000000, 6.257160, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:3v3_now_v", timestamp, payload[39:42], 0.000000, 6.257160, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:min_t", timestamp, payload[42:45], -413.714286, 171.285714, 1.000000, 273.150000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:max_t", timestamp, payload[45:48], -413.714286, 171.285714, 1.000000, 273.150000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:avg_t", timestamp, payload[48:51], -413.714286, 171.285714, 1.000000, 273.150000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:now_t", timestamp, payload[51:54], -413.714286, 171.285714, 1.000000, 273.150000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:min_c", timestamp, payload[54:57], 0.000000, 0.293980, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:max_c", timestamp, payload[57:60], 0.000000, 0.293980, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:avg_c", timestamp, payload[60:63], 0.000000, 0.293980, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:now_c", timestamp, payload[63:66], 0.000000, 0.293980, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:bat_min_v", timestamp, payload[66:69], 0.000000, 9.801528, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:bat_max_v", timestamp, payload[69:72], 0.000000, 9.801528, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:bat_avg_v", timestamp, payload[72:75], 0.000000, 9.801528, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:bat_now_v", timestamp, payload[75:78], 0.000000, 9.801528, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:bat_min_c", timestamp, payload[78:81], 0.000000, 3.315940, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:bat_max_c", timestamp, payload[81:84], 0.000000, 3.315940, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:bat_avg_c", timestamp, payload[84:87], 0.000000, 3.315940, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:bat_now_c", timestamp, payload[87:90], 0.000000, 3.315940, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:charge_min_c", timestamp, payload[90:93], -3.484321, 3.482619, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:charge_max_c", timestamp, payload[93:96], -3.484321, 3.482619, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:charge_avg_c", timestamp, payload[96:99], -3.484321, 3.482619, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeADC("d:techedsat:charge_now_c", timestamp, payload[99:102], -3.484321, 3.482619, 1.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

// d:techedsat:sa_status payload[102:104]
d, err = techedsat_DecodeDouble("d:techedsat:nominal_mode_s", timestamp, payload[104:108], 60.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeDouble("d:techedsat:safe_mode_s", timestamp, payload[108:112], 60.000000, 0.000000)
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeInt("d:techedsat:single_errors", timestamp, payload[112:114])
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeInt("d:techedsat:all_errors", timestamp, payload[114:116])
if err != nil {
  return nil, err
}
data = append(data, d)

d, err = techedsat_DecodeInt("d:techedsat:all_errors_2", timestamp, payload[116:118])
if err != nil {
  return nil, err
}
data = append(data, d)

// d:techedsat:crc_valid payload[118:122]

	sa_data, err := techedsat_DecodeSolarArrayStatus(
		timestamp, payload[102:104])
	if err != nil {
		return nil, err
	}
	data = append(data, sa_data...)

	return data, nil
}

func DecodeFrame_techedsat(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {
	if len(frame) < 16 {
		return nil, errors.New("Frame too short")
	}
	return techedsat_DecodePayload((string)(frame[16:]), timestamp)
}
