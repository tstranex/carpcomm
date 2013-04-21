// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"

func TestDecodeKISS(t *testing.T) {
	frames := DecodeKISS([]byte("\x00test1\xc0\x00test2\xc0"))
	if len(frames) != 2 {
		t.Errorf("Wrong number of frames: %d", len(frames))
		return
	}
	if string(frames[0]) != "test1" {
		t.Errorf("Frame 0 wrong")
	}
	if string(frames[1]) != "test2" {
		t.Errorf("Frame 1 wrong")
	}
}