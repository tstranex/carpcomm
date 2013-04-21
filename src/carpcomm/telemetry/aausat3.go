// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "errors"
import "strings"
import "carpcomm/pb"
import "carpcomm/util"
import "strconv"

const kAausat3Call = "OZ3CUB"
const kAausat3MinLength = 12
const kAausat3MaxLength = 13

// Example: OZ3CUBB7.5T25
func aausat3Decode(s string, timestamp int64) (
	data []pb.TelemetryDatum, err error) {
	if len(s) < kAausat3MaxLength {
		return nil, errors.New("Too short")
	}

	if s[:len(kAausat3Call)] != kAausat3Call {
		return nil, errors.New("Wrong callsign")
	}
	s = s[len(kAausat3Call):]

	if s[0] != 'B' || s[2] != '.' || s[4] != 'T' {
		return nil, errors.New("Invalid format")
	}

	v1, err := strconv.Atoi(s[1:2])
	if err != nil {
		return nil, err
	}
	v2, err := strconv.Atoi(s[3:4])
	if err != nil {
		return nil, err
	}
	data = append(data, NewDoubleDatum(
		"d:aausat3:bat_v", timestamp, float64(v1) + 0.1*float64(v2)))

	// Temperature can be either 1 or 2 digits.
	t, err := strconv.Atoi(s[5:7])
	if err != nil {
		t, err = strconv.Atoi(s[5:6])
		if err != nil {
			return nil, err
		}
	}
	data = append(data, NewDoubleDatum(
		"d:aausat3:beacon_t", timestamp, float64(t) + 273.15))

	return data, nil
}

func DecodeAausat3Morse(decoded_morse string, timestamp int64)  (
	data []pb.TelemetryDatum, err error) {
	decoded_morse = strings.ToUpper(util.StripWhitespace(decoded_morse)) +
		" ";

	for i := 0; i <= len(decoded_morse) - kAausat3MaxLength; i++ {
		d, _ := aausat3Decode(decoded_morse[i:], timestamp)
		data = append(data, d...)
		if len(d) > 0 {
			i += kAausat3MinLength
		}
	}

	if len(data) == 0 {
		return nil, errors.New("No valid messages found")
	}

	return data, nil
}