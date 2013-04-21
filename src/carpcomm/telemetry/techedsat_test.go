// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"

func Test_techedsat_DecodeADC(t *testing.T) {
	d, err := techedsat_DecodeADC("test", 0, "A64", -1.0, 9.0, 0.1, -2.0)
	if err != nil {
		t.Error(err)
	}
	ExpectDoubleDatum(t, d, "test", -1.450427)
}


func TestDecodeFrame_techedsat_1(t *testing.T) {
	frame := "0000000000000000ncasst.org0000057fcadcadcadcad85d85d85d85dbf4bf6bf5bf585d85d85d85dcc3cc3cc3cc332834a33933482d82d82d82d03000000170000012d2b"
	data, err := DecodeFrame_techedsat(([]byte)(frame), 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 39 {
		t.Errorf("Wrong number of datums: %d, expected 39", len(data))
	}

	ExpectDoubleDatum(t, data[0], "d:techedsat:elapsed_s", 1407.0)

	// ADC values
	ExpectDoubleDatum(t, data[1], "d:techedsat:5v_min_v", 4.95836)
	ExpectDoubleDatum(t, data[2], "d:techedsat:5v_max_v", 4.95836)
	ExpectDoubleDatum(t, data[3], "d:techedsat:5v_avg_v", 4.95836)
	ExpectDoubleDatum(t, data[4], "d:techedsat:5v_now_v", 4.95836)
	ExpectDoubleDatum(t, data[5], "d:techedsat:3v3_min_v", 3.271448)
	ExpectDoubleDatum(t, data[6], "d:techedsat:3v3_max_v", 3.271448)
	ExpectDoubleDatum(t, data[7], "d:techedsat:3v3_avg_v", 3.271448)
	ExpectDoubleDatum(t, data[8], "d:techedsat:3v3_now_v", 3.271448)
	ExpectDoubleDatum(t, data[9], "d:techedsat:min_t", 23.4285714285714 + 273.15)
	ExpectDoubleDatum(t, data[10], "d:techedsat:max_t", 23.7142857142857 + 273.15)
	ExpectDoubleDatum(t, data[11], "d:techedsat:avg_t", 23.5714285714286 + 273.15)
	ExpectDoubleDatum(t, data[12], "d:techedsat:now_t", 23.5714285714286 + 273.15)
	ExpectDoubleDatum(t, data[13], "d:techedsat:min_c", 0.15370239)
	ExpectDoubleDatum(t, data[14], "d:techedsat:max_c", 0.15370239)
	ExpectDoubleDatum(t, data[15], "d:techedsat:avg_c", 0.15370239)
	ExpectDoubleDatum(t, data[16], "d:techedsat:now_c", 0.15370239)
	ExpectDoubleDatum(t, data[17], "d:techedsat:bat_min_v", 7.81968060661765)
	ExpectDoubleDatum(t, data[18], "d:techedsat:bat_max_v", 7.81968060661765)
	ExpectDoubleDatum(t, data[19], "d:techedsat:bat_avg_v", 7.81968060661765)
	ExpectDoubleDatum(t, data[20], "d:techedsat:bat_now_v", 7.81968060661765)
	ExpectDoubleDatum(t, data[21], "d:techedsat:bat_min_c", 0.654280679933665)
	ExpectDoubleDatum(t, data[22], "d:techedsat:bat_max_c", 0.681812292703151)
	ExpectDoubleDatum(t, data[23], "d:techedsat:bat_avg_c", 0.668046486318408)
	ExpectDoubleDatum(t, data[24], "d:techedsat:bat_now_c", 0.66399771973466)
	ExpectDoubleDatum(t, data[25], "d:techedsat:charge_min_c", 0.076559777874564)
	ExpectDoubleDatum(t, data[26], "d:techedsat:charge_max_c", 0.076559777874564)
	ExpectDoubleDatum(t, data[27], "d:techedsat:charge_avg_c", 0.076559777874564)
	ExpectDoubleDatum(t, data[28], "d:techedsat:charge_now_c", 0.076559777874564)

	ExpectDoubleDatum(t, data[29], "d:techedsat:nominal_mode_s", 0.0)
	ExpectDoubleDatum(t, data[30], "d:techedsat:safe_mode_s", 1380.0)

	ExpectInt64Datum(t, data[31], "d:techedsat:single_errors", 0)
	ExpectInt64Datum(t, data[32], "d:techedsat:all_errors", 0)
	ExpectInt64Datum(t, data[33], "d:techedsat:all_errors_2", 1)

	ExpectBoolDatum(t, data[34], "d:techedsat:cell5_active", true)
	ExpectBoolDatum(t, data[35], "d:techedsat:cell4_active", true)
	ExpectBoolDatum(t, data[36], "d:techedsat:cell3_active", false)
	ExpectBoolDatum(t, data[37], "d:techedsat:cell2_active", false)
	ExpectBoolDatum(t, data[38], "d:techedsat:cell1_active", false)
}

func TestDecodeFrame_techedsat_2(t *testing.T) {
	frame := "0000000000000000ncasst.org00000c8dcadcadcadcad85d85d85d85dc02c03c02c0285d85d85d85dd31d31d31d3114b14c14b14c82e83082f8300300350000000001d25f"
	data, err := DecodeFrame_techedsat(([]byte)(frame), 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 39 {
		t.Errorf("Wrong number of datums: %d, expected 39", len(data))
	}

	ExpectDoubleDatum(t, data[0], "d:techedsat:elapsed_s", 3213.0)

	// ADC values
	ExpectDoubleDatum(t, data[1], "d:techedsat:5v_min_v", 4.95836)
	ExpectDoubleDatum(t, data[2], "d:techedsat:5v_max_v", 4.95836)
	ExpectDoubleDatum(t, data[3], "d:techedsat:5v_avg_v", 4.95836)
	ExpectDoubleDatum(t, data[4], "d:techedsat:5v_now_v", 4.95836)
	ExpectDoubleDatum(t, data[5], "d:techedsat:3v3_min_v", 3.271448)
	ExpectDoubleDatum(t, data[6], "d:techedsat:3v3_max_v", 3.271448)
	ExpectDoubleDatum(t, data[7], "d:techedsat:3v3_avg_v", 3.271448)
	ExpectDoubleDatum(t, data[8], "d:techedsat:3v3_now_v", 3.271448)
	ExpectDoubleDatum(t, data[9], "d:techedsat:min_t", 25.4285714285714 + 273.15)
	ExpectDoubleDatum(t, data[10], "d:techedsat:max_t", 25.5714285714286 + 273.15)
	ExpectDoubleDatum(t, data[11], "d:techedsat:avg_t", 25.4285714285714 + 273.15)
	ExpectDoubleDatum(t, data[12], "d:techedsat:now_t", 25.4285714285714 + 273.15)
	ExpectDoubleDatum(t, data[13], "d:techedsat:min_c", 0.15370239)
	ExpectDoubleDatum(t, data[14], "d:techedsat:max_c", 0.15370239)
	ExpectDoubleDatum(t, data[15], "d:techedsat:avg_c", 0.15370239)
	ExpectDoubleDatum(t, data[16], "d:techedsat:now_c", 0.15370239)
	ExpectDoubleDatum(t, data[17], "d:techedsat:bat_min_v", 8.08296951593137)
	ExpectDoubleDatum(t, data[18], "d:techedsat:bat_max_v", 8.08296951593137)
	ExpectDoubleDatum(t, data[19], "d:techedsat:bat_avg_v", 8.08296951593137)
	ExpectDoubleDatum(t, data[20], "d:techedsat:bat_now_v", 8.08296951593137)
	ExpectDoubleDatum(t, data[21], "d:techedsat:bat_min_c", 0.268028347844113)
	ExpectDoubleDatum(t, data[22], "d:techedsat:bat_max_c", 0.268838101160862)
	ExpectDoubleDatum(t, data[23], "d:techedsat:bat_avg_c", 0.268028347844113)
	ExpectDoubleDatum(t, data[24], "d:techedsat:bat_now_c", 0.268838101160862)
	ExpectDoubleDatum(t, data[25], "d:techedsat:charge_min_c", 0.078261106271777)
	ExpectDoubleDatum(t, data[26], "d:techedsat:charge_max_c", 0.081663763066202)
	ExpectDoubleDatum(t, data[27], "d:techedsat:charge_avg_c", 0.07996243466899)
	ExpectDoubleDatum(t, data[28], "d:techedsat:charge_now_c", 0.081663763066202)

	ExpectDoubleDatum(t, data[29], "d:techedsat:nominal_mode_s", 3180.0)
	ExpectDoubleDatum(t, data[30], "d:techedsat:safe_mode_s", 0.0)

	ExpectInt64Datum(t, data[31], "d:techedsat:single_errors", 0)
	ExpectInt64Datum(t, data[32], "d:techedsat:all_errors", 0)
	ExpectInt64Datum(t, data[33], "d:techedsat:all_errors_2", 1)

	ExpectBoolDatum(t, data[34], "d:techedsat:cell5_active", true)
	ExpectBoolDatum(t, data[35], "d:techedsat:cell4_active", true)
	ExpectBoolDatum(t, data[36], "d:techedsat:cell3_active", false)
	ExpectBoolDatum(t, data[37], "d:techedsat:cell2_active", false)
	ExpectBoolDatum(t, data[38], "d:techedsat:cell1_active", false)
}

func TestDecodeFrame_techedsat_3(t *testing.T) {
	frame := "0000000000000000ncasst.org0004fccccb1cb2cb1cb286b86c86b86bb98b99b98b9886a86b86a86bda0da1da0da07967b67a57b684284384284303150c0000000001af9a"
	data, err := DecodeFrame_techedsat(([]byte)(frame), 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 39 {
		t.Errorf("Wrong number of datums: %d, expected 39", len(data))
	}

	ExpectDoubleDatum(t, data[0], "d:techedsat:elapsed_s", 326860.0)

	// ADC values
	ExpectDoubleDatum(t, data[1], "d:techedsat:5v_min_v", 4.964472)
	ExpectDoubleDatum(t, data[2], "d:techedsat:5v_max_v", 4.966)
	ExpectDoubleDatum(t, data[3], "d:techedsat:5v_avg_v", 4.964472)
	ExpectDoubleDatum(t, data[4], "d:techedsat:5v_now_v", 4.966)
	ExpectDoubleDatum(t, data[5], "d:techedsat:3v3_min_v", 3.29284)
	ExpectDoubleDatum(t, data[6], "d:techedsat:3v3_max_v", 3.294368)
	ExpectDoubleDatum(t, data[7], "d:techedsat:3v3_avg_v", 3.29284)
	ExpectDoubleDatum(t, data[8], "d:techedsat:3v3_now_v", 3.29284)
	ExpectDoubleDatum(t, data[9], "d:techedsat:min_t", 10.2857142857143 + 273.15)
	ExpectDoubleDatum(t, data[10], "d:techedsat:max_t", 10.4285714285714 + 273.15)
	ExpectDoubleDatum(t, data[11], "d:techedsat:avg_t", 10.2857142857143 + 273.15)
	ExpectDoubleDatum(t, data[12], "d:techedsat:now_t", 10.2857142857143 + 273.15)
	ExpectDoubleDatum(t, data[13], "d:techedsat:min_c", 0.15463566)
	ExpectDoubleDatum(t, data[14], "d:techedsat:max_c", 0.15470745)
	ExpectDoubleDatum(t, data[15], "d:techedsat:avg_c", 0.15463566)
	ExpectDoubleDatum(t, data[16], "d:techedsat:now_c", 0.15470745)
	ExpectDoubleDatum(t, data[17], "d:techedsat:bat_min_v", 8.34865196078432)
	ExpectDoubleDatum(t, data[18], "d:techedsat:bat_max_v", 8.35104549632353)
	ExpectDoubleDatum(t, data[19], "d:techedsat:bat_avg_v", 8.34865196078432)
	ExpectDoubleDatum(t, data[20], "d:techedsat:bat_now_v", 8.34865196078432)
	ExpectDoubleDatum(t, data[21], "d:techedsat:bat_min_c", 1.5725409411277)
	ExpectDoubleDatum(t, data[22], "d:techedsat:bat_max_c", 1.59845304726368)
	ExpectDoubleDatum(t, data[23], "d:techedsat:bat_avg_c", 1.58468724087894)
	ExpectDoubleDatum(t, data[24], "d:techedsat:bat_now_c", 1.59845304726368)
	ExpectDoubleDatum(t, data[25], "d:techedsat:charge_min_c", 0.112287674216028)
	ExpectDoubleDatum(t, data[26], "d:techedsat:charge_max_c", 0.11398900261324)
	ExpectDoubleDatum(t, data[27], "d:techedsat:charge_avg_c", 0.112287674216028)
	ExpectDoubleDatum(t, data[28], "d:techedsat:charge_now_c", 0.11398900261324)

	ExpectDoubleDatum(t, data[29], "d:techedsat:nominal_mode_s", 5388*60.0)
	ExpectDoubleDatum(t, data[30], "d:techedsat:safe_mode_s", 0.0)

	ExpectInt64Datum(t, data[31], "d:techedsat:single_errors", 0)
	ExpectInt64Datum(t, data[32], "d:techedsat:all_errors", 0)
	ExpectInt64Datum(t, data[33], "d:techedsat:all_errors_2", 1)

	ExpectBoolDatum(t, data[34], "d:techedsat:cell5_active", true)
	ExpectBoolDatum(t, data[35], "d:techedsat:cell4_active", true)
	ExpectBoolDatum(t, data[36], "d:techedsat:cell3_active", false)
	ExpectBoolDatum(t, data[37], "d:techedsat:cell2_active", false)
	ExpectBoolDatum(t, data[38], "d:techedsat:cell1_active", false)
}
