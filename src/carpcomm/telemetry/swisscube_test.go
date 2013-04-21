// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"

func TestDecodePart0(t *testing.T) {
	data, err := DecodeSwisscube("hB9Eg/1", 123)
	if err != nil {
		t.Error(err)
	}
	if data != nil {
		t.Error("Expected nil data for callsign.")
	}
}

func TestDecodePart1(t *testing.T) {
	data, err := DecodeSwisscube("ATVV", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 6 {
		t.Errorf("Too few datums: %d, expected 6", len(data))
		return
	}
	ExpectBoolDatum(t, data[0], "d:swisscube:ads_power", false)
	ExpectBoolDatum(t, data[1], "d:swisscube:payload_power", true)
	ExpectBoolDatum(t, data[2], "d:swisscube:adcs_power", true)
	ExpectBoolDatum(t, data[3], "d:swisscube:cdms_power", false)
	ExpectBoolDatum(t, data[4], "d:swisscube:beacon_power", true)
	ExpectBoolDatum(t, data[5], "d:swisscube:com_power", true)
}

func TestDecodePart2(t *testing.T) {
	data, err := DecodeSwisscube("uAUV4e6", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 2 {
		t.Errorf("Too few datums: %d, expected 2", len(data))
		return
	}
	ExpectDoubleDatum(t, data[0], "d:swisscube:bat1_v", 1.621490)
	ExpectDoubleDatum(t, data[1], "d:swisscube:bat2_v", 5.899878)
}

func TestDecodePart3(t *testing.T) {
	data, err := DecodeSwisscube("VUTVtBTVB", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 7 {
		t.Errorf("Too few datums: %d, expected 7", len(data))
	}

	ExpectIntervalDatum(t, data[0],
		"d:swisscube:cell_x1_c", 0.250, 0.375)
	ExpectIntervalDatum(t, data[1],
		"d:swisscube:cell_x2_c", 0.000, 0.125)
	ExpectIntervalDatum(t, data[2],
		"d:swisscube:cell_y1_c", 0.375, 0.500)
	ExpectIntervalDatum(t, data[3],
		"d:swisscube:cell_y2_c", 0.000, 0.125)
	ExpectIntervalDatum(t, data[4],
		"d:swisscube:cell_z1_c", 0.875, 1.000)
	ExpectIntervalDatum(t, data[5],
		"d:swisscube:cell_z2_c", 0.000, 0.125)

	ExpectDoubleDatum(t, data[6],
		"d:swisscube:bat1_t", 269.15)
}

func TestDecodeMultiple(t *testing.T) {
	data, err := DecodeSwisscube("HB9EG/1 ATVV UAUV4E6", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 8 {
		t.Errorf("Too few datums: %d, expected 8", len(data))
		return
	}

	ExpectBoolDatum(t, data[0], "d:swisscube:ads_power", false)
	ExpectBoolDatum(t, data[1], "d:swisscube:payload_power", true)
	ExpectBoolDatum(t, data[2], "d:swisscube:adcs_power", true)
	ExpectBoolDatum(t, data[3], "d:swisscube:cdms_power", false)
	ExpectBoolDatum(t, data[4], "d:swisscube:beacon_power", true)
	ExpectBoolDatum(t, data[5], "d:swisscube:com_power", true)

	ExpectDoubleDatum(t, data[6], "d:swisscube:bat1_v", 1.621490)
	ExpectDoubleDatum(t, data[7], "d:swisscube:bat2_v", 5.899878)
}

func failIfNoError(t *testing.T, message string) {
	_, err := DecodeSwisscube(message, 123)
	if err == nil {
		t.Errorf("Should be invalid: %s", message)
	}
}

func TestInvalid(t *testing.T) {
	failIfNoError(t, "hB9Eg 1")
	failIfNoError(t, "AUTUVA")
	failIfNoError(t, "UAUV4E")
	failIfNoError(t, "VUTVTBTVX")
}

