// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"
import "encoding/hex"
import "strings"
import "carpcomm/util"
import "time"

func strand1ExpectDouble(t *testing.T, hexframe string,
	key string, v float64) {
	frame, _ := hex.DecodeString(hexframe)
	data, err := DecodeFrame_strand1(frame, 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Errorf("Wrong number of datums: %d, expected 1", len(data))
		return
	}
	ExpectDoubleDatum(t, data[0], key, v)
}

func strand1ExpectInt64(t *testing.T, hexframe string,
	key string, v int64) {
	frame, _ := hex.DecodeString(hexframe)
	data, err := DecodeFrame_strand1(frame, 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Errorf("Wrong number of datums: %d, expected 1", len(data))
		return
	}
	ExpectInt64Datum(t, data[0], key, v)
}

func TestDecodeFrame_strand1(t *testing.T) {
	strand1ExpectDouble(t, "C080A406022C0302AC00",
		"d:strand1:bat0_v", 8.1234)
	strand1ExpectDouble(t, "C080A506022C0802AD00",
		"d:strand1:bat1_v", 8.11602)
	strand1ExpectDouble(t, "C080A606022D0D024101",
		"d:strand1:cell_x2_c", 0.346304)
	strand1ExpectDouble(t, "C080A706022D0702FF03",
		"d:strand1:cell_x1_c", -0.0268354)
}

