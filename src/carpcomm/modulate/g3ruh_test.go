// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "testing"

func TestG3RUHScramble(t *testing.T) {
	original := "11010010 10111101 01001111 00101010 101"
	bits := stringToBits(original)
	G3RUHScramble(bits)
	G3RUHDescramble(bits)
	if bitstring(bits) != original {
		t.Errorf("G3RUHScrable/G3RUHDescramble didn't produce " +
			"original bits for %s", original)
	}

	// A short payload.
	short := "000"
	bits = stringToBits(short)
	G3RUHScramble(bits)
	G3RUHDescramble(bits)
	if bitstring(bits) != short {
		t.Errorf("G3RUHScrable/G3RUHDescramble didn't produce " +
			"original bits for %s", short)
	}
}