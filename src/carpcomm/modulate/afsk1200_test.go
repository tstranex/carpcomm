// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "testing"

func TestAFSK1200(t *testing.T) {
	payload := "Hello world!"
	rate := 22050.0
	samples := ModulateAFSK1200(EncodeHDLC([]byte(payload)), rate)
	packets := DemodulateAFSK1200(samples, rate)

	if len(packets) != 1 {
		t.Errorf("Wrong number of packets returned: %d", len(packets))
		return
	}

	if string(packets[0]) != payload {
		t.Errorf("Incorrect data received: '%s', expected '%s'",
			string(packets[0]), payload)
	}
}