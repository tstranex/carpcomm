// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"

func TestFitsat1Callsign(t *testing.T) {
	data, err := DecodeFitsat1Morse("hideniwakajapan", 123)
	if err != nil {
		t.Error(err)
	}
	if data != nil {
		t.Error("Expected nil data for callsign.")
	}
}

func TestFitsat1S1(t *testing.T) {
	data, err := DecodeFitsat1Morse("s13b2b1b0b", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 4 {
		t.Errorf("Expected data length 4, got %d", len(data))
		return
	}
	
	ExpectDoubleDatum(t, data[0], "d:fitsat1:rssi_437", 1.15234)
	ExpectDoubleDatum(t, data[1], "d:fitsat1:total_cell_v", 0.83984)
	ExpectDoubleDatum(t, data[2], "d:fitsat1:total_cell_c", 0.21093)
	ExpectDoubleDatum(t, data[3], "d:fitsat1:bat1_v", 0.21484)
}

func TestFitsat1S2(t *testing.T) {
	data, err := DecodeFitsat1Morse("s23b2b1b0b", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 4 {
		t.Errorf("Expected data length 4, got %d", len(data))
		return
	}
	
	ExpectDoubleDatum(t, data[0], "d:fitsat1:bat1_c", -0.5390625)
	ExpectDoubleDatum(t, data[1], "d:fitsat1:bat3_v", 0.83984)
	ExpectDoubleDatum(t, data[2], "d:fitsat1:bat3_c", -0.197265625)
	ExpectDoubleDatum(t, data[3], "d:fitsat1:2p5_bus_v", 0.21484)
}

func TestFitsat1S3(t *testing.T) {
	data, err := DecodeFitsat1Morse("s33b2b1b0b", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 4 {
		t.Errorf("Expected data length 4, got %d", len(data))
		return
	}
	
	ExpectDoubleDatum(t, data[0], "d:fitsat1:cell_x2_v", 2.07421875)
	ExpectDoubleDatum(t, data[1], "d:fitsat1:cell_y2_v", 1.51171875)
	ExpectDoubleDatum(t, data[2], "d:fitsat1:cell_x1_v", 0.94921875)
	ExpectDoubleDatum(t, data[3], "d:fitsat1:cell_y1_v", 0.38671875)
}

func TestFitsat1S4(t *testing.T) {
	data, err := DecodeFitsat1Morse("s43b2b1b0b", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 4 {
		t.Errorf("Expected data length 4, got %d", len(data))
		return
	}
	
	ExpectDoubleDatum(t, data[0], "d:fitsat1:bat3_t", 326.8609374)
	ExpectDoubleDatum(t, data[1], "d:fitsat1:bat1_t", 298.7359374)
	ExpectDoubleDatum(t, data[2], "d:fitsat1:panel_z2_t", 270.6109374)
	ExpectDoubleDatum(t, data[3], "d:fitsat1:panel_z1_t", 242.4859374)
}

func TestFitsat1S5(t *testing.T) {
	data, err := DecodeFitsat1Morse("s53b2b1b0b", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 2 {
		t.Errorf("Expected data length 2, got %d", len(data))
		return
	}
	
	ExpectDoubleDatum(t, data[0], "d:fitsat1:rssi_1200", 1.037109375)
	ExpectDoubleDatum(t, data[1], "d:fitsat1:time_since_boot_s", 2824971.0)
}

func TestFitsat1Multiple(t *testing.T) {
	data, err := DecodeFitsat1Morse("S1 F0 03 00 BA S2 88 D5 81 81", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 8 {
		t.Errorf("Expected data length 8, got %d", len(data))
		return
	}

	ExpectDoubleDatum(t, data[0], "d:fitsat1:rssi_437", 4.687500)
	ExpectDoubleDatum(t, data[1], "d:fitsat1:total_cell_v", 0.058594)
	ExpectDoubleDatum(t, data[2], "d:fitsat1:total_cell_c", 0.000000)
	ExpectDoubleDatum(t, data[3], "d:fitsat1:bat1_v", 3.632812)

	ExpectDoubleDatum(t, data[4], "d:fitsat1:bat1_c", 0.062500)
	ExpectDoubleDatum(t, data[5], "d:fitsat1:bat3_v", 4.160156)
	ExpectDoubleDatum(t, data[6], "d:fitsat1:bat3_c", 0.195312)
	ExpectDoubleDatum(t, data[7], "d:fitsat1:2p5_bus_v", 2.519531)
}

func TestFitsat1MultipleWithGarbageBehind(t *testing.T) {
	data, err := DecodeFitsat1Morse("S1 F0 02 00 BA S289 DA 81 81 01 00 01", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 8 {
		t.Errorf("Expected data length 8, got %d", len(data))
		return
	}

	ExpectDoubleDatum(t, data[0], "d:fitsat1:rssi_437", 4.687500)
	ExpectDoubleDatum(t, data[1], "d:fitsat1:total_cell_v", 0.039062)
	ExpectDoubleDatum(t, data[2], "d:fitsat1:total_cell_c", 0.000000)
	ExpectDoubleDatum(t, data[3], "d:fitsat1:bat1_v", 3.632812)

	ExpectDoubleDatum(t, data[4], "d:fitsat1:bat1_c", 0.070312)
	ExpectDoubleDatum(t, data[5], "d:fitsat1:bat3_v", 4.257812)
	ExpectDoubleDatum(t, data[6], "d:fitsat1:bat3_c", 0.195312)
	ExpectDoubleDatum(t, data[7], "d:fitsat1:2p5_bus_v", 2.519531)
}

func TestFitsat1MultipleWithGarbageInFront(t *testing.T) {
	data, err := DecodeFitsat1Morse("F0 02 01 BA S289DA 81 81 S3 01 00 01 00 S41E 1F1A S51B 00 51 F0", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 10 {
		t.Errorf("Expected data length 10, got %d", len(data))
		return
	}

	ExpectDoubleDatum(t, data[0], "d:fitsat1:bat1_c", 0.070312)
	ExpectDoubleDatum(t, data[1], "d:fitsat1:bat3_v", 4.257812)
	ExpectDoubleDatum(t, data[2], "d:fitsat1:bat3_c", 0.195312)
	ExpectDoubleDatum(t, data[3], "d:fitsat1:2p5_bus_v", 2.519531)

	ExpectDoubleDatum(t, data[4], "d:fitsat1:cell_x2_v", 0.035156)
	ExpectDoubleDatum(t, data[5], "d:fitsat1:cell_y2_v", 0.0)
	ExpectDoubleDatum(t, data[6], "d:fitsat1:cell_x1_v", 0.035156)
	ExpectDoubleDatum(t, data[7], "d:fitsat1:cell_y1_v", 0.0)

	ExpectDoubleDatum(t, data[8], "d:fitsat1:rssi_1200", 0.474609)
	ExpectDoubleDatum(t, data[9], "d:fitsat1:time_since_boot_s", 20976.0)
}

func fitsat1FailIfNoError(t *testing.T, message string) {
	_, err := DecodeFitsat1Morse(message, 123)
	if err == nil {
		t.Errorf("Should be invalid: %s", message)
	}
}

func TestFitsat1Invalid(t *testing.T) {
	fitsat1FailIfNoError(t, "niwaka")
	fitsat1FailIfNoError(t, "S600000000")
	fitsat1FailIfNoError(t, "S10000000")
	fitsat1FailIfNoError(t, "S1000x0000")
}
