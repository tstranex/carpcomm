// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2012 Carpcomm GmbH

// Decoder for Hiscock Radiation Belt Explorer (hrbe).
// This is incomplete.

package telemetry

import "errors"
import "carpcomm/pb"

const hrbeCallsign = "K7MSU-1"

func DecodeFrame_hrbe(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {

	// First do some quick validation.
	if len(frame) < 81 {
		return nil, errors.New("Frame is too short")
	}
	if string(frame[4:11]) != hrbeCallsign {
		return nil, errors.New("Frame has wrong callsign")
	}

	// Battery current.
	InstCurrent := func(f []byte) float64 {
		c := (int16(f[0]) << 8) | int16(f[1])
		return 0.0390625 * float64(c) * 1e-3
	}
	data = append(data, NewDoubleDatum(
		"d:hrbe:bat1_c", timestamp, InstCurrent(frame[44:46])))
	data = append(data, NewDoubleDatum(
		"d:hrbe:bat2_c", timestamp, InstCurrent(frame[54:56])))

	// TODO: Decode further telemetry values here.

	return data, nil
}