func TestDecodeFrame_strand1_switch(t *testing.T) {
	frame, _ := hex.DecodeString("C080B40902669F05000059005A")
	data, err := DecodeFrame_strand1(frame, 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 3 {
		t.Errorf("Wrong number of datums: %d, expected 3", len(data))
		return
	}

	ExpectBoolDatum(t, data[0], "d:strand1:switch6_on", false)
	ExpectDoubleDatum(t, data[1], "d:strand1:switch6_c", 0.022271044)
	ExpectDoubleDatum(t, data[2], "d:strand1:switch6_v", -0.790967349)
}

func TestDecodeFrame_strand1_ModemBeacon(t *testing.T) {
	strand1ExpectDouble(t, "C080300601E000002104",
		"d:strand1:time_since_last_obc_packet", 8452.0)
	strand1ExpectInt64(t, "01E10000002E",
		"d:strand1:packet_up_count", 0x2e)
	strand1ExpectInt64(t, "01E200001EDD",
		"d:strand1:packet_down_count", 0x1edd)
	strand1ExpectInt64(t, "01E30000000F",
		"d:strand1:packet_up_dropped_count", 0xf)
	strand1ExpectInt64(t, "01E400000000",
		"d:strand1:packet_down_dropped_count", 0)
}

func TestDecodeFrame_strand1_OBC(t *testing.T) {
	frame, _ := hex.DecodeString("C080770C02800C0887110100D0E90400")
	data, err := DecodeFrame_strand1(frame, 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Errorf("Wrong number of datums: %d, expected 1", len(data))
		return
	}
	ExpectTimestampDatum(t,
		data[0], "d:strand1:obc_clock", time.Unix(70023, 0))
}

func TestDecodeFrame_strand1_Magnetometer(t *testing.T) {
	frame, _ := hex.DecodeString("C080AC0C02890308ECCCFFFF20DDFDFF")
	data, err := DecodeFrame_strand1(frame, 123)
	if err != nil {
		t.Error(err)
	}
	if len(data) != 2 {
		t.Errorf("Wrong number of datums: %d, expected 2", len(data))
		return
	}
	ExpectInt64Datum(t, data[0], "d:strand1:magnetometer_1x", -13076)
	ExpectInt64Datum(t, data[1], "d:strand1:magnetometer_1y", -140000)

	strand1ExpectInt64(t, "C080AD0C02890508533D030011770100",
		"d:strand1:magnetometer_2z", 212307)
}

var strand1Packets string = `
C0 80 A4 06 02 2C 03 02 AC 00 
C0 80 CE 06 02 2C 03 02 AC 00 
C0 80 CF 06 02 2C 08 02 AC 00 
C0 80 A5 06 02 2C 08 02 AD 00 
C0 80 D2 06 02 2D 01 02 04 01 
C0 80 A9 06 02 2D 04 02 CC 03 
C0 80 B4 09 02 66 9F 05 00 00 59 00 5A 
C0 80 A8 06 02 2D 01 02 FF 03 
C0 80 D4 06 02 2D 0A 02 EC 02 
C0 80 DB 09 02 66 90 05 00 00 59 00 59 
C0 80 DD 09 02 66 9A 05 00 00 57 00 58 
C0 80 9E 09 02 66 9A 05 00 00 59 00 5A 
C0 80 D3 06 02 2D 04 02 CF 03 
C0 80 A6 06 02 2D 0D 02 41 01 
C0 80 AB 06 02 2D 1F 02 5E 02 
C0 80 DC 09 02 66 95 05 00 00 58 00 58 
C0 80 DE 09 02 66 9F 05 00 00 58 00 58 
C0 80 9B 09 02 66 8B 05 00 00 58 00 5A 
C0 80 B2 09 02 66 95 05 00 00 58 00 57 
C0 80 B1 09 02 66 90 05 00 00 5A 00 59 
C0 80 A7 06 02 2D 07 02 FF 03 
C0 80 AF 09 02 66 86 05 00 03 08 02 05 
C0 80 D1 06 02 2D 07 02 FF 03 
C0 80 D5 06 02 2D 1F 02 FF 03 
C0 80 AA 06 02 2D 0A 02 FF 03 
C0 80 D9 09 02 66 86 05 00 03 07 02 03 
C0 80 9F 09 02 66 9F 05 00 00 5A 00 57 
C0 80 95 06 02 2D 0A 02 4B 01 
C0 80 9D 09 02 66 95 05 00 00 59 00 58 
C0 80 AE 09 02 66 81 05 00 00 58 00 5C 
C0 80 D0 06 02 2D 0D 02 CB 03 
C0 80 D8 09 02 66 81 05 00 00 58 00 59 
C0 80 96 06 02 2D 1F 02 FF 03 
C0 80 9A 09 02 66 86 05 00 03 07 02 03 
C0 80 99 09 02 66 81 05 00 00 57 00 5A 
C0 80 DA 09 02 66 8B 05 00 00 58 00 59 
C0 80 A0 09 02 66 A4 05 00 00 56 00 5B
C0 80 CD 0C 02 80 0C 08 22 77 01 00 78 3F 07 00
C0 80 A3 0C 02 80 0C 08 0C 77 01 00 B0 36 00 00
C0 80 AD 0C 02 89 05 08 53 3D 03 00 11 77 01 00
C0 80 AC 0C 02 89 03 08 EC CC FF FF 20 DD FD FF
C0 80 D6 0C 02 89 03 08 7E A6 03 00 4B B7 00 00
C0 80 97 0C 02 89 03 08 A7 24 FF FF 72 A3 FE FF
C0 80 98 0C 02 89 05 08 A2 FD 03 00 06 77 01 00
C0 80 D7 0C 02 89 05 08 70 11 01 00 27 77 01 00
`

var strand1ErrorPackets = `
C0 80 A1 09 02 66 A9 05 00 00 00 00 00 
C0 80 A2 09 02 66 AC 05 00 00 00 00 00 
C0 80 CC 09 02 66 AC 05 00 00 00 00 00 
`

func TestDecodeFrame_strand1_all(t *testing.T) {
	hexpackets := strings.Split(strand1Packets, "\n")
	for _, p := range hexpackets {
		stripped := util.StripWhitespace(p)
		if stripped == "" {
			continue
		}
		frame, _ := hex.DecodeString(stripped)
		data, err := DecodeFrame_strand1(frame, 123)
		if err != nil {
			t.Errorf("Error for %s: %s", stripped, err)
		}
		if len(data) == 0 {
			t.Errorf("No data decoded for %s", stripped)
		}
	}
}