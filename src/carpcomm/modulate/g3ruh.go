// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

//import "fmt"

func G3RUHScramble(bits []bool) {
	// See http://www.amsat.org/amsat/articles/g3ruh/109/fig03.gif
	shift_reg := 0xffffff  // We use 17 bits.
	for i, b := range bits {
		var in int
		if b {
			in = 1
		}
		new_bit := ((shift_reg >> 16) ^ (shift_reg >> 11) ^ in) & 1
		shift_reg = (shift_reg << 1) | new_bit
		bits[i] = new_bit > 0
	}
}

func G3RUHDescramble(bits []bool) {
	// See http://www.amsat.org/amsat/articles/g3ruh/109/fig03.gif
	shift_reg := 0xffffff  // We use 17 bits.
	for i, b := range bits {
		var in int
		if b {
			in = 1
		}
		new_bit := ((shift_reg >> 16) ^ (shift_reg >> 11) ^ in) & 1
		shift_reg = (shift_reg << 1) | in
		bits[i] = new_bit > 0
	}
}

func modulateNRZ(bits []bool, sample_rate float64, baud float64) []float64 {
	samples_per_bit := sample_rate / baud
	samples := make([]float64, 0)
	num_samples := 0.0
	for _, b := range bits {
		num_samples += samples_per_bit
		for ; num_samples > 0.5; num_samples -= 1.0 {
			var s float64
			if b {
				s = 1.0
			} else {
				s = -1.0
			}
			samples = append(samples, s)
		}
	}
	return samples
}

// G3RUH 9600-baud modulator.
// This is based on the information here:
// http://www.amsat.org/amsat/articles/g3ruh/109.html
// Note that it outputs a digital signal which is meant to be fed into an
// NBFM modulator.
func ModulateG3RUH(bits []bool, sample_rate float64) []float64 {
	const N = 100
	bits = append(append(make([]bool, N), bits...), make([]bool, N)...)
	NRZIEncode(bits)
	G3RUHScramble(bits)
	return modulateNRZ(bits, sample_rate, 9600)
}

func ModulateG3RUHComplex(bits []bool, sample_rate, carrier float64) []complex64 {
	const N = 100
	bits = append(append(make([]bool, N), bits...), make([]bool, N)...)
	NRZIEncode(bits)
	G3RUHScramble(bits)
	markHz := carrier + 2400
	spaceHz := carrier - 2400
	return modulateFSKComplex(bits, sample_rate, markHz, spaceHz, 9600)
}


func DemodulateG3RUHIQ(samples []complex64, sample_rate, carrier float64) (
	packets [][]byte) {
	const baud = 9600.0
	markHz := carrier + 2400
	spaceHz := carrier - 2400
	
	samples_per_bit := sample_rate / baud
	clen := int(samples_per_bit)

	// For clock recovery / bit synchronization, we simply try many
	// different bit streams. At least one of them will be synchronized.
 	bits := make([][]bool, clen)

	excess := 0.0
	for i := clen; i < len(samples); {
		for j := i; j < i + clen; j++ {
			bit := correlateBitComplex(
				samples[j-clen:j], sample_rate, markHz, spaceHz)
			bits[j-i] = append(bits[j-i], bit)

			if (bit) {
				//fmt.Printf("1")
			} else {
				//fmt.Printf("0")
			}
		}
		//fmt.Printf("\n")
		skip := samples_per_bit + excess
		i += int(skip)
		excess = skip - float64(int(skip))
	}

	// We may detect the same packet multiple times.
	// Store them in a map to remove duplicates.
	packet_set := make(map[string][]byte)

	for i := 0; i < len(bits); i++ {
		G3RUHDescramble(bits[i])
		NRZIDecode(bits[i])
		for _, p := range DecodeHDLC(bits[i]) {
			packet_set[string(p)] = p
		}
	}

	for _, v := range packet_set {
		packets = append(packets, v)
	}
	return packets

}