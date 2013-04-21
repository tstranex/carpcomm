// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "errors"
import "strings"
import "encoding/hex"
import "carpcomm/pb"
import "fmt"
import "carpcomm/util"

const kFitsat1Call = "HIDENIWAKAJAPAN"

func fitsat1DecodeByte(key string, timestamp int64, b byte, m, c float64) (
	pb.TelemetryDatum, error) {
	v := float64(b) * m + c
	return NewDoubleDatum(key, timestamp, v), nil
}

func fitsat1DecodePart(decoded_morse string, timestamp int64) (
	data []pb.TelemetryDatum, err error) {

	if len(decoded_morse) < 10 {
		return nil, errors.New("Message too short")
	}
	
	part := decoded_morse[:2]
	bytes, err := hex.DecodeString(decoded_morse[2:10])
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Hex decode error: %s", err.Error()))
	}
	v := make([]float64, len(bytes))
	for i, b := range bytes {
		v[i] = (float64)(b)
	}

	if part == "S1" {
		const f = 5.0 / 256.0
		data = append(data, NewDoubleDatum(
			"d:fitsat1:rssi_437", timestamp, v[0]*f))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:total_cell_v", timestamp, v[1]*f))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:total_cell_c", timestamp, v[2]*f*0.4))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:bat1_v", timestamp, v[3]*f))

	} else if part == "S2" {
		const f = 5.0 / 256.0
		data = append(data, NewDoubleDatum(
			"d:fitsat1:bat1_c", timestamp, (v[0]*f - 2.5)*0.4))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:bat3_v", timestamp, v[1]*f))

		c := v[2]*f - 2.5
		if c > 0.0 {
			c *= 10.0
		} else {
			c *= 0.1
		}
		data = append(data, NewDoubleDatum(
			"d:fitsat1:bat3_c", timestamp, c))

		data = append(data, NewDoubleDatum(
			"d:fitsat1:2p5_bus_v", timestamp, v[3]*f))

	} else if part == "S3" {
		const f = 4.5 / 256.0 * 2
		data = append(data, NewDoubleDatum(
			"d:fitsat1:cell_x2_v", timestamp, v[0]*f))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:cell_y2_v", timestamp, v[1]*f))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:cell_x1_v", timestamp, v[2]*f))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:cell_y1_v", timestamp, v[3]*f))

	} else if part == "S4" {
		tr := func(v float64) float64 {
			return (v*4.5 / 256.0 - 0.5) / 0.01 + 273.15
		}
		data = append(data, NewDoubleDatum(
			"d:fitsat1:bat3_t", timestamp, tr(v[0])))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:bat1_t", timestamp, tr(v[1])))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:panel_z2_t", timestamp, tr(v[2])))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:panel_z1_t", timestamp, tr(v[3])))

	} else if part == "S5" {
		data = append(data, NewDoubleDatum(
			"d:fitsat1:rssi_1200", timestamp, v[0] * 4.5 / 256.0))
		data = append(data, NewDoubleDatum(
			"d:fitsat1:time_since_boot_s",
			timestamp,
			v[1]*65536 + v[2]*256 + v[3]))

	} else {
		return nil, errors.New("Invalid message type")
	}

	return data, nil
}

func DecodeFitsat1Morse(decoded_morse string, timestamp int64)  (
	data []pb.TelemetryDatum, err error) {

	decoded_morse = util.StripWhitespace(decoded_morse)

	decoded_morse = strings.ToUpper(decoded_morse)
	if decoded_morse == kFitsat1Call {
		return nil, nil
	}

	for i := 0; i < len(decoded_morse)-9; i++ {
		d, _ := fitsat1DecodePart(decoded_morse[i:], timestamp)
		data = append(data, d...)
		if len(d) > 0 {
			i += 9
		}
	}

	if len(data) == 0 {
		return nil, errors.New("No valid messages found")
	}

	return data, nil
}