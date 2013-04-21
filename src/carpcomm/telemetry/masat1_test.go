// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"

func TestMasat1Callsign(t *testing.T) {
	data, err := DecodeMasat1("ha5MaSAT", 123)
	if err != nil {
		t.Error(err)
	}
	if data != nil {
		t.Error("Expected nil data for callsign.")
	}
}

func TestMasat1DecodeBatteryVoltage(t *testing.T) {
	data, err := DecodeMasat1("406", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Errorf("Wrong number of datums: %d, expected 1", len(data))
		return
	}
	ExpectDoubleDatum(t, data[0], "d:masat1:bat_v", 4.06)
}

func TestMasat1DecodeBatteryTemperature(t *testing.T) {
	data, err := DecodeMasat1("05", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Errorf("Wrong number of datums: %d, expected 1", len(data))
		return
	}
	ExpectDoubleDatum(t, data[0], "d:masat1:bat_t", 278.15)
}

func TestMasat1Multiple(t *testing.T) {
	data, err := DecodeMasat1("ha5MaSAT 406 05", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 2 {
		t.Errorf("Wrong number of datums: %d, expected 2", len(data))
		return
	}
	ExpectDoubleDatum(t, data[0], "d:masat1:bat_v", 4.06)
	ExpectDoubleDatum(t, data[1], "d:masat1:bat_t", 278.15)
}

func masat1FailIfNoError(t *testing.T, message string) {
	_, err := DecodeMasat1(message, 123)
	if err == nil {
		t.Errorf("Should be invalid: %s", message)
	}
}

func TestMasat1Invalid(t *testing.T) {
	masat1FailIfNoError(t, "HA6MASAT")  // Wrong char 6.
	masat1FailIfNoError(t, "00B")  // Invalid decimal digit B.
	masat1FailIfNoError(t, "x1")  // Invalid decimal digit x.
	masat1FailIfNoError(t, "1234")  // Too long.
	masat1FailIfNoError(t, "1")  // Too short.
}
