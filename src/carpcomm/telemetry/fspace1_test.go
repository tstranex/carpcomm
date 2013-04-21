// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"
import "encoding/hex"
import "time"

func TestFSpace1MorseDecode1(t *testing.T) {
	// 12345 = 00001 00010 00011 00100 00100
	//       = 00001000 10000110 01000010 0
	//       = 8 134 66
	data, err := DecodeFSpace1Morse("zzxv1vn12344", 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 3 {
		t.Errorf("Wrong number of datums: %d, expected 3", len(data))
		return
	}

	ExpectInt64Datum(t, data[0], "d:fspace1:reset_count", 8)
	ExpectDoubleDatum(t, data[1], "d:fspace1:obc_t", 307.15)
	ExpectDoubleDatum(t, data[2], "d:fspace1:out_y1_t", 239.15)
}

func TestFSpace1MorseDecode2(t *testing.T) {
	// ABcdE = 01010 01011 01100 01101 01111
	//       = 01010010 11011000 11010111 1
	//       = 82 216 215
	data, err := DecodeFSpace1Morse("zzxv1vnABcdF", 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 3 {
		t.Errorf("Wrong number of datums: %d, expected 3", len(data))
		return
	}

	ExpectInt64Datum(t, data[0], "d:fspace1:reset_count", 82)
	ExpectDoubleDatum(t, data[1], "d:fspace1:obc_t", 389.15)
	ExpectDoubleDatum(t, data[2], "d:fspace1:out_y1_t", 388.15)
}

// Example from:
// http://fspace.edu.vn/F_1_packet_for_radio_operators/GSS_UserManual.pdf
func TestFSpace1MorseDecode3(t *testing.T) {
	data, err := DecodeFSpace1Morse("zzXV1VN09FNQ", 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 3 {
		t.Errorf("Wrong number of datums: %d, expected 3", len(data))
		return
	}

	ExpectInt64Datum(t, data[0], "d:fspace1:reset_count", 2)
	ExpectDoubleDatum(t, data[1], "d:fspace1:obc_t", 268.15)
	ExpectDoubleDatum(t, data[2], "d:fspace1:out_y1_t", 298.15)
}

// Example from:
// http://fspace.edu.vn/F_1_packet_for_radio_operators/F-1_OBC1_PWM_CW_beacon.mp3
func TestFSpace1MorseDecode4(t *testing.T) {
	data, err := DecodeFSpace1Morse("ZZXV1VN1PV8R", 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 3 {
		t.Errorf("Wrong number of datums: %d, expected 3", len(data))
		return
	}

	ExpectInt64Datum(t, data[0], "d:fspace1:reset_count", 14)
	ExpectDoubleDatum(t, data[1], "d:fspace1:obc_t", 299.15)
	ExpectDoubleDatum(t, data[2], "d:fspace1:out_y1_t", 314.15)
}

func fspace1MorseFailIfNoError(t *testing.T, message string) {
	_, err := DecodeFSpace1Morse(message, 123)
	if err == nil {
		t.Errorf("Should be invalid: %s", message)
	}
}

func TestFSpace1MorseInvalid(t *testing.T) {
	fspace1MorseFailIfNoError(t, "zzJG6KBWHORYUzz")  // wrong callsign
	fspace1MorseFailIfNoError(t, "zzxv1vn1234zz")  // wrong length
	fspace1MorseFailIfNoError(t, "XV1VN09FNQ")  // missing z padding

	fspace1MorseFailIfNoError(t, "zzxv1vnABcdE")  // invalid checksum
	fspace1MorseFailIfNoError(t, "zzxv1vn12345")  // invalid checksum
	fspace1MorseFailIfNoError(t, "zzxv1vnnd3fs")  // invalid checksum
}



// Example from:
// http://fspace.edu.vn/F_1_packet_for_radio_operators/F-1_OBC2_KISS_beacon.mp3
func TestFSpace1FrameDecode(t *testing.T) {
	const hexFrame = "020000088000817E2888938E8C91908F8F"
	frame, _ := hex.DecodeString(hexFrame)
	data, err := DecodeFrame_fspace1(frame, 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 11 {
		t.Errorf("Wrong number of datums: %d, expected 11", len(data))
		return
	}

	ExpectTimestampDatum(t, data[0], "d:fspace1:clock",
		time.Date(2012, time.January, 1, 0, 0, 16, 0, time.UTC))
	ExpectDoubleDatum(t, data[1], "d:fspace1:bat_v", 3.82)
	ExpectDoubleDatum(t, data[2], "d:fspace1:cell_v", 4.0)
	ExpectDoubleDatum(t, data[3], "d:fspace1:out_y2_t", 36 + 273.15)
	ExpectDoubleDatum(t, data[4], "d:fspace1:out_y1_t", 47 + 273.15)
	ExpectDoubleDatum(t, data[5], "d:fspace1:out_x1_t", 42 + 273.15)
	ExpectDoubleDatum(t, data[6], "d:fspace1:out_z2_t", 40 + 273.15)
	ExpectDoubleDatum(t, data[7], "d:fspace1:out_z1_t", 45 + 273.15)
	ExpectDoubleDatum(t, data[8], "d:fspace1:out_x2_t", 44 + 273.15)
	ExpectDoubleDatum(t, data[9], "d:fspace1:in_z1_t", 43 + 273.15)
	ExpectDoubleDatum(t, data[10], "d:fspace1:radio_t", 43 + 273.15)
}