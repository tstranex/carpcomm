// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "testing"

func expectChecksumCRC16HDLC(t *testing.T, data []byte, expected uint16) {
	actual := ChecksumCRC16HDLC(data)
	if actual != expected {
		t.Errorf("ChecksumCRC16HDLC returned %04x, expected %04x, " +
			"for data: %v", actual, expected, data)
	}
}

func TestChecksumCRC16HDLC(t *testing.T) {
	expectChecksumCRC16HDLC(t, []byte(""), 0x0000)
	expectChecksumCRC16HDLC(t, []byte("A"), 0xa3f5)
	expectChecksumCRC16HDLC(t, []byte("AB"), 0x31ef)
	expectChecksumCRC16HDLC(t, []byte("123456789"), 0x906e)
}

func bitstring(bits []bool) string {
	s := ""
	for i, b := range bits {
		if i > 0 && i % 8 == 0 {
			s += " "
		}
		if b {
			s += "1"
		} else {
			s += "0"
		}
	}
	return s
}

func stringToBits(s string) (r []bool) {
	r = make([]bool, 0)
	for _, c := range s {
		if c != ' ' {
			r = append(r, c == '1')
		}
	}
	return r
}

func expectEncodeHDLC(t *testing.T, data []byte, expected string) {
	actual := bitstring(EncodeHDLC(data))
	if actual != expected {
		t.Errorf("EncodeHDLC returned\n%s\n, expected\n%s\n, " +
			"for data %v", actual, expected, data)
	}
}

func TestEncodeHDLC(t *testing.T) {
	expectEncodeHDLC(t, []byte(""),
		"01111110 00000000 00000000 01111110")
	expectEncodeHDLC(t, []byte("A"),
		"01111110 10000010 10101111 10100010 10111111 0")
	expectEncodeHDLC(t, []byte{0xff},
		"01111110 11111011 10000000 01111101 11011111 10")
}

func expectDecodeHDLC(t *testing.T, bits string, expected [][]byte) {
	actual := DecodeHDLC(stringToBits(bits))
	if len(actual) != len(expected) {
		t.Errorf("DecodeHDLC returned %d packets, %d expected, " +
			"for data %v", len(actual), len(expected), bits)
		return
	}
	for i := 0; i < len(actual); i++ {
		ok := true
		for j := 0; j < len(actual[i]); j++ {
			if actual[i][j] != expected[i][j] {
				ok = false
				break
			}
		}
		if !ok {
			t.Errorf("DecodeHDLC packet %d differs from expected " +
			"for data %v", i, bits)
		}
	}
}

func TestDecodeHDLC(t *testing.T) {
	expectDecodeHDLC(t, "11001" + bitstring(EncodeHDLC([]byte{})) + "111",
		[][]byte{[]byte{}})
	expectDecodeHDLC(t, "01111110" + bitstring(EncodeHDLC([]byte("A"))),
		[][]byte{[]byte("A")})
	expectDecodeHDLC(t, bitstring(EncodeHDLC([]byte{0xff})) + "01111110",
		[][]byte{[]byte{0xff}})

	// Multiple packets.
	expectDecodeHDLC(t,
		"1111" + bitstring(EncodeHDLC([]byte("hello"))) +
		"0101011110" + bitstring(EncodeHDLC([]byte("world"))) +
		"0010",
		[][]byte{[]byte("hello"), []byte("world")})

	// Invalid checksum.
	expectDecodeHDLC(t, "01111110 10000100 10101111 10100010 10111111 0",
		nil)
}