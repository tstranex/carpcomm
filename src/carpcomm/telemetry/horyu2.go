// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "fmt"
import "errors"
import "strings"
import "encoding/hex"
import "carpcomm/pb"

const kHoryu2Callsign = "JG6YBWHORYU"

// Documentation:
// http://kitsat.ele.kyutech.ac.jp/Documents/ground_station/english/explain_downlink_information_Eng.pdf
// http://kitsat.ele.kyutech.ac.jp/Documents/ground_station/explain_downlink_information_jap.pdf
func decodeHoryu2(decoded_morse string, timestamp int64)  (
	data []pb.TelemetryDatum, err error) {

	decoded_morse = strings.ToUpper(decoded_morse)

	// Callsign
	if decoded_morse == kHoryu2Callsign {
		return nil, nil
	}

	if len(decoded_morse) != 15 {
		return nil, errors.New(fmt.Sprintf(
			"Invalid length: %d, expected 15.",
			len(decoded_morse)))
	}

	// Append a 0 digit so that we have an even number of digits.
	hex_digits := decoded_morse + "0"
	bytes, err := hex.DecodeString(hex_digits)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Hex decoding error: %s", err.Error()))
	}
	if len(bytes) != 8 {
		// This shouldn't happen.
		return nil, errors.New(fmt.Sprintf(
			"Internal hex decode error: %s", err.Error()))
	}

	temperature := func(b byte, m0, m1, m2, m3 float64) float64 {
		x := float64(b) * 16
		t_celcius := m0 + m1*x + m2*x*x + m3*x*x*x
		return t_celcius + 273.15
	}

	// Byte 0: Vref
	// Calibration formula is unknown!

	// Byte 1: Battery top temperature
	bat_top_t := temperature(bytes[1],
		-18762.6337, 20.1087835, -7.21222547E-03, 8.65463839E-07)
	data = append(data, NewDoubleDatum(
		"d:horyu2:bat_top_t", timestamp, bat_top_t))

	// Byte 2: Battery bottom temperature
	bat_bottom_t := temperature(bytes[2],
		-76833.1906, 79.1616336, -2.7213206E-02, 3.12158024E-06)
	data = append(data, NewDoubleDatum(
		"d:horyu2:bat_bottom_t", timestamp, bat_bottom_t))

	// Byte 3: Comm temperature
	com_t := temperature(bytes[3],
		-64672.493, 67.3602286, -2.34063489E-02, 2.71339774E-06)
	data = append(data, NewDoubleDatum("d:horyu2:comm_t", timestamp, com_t))

	// Byte 4: Battery current
	bat_c := (1.596 * 16 * float64(bytes[4]) - 3309) / 1000.0
	data = append(data, NewDoubleDatum("d:horyu2:bat_c", timestamp, bat_c))

	// Byte 5: Battery voltage
	bat_v := (1.22 * 16 * float64(bytes[5]) + 63) / 1000.0
	data = append(data, NewDoubleDatum("d:horyu2:bat_v", timestamp, bat_v))

	// Byte 6: Status bits
	data = append(data, NewBoolDatum(
		"d:horyu2:clock_normal", timestamp, bytes[6] & 128 > 0))
	data = append(data, NewBoolDatum(
		"d:horyu2:flash_main_normal", timestamp, bytes[6] & 64 > 0))
	data = append(data, NewBoolDatum(
		"d:horyu2:flash_share_normal", timestamp, bytes[6] & 32 > 0))
	data = append(data, NewBoolDatum(
		"d:horyu2:flash_300_normal", timestamp, bytes[6] & 16 > 0))

	data = append(data, NewBoolDatum(
		"d:horyu2:switch_share_normal", timestamp, bytes[6] & 8 > 0))
	data = append(data, NewBoolDatum(
		"d:horyu2:switch_300_normal", timestamp, bytes[6] & 4 > 0))
	data = append(data, NewBoolDatum(
		"d:horyu2:debris_collision", timestamp, bytes[6] & 2 == 0))
	data = append(data, NewBoolDatum(
		"d:horyu2:reserve_command", timestamp, bytes[6] & 1 > 0))

	// Byte 7: More status bits
	data = append(data, NewBoolDatum(
		"d:horyu2:mission_mode", timestamp, bytes[7] & 64 > 0))
	data = append(data, NewBoolDatum(
		"d:horyu2:kill_switch_main", timestamp, bytes[7] & 32 == 0))
	data = append(data, NewBoolDatum(
		"d:horyu2:kill_switch_comm", timestamp, bytes[7] & 16 == 0))

	return data, nil
}

func DecodeHoryu2(decoded_morse string, timestamp int64)  (
	data []pb.TelemetryDatum, err error) {
	return decodeMorseEachWord(decoded_morse, timestamp, decodeHoryu2)
}