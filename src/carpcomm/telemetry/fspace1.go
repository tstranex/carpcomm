// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "carpcomm/pb"
import "errors"
import "fmt"
import "strings"
import "time"

const fspace1Callsign = "XV1VN"

func fspace1MorseCharToInt(c byte) (byte, error) {
	if c >= '0' && c <= '9' {
		return c - '0', nil
	} else if c >= 'A' && c <= 'V' {
		return c - 'A' + 10, nil
	}
	return 0, errors.New(fmt.Sprintf("Unknown char: %c", c))
}

func fspace1Temperature(b byte) float64 {
	return float64((int32)(b)) - 100.0 + 273.15
}

// Documentation: http://fspace.edu.vn/?page_id=27
func decodeFSpace1Morse(morse_chars string, timestamp int64) (
	data []pb.TelemetryDatum, err error) {

	morse_chars = strings.ToUpper(morse_chars)

	if len(morse_chars) != 12 {
		return nil, errors.New(fmt.Sprintf(
			"Frame has wrong length %d, expected 12",
			len(morse_chars)))
	}

	// Ignore the z padding.
	morse_chars = morse_chars[2:]

	if morse_chars[:5] != fspace1Callsign {
		return nil, errors.New(fmt.Sprintf(
			"Wrong callsign: %s", morse_chars[:5]))
	}

	var bits uint32
	for _, c := range morse_chars[5:] {
		i, err := fspace1MorseCharToInt((byte)(c))
		if err != nil {
			return nil, err
		}
		bits = (bits << 5) | (uint32)(i)
	}

	byte1 := (byte)((bits >> 17) & 255)
	byte2 := (byte)((bits >> 9) & 255)
	byte3 := (byte)((bits >> 1) & 255)

	sum := (uint32)(byte1) + (uint32)(byte2) + (uint32)(byte3)
	if sum % 2 != (bits & 1) {
		return nil, errors.New("Invalid checksum")
	}

	data = append(data, NewInt64Datum(
		"d:fspace1:reset_count", timestamp, (int64)(byte1)))
	data = append(data, NewDoubleDatum(
		"d:fspace1:obc_t", timestamp, fspace1Temperature(byte2)))
	data = append(data, NewDoubleDatum(
		"d:fspace1:out_y1_t", timestamp, fspace1Temperature(byte3)))

	return data, nil
}

func DecodeFSpace1Morse(morse_chars string, timestamp int64) (
	data []pb.TelemetryDatum, err error) {
	return decodeMorseEachWord(morse_chars, timestamp, decodeFSpace1Morse)
}

// Documentation: http://fspace.edu.vn/?page_id=27
func DecodeFrame_fspace1(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {

	// The spec says the frame should be 14 bytes but the sample packet has
	// 17 bytes.
	if len(frame) != 17 {
		return nil, errors.New(fmt.Sprintf(
			"Frame too short: %d, expected %d.", len(frame), 14))
	}

	// Ignore the first three bytes. Not sure what they're for.
	frame = frame[3:]

	var bits uint32
	for i := 0; i < 4; i++ {
		bits = (bits << 8) | (uint32)(frame[i])
	}
	bits = bits >> 3

	second := (int)(bits & 63)
	minute := (int)((bits >> 6) & 63)
	hour := (int)((bits >> 12) & 31)
	year := (int)(2012 + (bits >> 17) & 7)
	month := time.January + (time.Month)((bits >> 20) & 15) - 1
	day := (int)((bits >> 24) & 31)

	date := time.Date(year, month, day, hour, minute, second, 0, time.UTC)
	data = append(data, NewTimestampDatum(
		"d:fspace1:clock", timestamp, date))

	v := (((uint32)(frame[3]) << 8) | (uint32)(frame[4])) & 2047
	data = append(data, NewDoubleDatum(
		"d:fspace1:bat_v", timestamp, (float64)(v) * 0.01))
	data = append(data, NewDoubleDatum(
		"d:fspace1:cell_v", timestamp, (float64)(frame[5]) * 0.1))

	temperature_keys := []string{
		"d:fspace1:out_y2_t",
		"d:fspace1:out_y1_t",
		"d:fspace1:out_x1_t",
		"d:fspace1:out_z2_t",
		"d:fspace1:out_z1_t",
		"d:fspace1:out_x2_t",
		"d:fspace1:in_z1_t",
		"d:fspace1:radio_t"}
	for i, key := range temperature_keys {
		data = append(data, NewDoubleDatum(
			key, timestamp, fspace1Temperature(frame[6+i])))
	}

	return data, nil
}