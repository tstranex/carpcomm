// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"

func TestHoryuCallsign(t *testing.T) {
	data, err := DecodeHoryu2("jG6yBWHoRYu", 123)
	if err != nil {
		t.Error(err)
	}
	if data != nil {
		t.Error("Expected nil data for callsign.")
	}
}

func TestHoryuDecode(t *testing.T) {
	data, err := DecodeHoryu2("JG6YBWHORYU d1b8bab987d4fe3", 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 16 {
		t.Errorf("Wrong number of datums: %d, expected 16", len(data))
		return
	}

	ExpectDoubleDatum(t, data[0], "d:horyu2:bat_top_t", 284.653103)
	ExpectDoubleDatum(t, data[1], "d:horyu2:bat_bottom_t", 285.175455)
	ExpectDoubleDatum(t, data[2], "d:horyu2:comm_t", 280.035825)
	ExpectDoubleDatum(t, data[3], "d:horyu2:bat_c", 0.138360)
	ExpectDoubleDatum(t, data[4], "d:horyu2:bat_v", 4.201240)

	ExpectBoolDatum(t, data[5], "d:horyu2:clock_normal", true)
	ExpectBoolDatum(t, data[6], "d:horyu2:flash_main_normal", true)
	ExpectBoolDatum(t, data[7], "d:horyu2:flash_share_normal", true)
	ExpectBoolDatum(t, data[8], "d:horyu2:flash_300_normal", true)
	ExpectBoolDatum(t, data[9], "d:horyu2:switch_share_normal", true)
	ExpectBoolDatum(t, data[10], "d:horyu2:switch_300_normal", true)
	ExpectBoolDatum(t, data[11], "d:horyu2:debris_collision", false)
	ExpectBoolDatum(t, data[12], "d:horyu2:reserve_command", false)
	ExpectBoolDatum(t, data[13], "d:horyu2:mission_mode", false)
	ExpectBoolDatum(t, data[14], "d:horyu2:kill_switch_main", false)
	ExpectBoolDatum(t, data[15], "d:horyu2:kill_switch_comm", false)
}

func horyuFailIfNoError(t *testing.T, message string) {
	_, err := DecodeHoryu2(message, 123)
	if err == nil {
		t.Errorf("Should be invalid: %s", message)
	}
}

func TestHoryuInvalid(t *testing.T) {
	horyuFailIfNoError(t, "JG6KBWHORYU")  // Wrong char K.
	horyuFailIfNoError(t, "d1b8bab987d4feg")  // Invalid hex digit g.
	horyuFailIfNoError(t, "d1b8bab987d4fe33")  // Too long.
	horyuFailIfNoError(t, "d1b8bab987d4fe")  // Too short.
}
