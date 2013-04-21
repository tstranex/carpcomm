// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

// 0 set(['System time'])
// 1 set(['System time'])
// 2 set(['System time'])
// 3 set(['System time'])
// 18 set(['RAAC'])
// 19 set(['RAAC'])
// 20 set(['Avg. Current'])
// 21 set(['Avg. Current'])
// 22 set(['Temp'])
// 23 set(['Temp'])
// 24 set(['Volts'])
// 25 set(['Volts'])
// 26 set(['Inst. Current'])
// 27 set(['Inst. Current'])
// 28 set(['RAAC'])
// 29 set(['RAAC'])
// 30 set(['Avg. Current'])
// 31 set(['Avg. Current'])
// 32 set(['Temp'])
// 33 set(['Temp'])
// 34 set(['Volts'])
// 35 set(['Volts'])
// 36 set(['Inst. Current'])
// 37 set(['Inst. Current'])
// 38 set(['SA1 ( +X )'])
// 39 set(['SA2 ( -X )'])
// 40 set(['SA3 ( -Y )'])
// 41 set(['SA4 ( +Y )'])
// 42 set(['SA5 ( -Z )'])
// 43 set(['SA6 ( +Z )'])
// 44 set(['3.3V Reg. Cur'])
// 45 set(['5V Reg Cur'])
// 46 set(['3.3V Reg Volts'])
// 47 set(['3.3V Reg Volts'])
// 48 set(['5V Reg Volts'])
// 49 set(['5V Reg Volts'])
// 50 set(['3.3V Reg Temp'])
// 51 set(['5V Reg Temp'])
// 52 set(['Controller Temp'])
// 53 set(['Amplifier Temp'])
// 54 set(['RSSI [dB]'])
// 55 set(['File System Pointer'])
// 56 set(['File System Pointer'])
// 57 set(['Failed Packet Count'])
// 58 set(['Failed Packet Count'])
// 59 set(['Comm RX Count'])
// 60 set(['Comm RX Count'])
// 61 set(['Comm RX Count'])
// 62 set(['Comm RX Count'])
// 63 set(['Subsystem Status'])
// 64 set(['Label5'])

import "errors"
import "carpcomm/pb"

const hrbeCallsign = "K7MSU-1"

/*
func loadUint32BigEndian(f []byte) uint32 {
	return (uint32(f[0]) << 24) + (uint32(f[1]) << 16) +
		(uint32(f[2]) << 8) + uint32(f[3])
}
*/

func DecodeFrame_hrbe(frame []byte, timestamp int64) (
	data []pb.TelemetryDatum, err error) {

	// First do some quick validation.
	if len(frame) < 81 {
		return nil, errors.New("Frame is too short")
	}
	if string(frame[4:11]) != hrbeCallsign {
		return nil, errors.New("Frame has wrong callsign")
	}

	// System time.
	// uint32, big endian
	/*
	system_time := loadUint32BigEndian(frame[17:21])
	data = append(data, NewDoubleDatum(
		"d:hrbe:elapsed_s", timestamp, float64(system_time)))
	 */

	// Comm rx count.
	// uint32, big endian
		/*
	comm_rx_count := loadUint32BigEndian(frame[76:80])
	data = append(data, NewInt64Datum(
		"d:hrbe:rx_count", timestamp, int64(comm_rx_count)))
		 */

	// Battery current.
	InstCurrent := func(f []byte) float64 {
		c := (int16(f[0]) << 8) | int16(f[1])
		return 0.0390625 * float64(c) * 1e-3
	}
	data = append(data, NewDoubleDatum(
		"d:hrbe:bat1_c", timestamp, InstCurrent(frame[44:46])))
	data = append(data, NewDoubleDatum(
		"d:hrbe:bat2_c", timestamp, InstCurrent(frame[54:56])))

	// Regulator temperatures
	/*
	const reg_scale = 1.20698736
	data = append(data, NewDoubleDatum(
		"d:hrbe:3p3v_reg_t", timestamp, reg_scale*float64(frame[68])))
	data = append(data, NewDoubleDatum(
		"d:hrbe:5v_reg_t", timestamp, reg_scale*float64(frame[69])))
	 */

	return data, nil
}