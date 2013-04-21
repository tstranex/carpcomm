// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"

func TestAausat3Morse(t *testing.T) {
	data, err := DecodeAausat3Morse("OZ3CUBB7.5T25", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 2 {
		t.Errorf("Expected data length 2, got %d", len(data))
		return
	}
	
	ExpectDoubleDatum(t, data[0], "d:aausat3:bat_v", 7.5)
	ExpectDoubleDatum(t, data[1], "d:aausat3:beacon_t", 298.15)
}

func TestAausat3MorseWithJunk(t *testing.T) {
	data, err := DecodeAausat3Morse(
		"  oz3cub B 7.5 T 09 \n\n oz33cub B 7.5 T 25 \n dfgkfjglkfg " +
		"OZ3CUBB75T25 oz3CUB b1.2 \t t2", 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 4 {
		t.Errorf("Expected data length 4, got %d", len(data))
		return
	}
	
	ExpectDoubleDatum(t, data[0], "d:aausat3:bat_v", 7.5)
	ExpectDoubleDatum(t, data[1], "d:aausat3:beacon_t", 282.15)
	ExpectDoubleDatum(t, data[2], "d:aausat3:bat_v", 1.2)
	ExpectDoubleDatum(t, data[3], "d:aausat3:beacon_t", 275.15)
}