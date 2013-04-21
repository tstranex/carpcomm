// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "strconv"
import "errors"
import "strings"
import "carpcomm/pb"

const kMasat1Callsign = "HA5MASAT"

// Documentation: http://cubesat.bme.hu/en/radioamatoroknek/
func decodeMasat1(decoded_morse string, timestamp int64)  (
	data []pb.TelemetryDatum, err error) {

	decoded_morse = strings.ToUpper(decoded_morse)

	// Callsign
	if decoded_morse == kMasat1Callsign {
		return nil, nil
	}

	if !(len(decoded_morse) == 2 || len(decoded_morse) == 3) {
		return nil, errors.New(
			"Message length should be either 2 or 3 chars.")
	}

	i, err := strconv.Atoi(decoded_morse)
	if err != nil {
		return nil, errors.New(
			"Unable to convert message to an integer.")
	}

	// Battery voltage
	if len(decoded_morse) == 3 {
		data = append(data, NewDoubleDatum(
			"d:masat1:bat_v", timestamp, float64(i) / 100))
	}

	// Battery temperature
	if len(decoded_morse) == 2 {
		data = append(data, NewDoubleDatum(
			"d:masat1:bat_t", timestamp, float64(i) + 273.15))
	}

	return data, nil
}

func DecodeMasat1(decoded_morse string, timestamp int64) (
	data []pb.TelemetryDatum, err error) {
	return decodeMorseEachWord(decoded_morse, timestamp, decodeMasat1)
}