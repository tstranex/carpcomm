// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"

func TestDecodeFreeformMorse(t *testing.T) {
	data, _ := DecodeFreeform("fitsat1", ([]byte)(" s1 3b2\n b1b0b"), 0)
	if len(data) == 0 {
		t.Errorf("Valid morse decode failed.")
	}
}

func TestDecodeFreeformAX25PayloadText(t *testing.T) {
	data, _ := DecodeFreeform("techedsat", ([]byte)("ncasst.org0000057fcadcadcadcad85d85d85d85dbf4bf6bf5bf585d85d85d85dcc3cc3cc3cc332834a33933482d82d82d82d03000000170000012d2b"), 0)
	if len(data) == 0 {
		t.Errorf("Valid text AX.25 payload decode failed.")
	}
}

func TestDecodeFreeformAX25PayloadHex(t *testing.T) {
	data, _ := DecodeFreeform("techedsat", ([]byte)(" 6e 63617373742e6f72673030303030353766636164636164636164636164383564383564383564383564626634626636626635626635383     56\r\n\r\n43835643835643\n8356463633363633363633363633333323833346133333933333438326438326438326438326430333030303030\n30313730303030303132643262  "), 0)
	if len(data) == 0 {
		t.Errorf("Valid hex AX.25 payload decode failed.")
	}
}

func TestDecodeFreeformFrameText(t *testing.T) {
	data, _ := DecodeFreeform("techedsat", ([]byte)("0000000000000000ncasst.org0000057fcadcadcadcad85d85d85d85dbf4bf6bf5bf585d85d85d85dcc3cc3cc3cc332834a33933482d82d82d82d03000000170000012d2b"), 0)
	if len(data) == 0 {
		t.Errorf("Valid text frame decode failed.")
	}
}

func TestDecodeFreeformFrameHex(t *testing.T) {
	data, _ := DecodeFreeform("techedsat", ([]byte)("\n000000000000000000000000000000006E 63617373742e\n6f72673030303030353766636164636164636164636164383564383564383564383 5646266346266366 2663562663538\r35643835643835643835646363336363336363336363333332383334613333393333343832643832  643832643832643033303030303030313730303030303132643262\n\n"), 0)
	if len(data) == 0 {
		t.Errorf("Valid hex frame decode failed.")
	}
}

func TestDecodeFreeformMultiline(t *testing.T) {
	input := `
C0 80 D2 06 02 2D 01 02 04 01 
C0 80 A9 06 02 2D 04 02 CC 03 
C0 80 A1 09 02 66 A9 05 00 00 00 00 00 
C0 80 B4 09 02 66 9F 05 00 00 59 00 5A 
C0 80 A8 06 02 2D 01 02 FF 03 
C0 80 D4 06 02 2D 0A 02 EC 02 
C0 80 DB 09 02 66 90 05 00 00 59 00 59 
C0 80 DD 09 02 66 9A 05 00 00 57 00 58 
C0 80 9E 09 02 66 9A 05 00 00 59 00 5A 
`
	data, _ := DecodeFreeform("strand1", []byte(input), 0)
	if len(data) == 0 {
		t.Errorf("Valid KISS frames decode failed")
	}
}

func TestDecodeFreeformKISS(t *testing.T) {
	input := "C0 00 DB DC 80 2D 09 02 66 81 05 00 00 58 00 59 C0" +
		"C0 00 DB DC 80 2F 09 02 66 8B 05 00 00 59 00 5B C0"
	data, _ := DecodeFreeform("strand1", []byte(input), 0)
	if len(data) == 0 {
		t.Errorf("Valid KISS frames decode failed")
	}
}

func TestDecodeFreeformKISSMultiline(t *testing.T) {
	input := `
C0 00 DB DC 80 2D 09 02 66 81 05 00 00 58 00 59 C0
C0 00 DB DC 80 2E 09 02 66 86 05 00 03 0A 02 04 C0
C0 00 DB DC 80 2F 09 02 66 8B 05 00 00 59 00 5B C0
C0 00 DB DC 80 30 09 02 66 90 05 00 00 57 00 5A C0
C0 00 DB DC 80 31 09 02 66 95 05 00 00 59 00 57 C0
C0 00 DB DC 80 32 09 02 66 9A 05 00 00 58 00 59 C0
`
	data, _ := DecodeFreeform("strand1", []byte(input), 0)
	if len(data) == 0 {
		t.Errorf("Valid KISS frames decode failed")
	}
}

func TestDecodeFreeformKISSWithCrap(t *testing.T) {
	input := `
1[14 Bytes unknown Protocol]   1 > C0 10 DB DC 80 5A 06 01 E0 00 00 21 04 C0 À.ÛÜ€Z..à..!.À
[9 Bytes unknown Protocol]   1 > C0 10 01 E1 00 00 00 2E C0 À..á....À
[9 Bytes unknown Protocol]
   1 > C0 10 01 E3 00 00 00 0F C0 
À..ã....À

[9 Bytes unknown Protocol]
   1 > C0 10 01 E2 00 00 1F 7D C0 
À..â...}À[9 Bytes unknown Protocol]
   1 > C0 10 01 E4 00 00 00 00 C0 
À..ä....À
`
	data, _ := DecodeFreeform("strand1", []byte(input), 0)
	if len(data) == 0 {
		t.Errorf("Valid KISS frames (with crap) decode failed")
	}
}

func TestDecodeFreeformInvalid(t *testing.T) {
	data, _ := DecodeFreeform("fitsat1", ([]byte)("invalid stuff"), 0)
	if len(data) != 0 {
		t.Errorf("Data should be invalid.")
	}
}

