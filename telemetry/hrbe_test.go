// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2012 Carpcomm GmbH

// Tests for hrbe.go.

package telemetry

import "testing"
import "encoding/hex"

const hrbeTestFrame = "337205024b374d53552d312009200800841dc63d3c0020fecd290c31309090bb000ff08f11d7da5cfd005b80fe3e0453f62504805ba0fc0709060607070352210267026ed3d45055d5a000000000000000bf0258353535353535000000000000000000000000000000000000000000000001000100000000000000000000000000000000000000000000000000000001000000010000"

func TestDecodeFrame_hrbe(t *testing.T) {
	frame, _ := hex.DecodeString(hrbeTestFrame)

	data, err := DecodeFrame_hrbe(frame, 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 2 {
		t.Errorf("Wrong number of datums: %d", len(data))
		return
	}

	ExpectDoubleDatum(t, data[0], "d:hrbe:bat1_c", -0.01758)
	ExpectDoubleDatum(t, data[1], "d:hrbe:bat2_c", -0.03973)
}