// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"
import "encoding/hex"

const aeneasFrame1 = "4341455255533B00020000060A0A05003A1E290000000000401701030000000D0008A6F1549026021CC5C6D02AF9F80104100A040C1630F8BFFDBF01802C8038803E80DD49"
// Decoded version of aeneasFrame1 according to beacondecoder.jar.
/*
Beacon TLM Packet found:
---------TELEMETRY-BEGIN----------
Total length: 59
Type:2
Unused: 0x00 0x00 
Time: 6/10/10-5-0:58:30
Reboots: 41
Last Reboot Cause: 0
Flash Status: 0
Bit 0 = 0 : Device is ready.
Bit 1 = 0 : Device is NOT write-enabled.
Bit3&4= 1 : All sectors software-protected.
Bit 4 = 0 : WP is asserted.
Bit 5 = 0 : Last op successful.
Bit 7 = 0 : Sector Protection Registers unlocked.
Radio Status: 64
No detailed Radio info yet
TLM Pointer: 196887
PLY Pointer: 851968
Last Sent Command: 0x08 0xa6 0xf1 0x54 0x90 0x26 0x02 0x1c 0xc5 0xc6 0xd0 0x2a 0xf9 0xf8 0x01 
Time: 4/16/10-4-12:22:48
Gyro Rates Data Bytes: 0xf8 0xbf 0xfd 0xbf 0x01 0x80 
Gyro Rates (deg/s): A: -0.58608 B: -0.21978 C: 0.07326
Gyro Rates New Data Flags: A: true B: true C: true
Gyro Rates Error Flags: A: false B: false C: false
Gyro Temp Data Bytes: 0x2c 0x80 0x38 0x80 0x3e 0x80 
Gyro Temps (deg C): A: 31.3932 B: 33.1368 C: 34.0086
Gyro Temps New Data Flags: A: true B: true C: true
Gyro Temps Error Flags: A: false B: false C: false
Total count: 59 (Should be 83 or 59)
*/

func TestDecodeAeneas1(t *testing.T) {
	frame, _ := hex.DecodeString(aeneasFrame1)

	data, err := DecodeAeneas(frame, 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 7 {
		t.Errorf("Wrong number of datums: %d, expected 7", len(data))
		return
	}

	ExpectInt64Datum(t, data[0], "d:aeneas:num_reboots", 41)

	// These values differ slightly from beacondecoder.jar (ref).
	// The exact calibration formulae are unknown.

	// ref: -0.010229025680088367
	ExpectDoubleDatum(t, data[1], "d:aeneas:gyro1_r", -0.010261)
	// ref: -0.0038358846300331375
	ExpectDoubleDatum(t, data[2], "d:aeneas:gyro2_r", -0.003758)
	// ref: 0.0012786282100110458
	ExpectDoubleDatum(t, data[3], "d:aeneas:gyro3_r", 0.001124)

	// ref: 304.54319999999996
	ExpectDoubleDatum(t, data[4], "d:aeneas:gyro1_t", 304.548733)
	// ref: 306.28679999999997
	ExpectDoubleDatum(t, data[5], "d:aeneas:gyro2_t", 306.282567)
	// ref: 307.15859999999998
	ExpectDoubleDatum(t, data[6], "d:aeneas:gyro3_t", 307.149484)
}


const aeneasFrame2 = "4341455255533B00020000060A0A0500311E270000000000405D01030000000D0008A2F1549026021CC5C6D02AF9F80104100A040C1612F9BFFFBFFDBFF38FF58FF58F7117"
// Decoded version of aeneasFrame2 according to beacondecoder.jar.
/*
Beacon TLM Packet found:
---------TELEMETRY-BEGIN----------
Total length: 59
Type:2
Unused: 0x00 0x00 
Time: 6/10/10-5-0:49:30
Reboots: 39
Last Reboot Cause: 0
Flash Status: 0
Bit 0 = 0 : Device is ready.
Bit 1 = 0 : Device is NOT write-enabled.
Bit3&4= 1 : All sectors software-protected.
Bit 4 = 0 : WP is asserted.
Bit 5 = 0 : Last op successful.
Bit 7 = 0 : Sector Protection Registers unlocked.
Radio Status: 64
No detailed Radio info yet
TLM Pointer: 196957
PLY Pointer: 851968
Last Sent Command: 0x08 0xa2 0xf1 0x54 0x90 0x26 0x02 0x1c 0xc5 0xc6 0xd0 0x2a 0xf9 0xf8 0x01 
Time: 4/16/10-4-12:22:18
Gyro Rates Data Bytes: 0xf9 0xbf 0xff 0xbf 0xfd 0xbf 
Gyro Rates (deg/s): A: -0.51282 B: -0.07326 C: -0.21978
Gyro Rates New Data Flags: A: true B: true C: true
Gyro Rates Error Flags: A: false B: false C: false
Gyro Temp Data Bytes: 0xf3 0x8f 0xf5 0x8f 0xf5 0x8f 
Gyro Temps (deg C): A: 23.1111 B: 23.4017 C: 23.4017
Gyro Temps New Data Flags: A: true B: true C: true
Gyro Temps Error Flags: A: false B: false C: false
Total count: 59 (Should be 83 or 59)
---------TELEMETRY-END----------
*/
func TestDecodeAeneas2(t *testing.T) {
	frame, _ := hex.DecodeString(aeneasFrame2)

	data, err := DecodeAeneas(frame, 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 7 {
		t.Errorf("Wrong number of datums: %d, expected 7", len(data))
		return
	}

	ExpectInt64Datum(t, data[0], "d:aeneas:num_reboots", 39)

	// These values differ slightly from beacondecoder.jar (ref).
	// The exact calibration formulae are unknown.

	// ref: -0.0089503974700773214
	ExpectDoubleDatum(t, data[1], "d:aeneas:gyro1_r", -0.008960)
	// ref: -0.0012786282100110458
	ExpectDoubleDatum(t, data[2], "d:aeneas:gyro2_r", -0.001157)
	// ref: -0.0038358846300331375
	ExpectDoubleDatum(t, data[3], "d:aeneas:gyro3_r", -0.003758)

	// ref: 296.2611
	ExpectDoubleDatum(t, data[4], "d:aeneas:gyro1_t", 296.321489)
	// ref: 296.55169999999998
	ExpectDoubleDatum(t, data[5], "d:aeneas:gyro2_t", 296.610461)
	// ref: 296.55169999999998
	ExpectDoubleDatum(t, data[6], "d:aeneas:gyro3_t", 296.610461)
}


func TestAeneasWrongLength(t *testing.T) {
	good_frame, _ := hex.DecodeString(aeneasFrame1)
	bad_frame := good_frame[:40]
	_, err := DecodeAeneas(bad_frame, 123)
	if err == nil {
		t.Errorf("Frame should be invalid.")
	}
}