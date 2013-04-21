// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "fmt"
import "errors"
import "carpcomm/pb"
import "math"

// Some documentation:
// http://www.isi.edu/projects/serc/aeneas_telemetry_data_qsl_card
func DecodeAeneas(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {

	if len(frame) != 69 {
		return nil, errors.New(fmt.Sprintf(
			"Invalid length: %d, expected 69.", len(frame)))
	}

	// Packet type
	packet_type := frame[8]
	if packet_type != 0x02 {
		return nil, errors.New(fmt.Sprintf(
			"Unknown packet type: %d, expected 2.", packet_type));
	}

	// Reboot count
	data = append(data, NewInt64Datum(
		"d:aeneas:num_reboots", timestamp, int64(frame[18])))

	// Gyro rates
	gyroRate := func(buf []byte) (rate_rad_per_s float64) {
		raw := int16(buf[1]) | int16(buf[0])<<8
		rate_deg_per_s := 0.000291075 * float64(raw) - 0.0473675
		return rate_deg_per_s * math.Pi / 180.0
	}
	data = append(data, NewDoubleDatum(
		"d:aeneas:gyro1_r", timestamp, gyroRate(frame[55:57])))
	data = append(data, NewDoubleDatum(
		"d:aeneas:gyro2_r", timestamp, gyroRate(frame[57:59])))
	data = append(data, NewDoubleDatum(
		"d:aeneas:gyro3_r", timestamp, gyroRate(frame[59:61])))

	// Gyro temperatures
	gyroTemp := func(buf []byte) (temp_kelvin float64) {
		raw := int16(buf[1]) | int16(buf[0])<<8
		temp_celcius := 0.000564399 * float64(raw) + 24.9691
		return temp_celcius + 273.15
	}
	data = append(data, NewDoubleDatum(
		"d:aeneas:gyro1_t", timestamp, gyroTemp(frame[61:63])))
	data = append(data, NewDoubleDatum(
		"d:aeneas:gyro2_t", timestamp, gyroTemp(frame[63:65])))
	data = append(data, NewDoubleDatum(
		"d:aeneas:gyro3_t", timestamp, gyroTemp(frame[65:67])))

	return data, nil
}