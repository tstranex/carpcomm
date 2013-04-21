// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"
import "encoding/hex"

const hrbeTestFrame = "337205024b374d53552d312009200800841dc63d3c0020fecd290c31309090bb000ff08f11d7da5cfd005b80fe3e0453f62504805ba0fc0709060607070352210267026ed3d45055d5a000000000000000bf0258353535353535000000000000000000000000000000000000000000000001000100000000000000000000000000000000000000000000000000000001000000010000"

// System time;RAAC;Avg. Current;Temp;Volts;Inst. Current;RAAC;Avg. Current;Temp;Volts;Inst. Current;SA1 ( +X );SA2 ( -X );SA3 ( -Y );SA4 ( +Y );SA5 ( -Z );SA6 ( +Z );3.3V Reg. Cur;5V Reg Cur;3.3V Reg Volts;5V Reg Volts;3.3V Reg Temp;5V Reg Temp;Controller Temp;Amplifier Temp;RSSI [dB];File System Pointer;Failed Packet Count;Comm RX Count;Subsystem Status;Label5;
// 499531068;36625;-112.8125;125;3.57216;-17.578125;1107;-197.109375;4.5;3.57704;-39.7265625;5.68704044117647;3.79136029411765;3.79136029411765;4.42325367647059;4.42325367647059;1.89568014705882;77.7228860294118;31.2787224264706;3.3098291015625;5.01123046875;-18.4756671348315;-17.2686797752809;-15.3375;0.775781250000023;-84.413642578125;40960;0;0;191;600;

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