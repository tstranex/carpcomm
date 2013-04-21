// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "testing"

func TestBytesToInt63(t *testing.T) {
	i := BytesToInt63([]byte{
		0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xee})
	if i != 0x7777777777777777 {
		t.Logf("i: %x", i)
		t.Fail()
	}
}