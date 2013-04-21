// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "errors"
import "strings"
import "fmt"
import "carpcomm/pb"

const kSwisscubeCallsign = "HB9EG/1"

var swisscubeEncoding = 
	map[byte]int{
	'T': 0,
	'A': 1,
	'U': 2,
	'V': 3,
	'4': 4,
	'E': 5,
	'6': 6,
	'B': 7,
	'D': 8,
	'N': 9}

// Unit: Amps
var solarCellCurrent = [][2]float64{
	0: {0.000, 0.125},
	1: {0.125, 0.250},
	2: {0.250, 0.375},
	3: {0.375, 0.500},
	4: {0.500, 0.625},
	5: {0.625, 0.750},
	6: {0.750, 0.875},
	7: {0.875, 1.000},
}

// Documentation:
// http://swisscube-live.ch/Publish/S3-D-COM-1-3-Beacon_format.pdf
func decodeSwisscube(decoded_morse string, timestamp int64) (
	data []pb.TelemetryDatum, err error) {
	decoded_morse = strings.ToUpper(decoded_morse)

	// Part 0: callsign
	if decoded_morse == kSwisscubeCallsign {
		return nil, nil
	}

	// Extract octal digits.
	digits := make([]int, len(decoded_morse))
	for i := 0; i < len(decoded_morse); i++ {
		digit, ok := swisscubeEncoding[decoded_morse[i]]
		if !ok {
			return nil, errors.New(fmt.Sprintf(
				"Invalid char in position %d: %c (%s)",
				i, decoded_morse[i], decoded_morse))
		}
		digits[i] = digit
	}
	if len(digits) == 0 {
		return nil, errors.New("Too few chars")
	}

	// Part 1
	// Note that the beacon specficiation and reality differ for this part.
	if len(digits) == 4 && digits[0] == 1 {
		// Ignore the first digit.
		// It may be related to the error flag, but seems to be always
		// zero.

		// Power flag
		data = append(data, NewBoolDatum(
			"d:swisscube:ads_power", timestamp,
			digits[2] & 4 > 0))
		data = append(data, NewBoolDatum(
			"d:swisscube:payload_power", timestamp,
			digits[2] & 2 > 0))
		data = append(data, NewBoolDatum(
			"d:swisscube:adcs_power", timestamp,
			digits[2] & 1 > 0))
		data = append(data, NewBoolDatum(
			"d:swisscube:cdms_power", timestamp,
			digits[3] & 4 > 0))
		data = append(data, NewBoolDatum(
			"d:swisscube:beacon_power", timestamp,
			digits[3] & 2 > 0))
		data = append(data, NewBoolDatum(
			"d:swisscube:com_power", timestamp,
			digits[3] & 1 > 0))

		return data, nil
	}

	// Part 2
	if len(digits) == 7 && digits[0] == 2 {
		x1 := (float64)(digits[1]<<6 + digits[2]<<3 + digits[3])
		data = append(data, NewDoubleDatum(
			"d:swisscube:bat1_v", timestamp, 80*x1 / 4095))
		x2 := (float64)(digits[4]<<6 + digits[5]<<3 + digits[6])
		data = append(data, NewDoubleDatum(
			"d:swisscube:bat2_v", timestamp, 80*x2 / 4095))
		return data, nil
	}

	// Part 3
	if len(digits) == 9 && digits[0] == 3 {
		data = append(data, NewIntervalDatum(
			"d:swisscube:cell_x1_c", timestamp,
			solarCellCurrent[digits[1]]))
		data = append(data, NewIntervalDatum(
			"d:swisscube:cell_x2_c", timestamp,
			solarCellCurrent[digits[2]]))
		data = append(data, NewIntervalDatum(
			"d:swisscube:cell_y1_c", timestamp,
			solarCellCurrent[digits[3]]))
		data = append(data, NewIntervalDatum(
			"d:swisscube:cell_y2_c", timestamp,
			solarCellCurrent[digits[4]]))
		data = append(data, NewIntervalDatum(
			"d:swisscube:cell_z1_c", timestamp,
			solarCellCurrent[digits[5]]))
		data = append(data, NewIntervalDatum(
			"d:swisscube:cell_z2_c", timestamp,
			solarCellCurrent[digits[6]]))

		x := (float64)(digits[7]<<3 + digits[8])
		T := 4.0*x - 128.0 + 273.15 // K
		data = append(data, NewDoubleDatum(
			"d:swisscube:bat1_t", timestamp, T))

		return data, nil
	}

	return nil, errors.New("Invalid message")
}

func DecodeSwisscube(decoded_morse string, timestamp int64) (
	data []pb.TelemetryDatum, err error) {
	return decodeMorseEachWord(decoded_morse, timestamp, decodeSwisscube)
}