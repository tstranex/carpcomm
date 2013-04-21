// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "testing"
import "encoding/hex"

// From http://www.dk3wn.info/p/?p=28430
const cssweFrame1 = "86A6A6AE8A6E6086A240404040E103F0420061F840C4FD61010000000000F4FBDA61D80E00000D00E10000CF250008FF03A90700CC042200FF022400000C160F0E000000000000CE00FEFEFEFECFCFCFCF000C006804F2AFFD6100CB220000FF03A9070000000100000000000004020203000000000000CE00FEFEFEFECFCFCFCF000C00680400CF280E0CFF04A90700D104C90DFF04ED11000E170810000000000000CE00FEFEFEFECFCFCFCF000D0071070ACD250206FF03A9070A4101760184016C030A080904090A0000000000CE0AFEFEFEFECFCFCFCF0A0C006D05001000000000000000000006000100000000000000001C2C0A000000600000009A0BAA000000748703000AC02200"

func TestDecodeCSSWE(t *testing.T) {
	frame, _ := hex.DecodeString(cssweFrame1)

	data, err := DecodeFrame_csswe(frame, 123)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 29 {
		t.Errorf("Wrong number of datums: %d, expected 29", len(data))
		return
	}

	ExpectDoubleDatum(t, data[0], "d:csswe:batt_v", 8.117646) //8.1569)
	ExpectDoubleDatum(t, data[1], "d:csswe:3p3v_bus_v", 3.313725) //3.3137)
	ExpectDoubleDatum(t, data[2], "d:csswe:5v_bus_v", 5.0)
	ExpectDoubleDatum(t, data[3], "d:csswe:batt_charge_c", 0.062745) //0.0862744)
	ExpectDoubleDatum(t, data[4], "d:csswe:batt_discharge_c", 0.0)
	ExpectDoubleDatum(t, data[5], "d:csswe:3p3v_bus_c", 0.0549019)
	ExpectDoubleDatum(t, data[6], "d:csswe:5v_bus_c", 0.0235294)

	ExpectDoubleDatum(t, data[7], "d:csswe:pvx1_v", 17.315999) //16.5521)
	ExpectDoubleDatum(t, data[8], "d:csswe:pvx2_v", 2.886000) //14.0056)
	ExpectDoubleDatum(t, data[9], "d:csswe:pvy1_v", 21.644999) //15.6184)
	ExpectDoubleDatum(t, data[10], "d:csswe:pvy2_v", 3.055765) //17.5706)
	ExpectDoubleDatum(t, data[11], "d:csswe:pvx1_c", 0.031373) //0.0156863)
	ExpectDoubleDatum(t, data[12], "d:csswe:pvx2_c", 0.000000) //0.0313725)
	ExpectDoubleDatum(t, data[13], "d:csswe:pvy1_c", 0.015686) //0.0)
	ExpectDoubleDatum(t, data[14], "d:csswe:pvy2_c", 0.000000) //0.0156863)

	ExpectDoubleDatum(t, data[15], "d:csswe:batt_t", 286.150000) //11.0 + 273.15)
	ExpectDoubleDatum(t, data[16], "d:csswe:cdh_t", 277.150000) //3.0 + 273.15)
	ExpectDoubleDatum(t, data[17], "d:csswe:radio_t", 12.0 + 273.15)
	ExpectDoubleDatum(t, data[18], "d:csswe:pvx1_t", 262.150000) //5.0 + 273.15)
	ExpectDoubleDatum(t, data[19], "d:csswe:pvx2_t", 274.150000) //14.0 + 273.15)
	ExpectDoubleDatum(t, data[20], "d:csswe:pvy1_t", 266.150000) //-11.0 + 273.15)
	ExpectDoubleDatum(t, data[21], "d:csswe:pvy2_t", 265.150000) //33.0 + 273.15)

	ExpectBoolDatum(t, data[22], "d:csswe:satellite_mode", false)
	ExpectInt64Datum(t, data[23], "d:csswe:cmds_since_boot", 0)
	ExpectDoubleDatum(t, data[24], "d:csswe:time_since_boot_s", 22799.998541) //34983.36)
	ExpectBoolDatum(t, data[25], "d:csswe:gpio_reptile_3p3v", true)
	ExpectBoolDatum(t, data[26], "d:csswe:gpio_reptile_5p0v", false)
	ExpectBoolDatum(t, data[27], "d:csswe:gpio_battery_heater", false)
	ExpectBoolDatum(t, data[28], "d:csswe:gpio_adm_resistor", false)
}
