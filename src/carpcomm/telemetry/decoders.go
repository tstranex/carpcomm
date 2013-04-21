// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "carpcomm/pb"
import "carpcomm/util"
import "encoding/hex"
import "strings"

// Returns nil,nil if no decoder is available for the satellite.
func DecodeMorse(satellite_id, decoded_morse string, timestamp int64) (
	[]pb.TelemetryDatum, error) {
	switch satellite_id {
	case "swisscube": return DecodeSwisscube(decoded_morse, timestamp)
	case "horyu2": return DecodeHoryu2(decoded_morse, timestamp)
	case "masat1": return DecodeMasat1(decoded_morse, timestamp)
	case "fspace1": return DecodeFSpace1Morse(decoded_morse, timestamp)
	case "fitsat1": return DecodeFitsat1Morse(decoded_morse, timestamp)
	case "aausat3": return DecodeAausat3Morse(decoded_morse, timestamp)
	}
	return nil, nil
}

// Returns nil,nil if no decoder is available for the satellite.
func DecodeFrame(satellite_id string, frame []byte, timestamp int64) (
	[]pb.TelemetryDatum, error) {
	switch satellite_id {
	case "aeneas": return DecodeAeneas(frame, timestamp)
	case "csswe": return DecodeFrame_csswe(frame, timestamp)
	case "fspace1": return DecodeFrame_fspace1(frame, timestamp)
	case "hrbe": return DecodeFrame_hrbe(frame, timestamp)
	case "techedsat": return DecodeFrame_techedsat(frame, timestamp)
	case "strand1": return DecodeFrame_strand1(frame, timestamp)
	}
	return nil,nil
}

func decodeFreeform(satellite_id string, frame []byte, timestamp int64) (
	[]pb.TelemetryDatum, [][]byte) {

	stripped := util.StripWhitespace((string)(frame))

	data, _ := DecodeMorse(satellite_id, stripped, timestamp)
	if data != nil {
		return data, nil
	}

	ax25_header := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// Text frame.
	data, _ = DecodeFrame(satellite_id, frame, timestamp)
	if data != nil {
		return data, [][]byte{frame}
	}
	ax25_frame := append(ax25_header, frame...)
	data, _ = DecodeFrame(
		satellite_id, ax25_frame, timestamp)
	if data != nil {
		return data, [][]byte{ax25_frame}
	}

	// Hex-encoded frame.
	hex_frame, err := hex.DecodeString(stripped)
	if err == nil {
		data, _ = DecodeFrame(satellite_id, hex_frame, timestamp)
		if data != nil {
			return data, [][]byte{hex_frame}
		}

		ax25_hex_frame := append(ax25_header, hex_frame...)
		data, _ = DecodeFrame(satellite_id, ax25_hex_frame, timestamp)
		if data != nil {
			return data, [][]byte{ax25_hex_frame}
		}

		frames := DecodeKISS(hex_frame)
		for _, frame := range frames {
			d, _ := DecodeFrame(satellite_id, frame, timestamp)
			if d != nil {
				data = append(data, d...)
			}
		}
		if data != nil {
			return data, frames
		}
	}

	return nil, nil
}

// Decode data in an unknown format entered manually on the website.
// Various methods and encoding are tried until something works.
func DecodeFreeform(satellite_id string, frame []byte, timestamp int64) (
	r []pb.TelemetryDatum, frames [][]byte) {

	r, frames = decodeFreeform(satellite_id, frame, timestamp)
	if r != nil {
		return r, frames
	}

	// Try handling multi-line output e.g. from a serial dump program.

	// [9 Bytes unknown Protocol]
	// 1 > C0 10 01 E1 00 00 00 2E C0 
	// À..á....À
	//
	// [14 Bytes unknown Protocol]
	// 1 > C0 10 DB DC 80 30 06 01 E0 00 00 21 04 C0 
	// À.ÛÜ€0..à..!.À
	
	for _, line := range strings.Split((string)(frame), "\n") {
		i := strings.Index(line, ">")
		if i > 0 {
			line = line[i+1:]
		}
		d, f := decodeFreeform(
			satellite_id, []byte(line), timestamp)
		if d != nil {
			r = append(r, d...)
		}
		if f != nil {
			frames = append(frames, f...)
		}
	}

	return r, frames
}