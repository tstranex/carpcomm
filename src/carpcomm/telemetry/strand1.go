// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "carpcomm/pb"
import "errors"
import "time"

func strand1ModemBeacon(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {
	if len(frame) < 5 {
		return nil, errors.New("Modem beacon frame too short")
	}

	var int_value int64
	for i := 1; i < 5; i++ {
		int_value = (int_value << 8) + int64(frame[i])
	}

	if frame[0] == 0xe0 {
		data = append(data, NewDoubleDatum(
			"d:strand1:time_since_last_obc_packet",
			timestamp,
			float64(int_value)))
		return data, nil
	} else if frame[0] == 0xe1 {
		data = append(data, NewInt64Datum(
			"d:strand1:packet_up_count",
			timestamp,
			int_value))
		return data, nil
	} else if frame[0] == 0xe2 {
		data = append(data, NewInt64Datum(
			"d:strand1:packet_down_count",
			timestamp,
			int_value))
		return data, nil
	} else if frame[0] == 0xe3 {
		data = append(data, NewInt64Datum(
			"d:strand1:packet_up_dropped_count",
			timestamp,
			int_value))
		return data, nil
	} else if frame[0] == 0xe4 {
		data = append(data, NewInt64Datum(
			"d:strand1:packet_down_dropped_count",
			timestamp,
			int_value))
		return data, nil
	}

	return nil, nil
}

type strand1ADC struct {
	m, c float64
	key string
	unit_scale float64
}
var strand1EPS = map[byte]strand1ADC{
	0x01: {-3.4969, 3185.1551, "d:strand1:bat0_c", 0.001},
	0x03: {-0.00945, 9.7488, "d:strand1:bat0_v", 1},
	0x04: {-0.163, 111.187 + 273.15, "d:strand1:bat0_t", 1},
	0x06: {-3.4768, 3173.1106, "d:strand1:bat1_c", 0.001},
	0x08: {-0.00946, 9.7526, "d:strand1:bat1_v", 1},
	0x09: {-0.163, 111.187 + 273.15, "d:strand1:bat1_t", 1},	
}
var strand1Battery = map[byte]strand1ADC{
	0x01: {-0.542490348, 528.0441026, "d:strand1:cell_y2_c", 0.001},
	0x02: {-0.163, 110.338 + 273.15, "d:strand1:cell_y2_t", 1},
	0x03: {-0.035254639, 34.6505381, "d:strand1:cell_y_v", 1},
	0x04: {-0.537846059, 523.1519466, "d:strand1:cell_y1_c", 0.001},
	0x05: {-0.163, 110.338 + 273.15, "d:strand1:cell_y1_t", 1},
	0x06: {-0.035579727, 34.76510021, "d:strand1:cell_x_v", 1},
	0x07: {-0.541228423, 526.8412823, "d:strand1:cell_x1_c", 0.001},
	0x08: {-0.163, 110.338 + 273.15, "d:strand1:cell_x1_t", 1},
	0x09: {-0.00914561, 8.782534345, "d:strand1:cell_z_v", 1},
	0x0a: {-0.52264946, 508.5204547, "d:strand1:cell_z2_c", 0.001},
	0x0b: {-0.163, 110.338 + 273.15, "d:strand1:cell_z2_t", 1},
	0x0d: {-0.518702129, 512.807352, "d:strand1:cell_x2_c", 0.001},
	0x0e: {-0.163, 110.338 + 273.15, "d:strand1:cell_x2_t", 1},
	0x11: {-4.926127936, 4414.027999, "d:strand1:bat_bus_c", 0.001},
	0x1a: {-5.431052862, 4636.008505, "d:strand1:5v_bus_c", 0.001},
	0x1b: {-3.626006798, 3080.538997, "d:strand1:3p3v_bus_c", 0.001},
	0x1e: {-0.163, 110.338 + 273.15, "d:strand1:cell_z1_t", 1},
	0x1f: {-0.52947555, 515.5141451, "d:strand1:cell_z1_c", 0.001},
}
var strand1SwitchCurrent = map[byte]strand1ADC{
	0x81: {0.259549, -1.516825, "d:strand1:switch0_c", 0.001},
	0x86: {0.258359, -1.554162, "d:strand1:switch1_c", 0.001},
	0x8b: {0.259325, -1.595903, "d:strand1:switch2_c", 0.001},
	0x90: {0.518526, -8.756971, "d:strand1:switch3_c", 0.001},
	0x95: {0.534516, -3.25046, "d:strand1:switch4_c", 0.001},
	0x9a: {0.528245, -2.974109, "d:strand1:switch5_c", 0.001},
	0x9f: {0.260476, -0.91132, "d:strand1:switch6_c", 0.001},
	0xa4: {0.532941, -3.152331, "d:strand1:switch7_c", 0.001},
}
var strand1SwitchVoltage = map[byte]strand1ADC{
	0x81: {2.300107, -1113.424579, "d:strand1:switch0_v", 0.001},
	0x86: {2.315349, -1136.056829, "d:strand1:switch1_v", 0.001},
	0x8b: {2.3315, -1187.043977, "d:strand1:switch2_v", 0.001},
	0x90: {3.667785, -7266.803691, "d:strand1:switch3_v", 0.001},
	0x95: {2.603641, -0.504061, "d:strand1:switch4_v", 0.001},
	0x9a: {2.233264, -930.303516, "d:strand1:switch5_v", 0.001},
	0x9f: {2.254974, -993.915009, "d:strand1:switch6_v", 0.001},
	0xa4: {2.592693, 3.656067, "d:strand1:switch7_v", 0.001},
}
var strand1SwitchOn = map[byte]string{
	0x81: "d:strand1:switch0_on",
	0x86: "d:strand1:switch1_on",
	0x8b: "d:strand1:switch2_on",
	0x90: "d:strand1:switch3_on",
	0x95: "d:strand1:switch4_on",
	0x9a: "d:strand1:switch5_on",
	0x9f: "d:strand1:switch6_on",
	0xa4: "d:strand1:switch7_on",
}

func strand1OBCBeacon(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {
	if len(frame) < 3 {
		return nil, errors.New("OBC beacon frame too short")
	}

	node := frame[0]
	channel := frame[1]
	size := frame[2]
	if len(frame) - 3 < int(size) {
		return nil, errors.New("Invalid size")
	}

	int_value := 0
	for i := 0; i < int(size); i++ {
		int_value = int_value << 8
		int_value = int_value | int(frame[len(frame)-i-1])
	}
	floatval := float64(int_value)

	if node == 0x2c {
		if channel == 0x00 {
			data = append(data, NewBoolDatum(
				"d:strand1:bat0_charge",
				timestamp,
				int_value < 30))
			return data, nil
		} else if channel == 0x05 {
			data = append(data, NewBoolDatum(
				"d:strand1:bat1_charge",
				timestamp,
				int_value < 30))
			return data, nil
		}
		eq, ok := strand1EPS[channel]
		if !ok {
			return data, errors.New("Unknown telemetry channel")
		}
		v := (eq.m * floatval + eq.c) * eq.unit_scale
		data = append(data, NewDoubleDatum(eq.key, timestamp, v))
		return data, nil

	} else if node == 0x2d {
		eq, ok := strand1Battery[channel]
		if !ok {
			return data, errors.New("Unknown telemetry channel")
		}
		v := (eq.m * floatval + eq.c) * eq.unit_scale
		data = append(data, NewDoubleDatum(eq.key, timestamp, v))
		return data, nil

	} else if node == 0x66 {
		// Switches

		on_key, ok := strand1SwitchOn[channel]
		if !ok {
			return nil, errors.New("Unknown switch")
		}

		state := frame[3]
		data = append(data, NewBoolDatum(
			on_key, timestamp, state == 0x01 || state == 0x05))

		c_eq, ok := strand1SwitchCurrent[channel]
		v_eq, ok := strand1SwitchVoltage[channel]
		if !ok {
			return data, nil
		}

		current := (frame[4] << 8) | frame[5]
		voltage := (frame[6] << 8) | frame[7]

		data = append(data, NewDoubleDatum(c_eq.key, timestamp,
			(c_eq.m * float64(current) + c_eq.c) * c_eq.unit_scale))
		data = append(data, NewDoubleDatum(v_eq.key, timestamp,
			(v_eq.m * float64(voltage) + v_eq.c) * v_eq.unit_scale))
		return data, nil

	} else if node == 0x80 {
		if channel == 0x0c {
			var clock int64 = int64(frame[3]) |
				(int64(frame[4]) << 8) |
				(int64(frame[5]) << 16) |
				(int64(frame[6]) << 24)
			data = append(data, NewTimestampDatum(
				"d:strand1:obc_clock", timestamp, 
				time.Unix(clock, 0)))
			return data, nil
		}

	} else if node == 0x89 {
		var num1 int32 = int32(frame[3]) |
			(int32(frame[4]) << 8) |
			(int32(frame[5]) << 16) |
			(int32(frame[6]) << 24)
		var num2 int32 = int32(frame[7]) |
			(int32(frame[8]) << 8) |
			(int32(frame[9]) << 16) |
			(int32(frame[10]) << 24)

		if channel == 0x03 {
			data = append(data, NewInt64Datum(
				"d:strand1:magnetometer_1x", timestamp,
				int64(num1)))
			data = append(data, NewInt64Datum(
				"d:strand1:magnetometer_1y", timestamp,
				int64(num2)))
			return data, nil
		} else if channel == 0x05 {
			data = append(data, NewInt64Datum(
				"d:strand1:magnetometer_2z", timestamp,
				int64(num1)))
			return data, nil
		}
	}

	return nil, errors.New("Unknown telemetry value")
}

// Example: C0 80 A5 06 02 2C 08 02 AD 00
func decodeFrame_strand1(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {

	// Some packets have the header, some don't. Weird...
	if len(frame) >= 4 && frame[0] == 0xc0 && frame[1] == 0x80 {
		frame = frame[4:]
	}

	if len(frame) < 1 {
		return nil, errors.New("Frame too short")
	}
	frame_type := frame[0]

	if frame_type == 0x01 {
		return strand1ModemBeacon(frame[1:], timestamp)
	} else if frame_type == 0x02 {
		return strand1OBCBeacon(frame[1:], timestamp)
	}
	return nil, errors.New("Unknown frame type")
}

// http://ukamsat.files.wordpress.com/2013/02/strand-1-telemetry-information.png
// 
func DecodeFrame_strand1(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {
	return decodeFrame_strand1(frame, timestamp)
}
