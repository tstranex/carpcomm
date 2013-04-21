// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "math"
import "math/cmplx"
import "os"
import "log"
import "carpcomm/util/binary"
import "fmt"
import "bufio"
import "encoding/hex"
import "strings"

const markHz = 1200
const spaceHz = 2200
const baud = 1200

// Non-Return-to-Zero Inverted (NRZI) encoding:
// 0 causes a state transition and 1 does not.
func NRZIEncode(bits []bool) {
	state := false
	for i, b := range bits {
		if !b {
			state = !state
		}
		bits[i] = state
	}
}

// Non-Return-to-Zero Inverted (NRZI) decoding.
func NRZIDecode(bits []bool) {
	for i := 0; i < len(bits)-1; i++ {
		bits[i] = bits[i] == bits[i+1]
	}
}

func modulateFSK(bits []bool, sample_rate float64,
	zero_hz, one_hz, baud float64) []float64 {
	const π = math.Pi
	Δt := 1.0 / sample_rate
	Δφ_one :=  2 * π * one_hz * Δt
	Δφ_zero :=  2 * π * zero_hz * Δt
	samples_per_bit := sample_rate / baud

	samples := make([]float64, 0)
	φ := 0.0
	num_samples := 0.0
	for _, b := range bits {
		num_samples += samples_per_bit
		for ; num_samples > 0.5; num_samples -= 1.0 {
			if b {
				φ += Δφ_one
			} else {
				φ += Δφ_zero
			}
			samples = append(samples, math.Sin(φ))
		}
	}
	return samples
}

func modulateFSKComplex(bits []bool, sample_rate float64,
	zero_hz, one_hz, baud float64) []complex64 {
	const π = math.Pi
	Δt := 1.0 / sample_rate
	Δφ_one :=  2 * π * one_hz * Δt
	Δφ_zero :=  2 * π * zero_hz * Δt
	samples_per_bit := sample_rate / baud

	samples := make([]complex64, 0)
	φ := 0.0
	num_samples := 0.0
	for _, b := range bits {
		num_samples += samples_per_bit
		for ; num_samples > 0.5; num_samples -= 1.0 {
			if b {
				φ += Δφ_one
			} else {
				φ += Δφ_zero
			}
			samples = append(samples,
				complex64(cmplx.Exp(complex(0, φ))))
		}
	}
	return samples
}

// Returns a 1200 baud Bell 202 FSK signal that encodes the bits.
// This modulation is often used for APRS and AX.25 ametuer packet radio.
func ModulateAFSK1200(bits []bool, sample_rate float64) []float64 {
	// Start and end with a burst to aid tuning.
	bits = append(append(make([]bool, 100), bits...), make([]bool, 100)...)
	NRZIEncode(bits)
	return modulateFSK(bits, sample_rate, markHz, spaceHz, baud)
}

func writeToFile(path string, samples []float64) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error opening file: %s", err.Error())
	}

	b := make([]byte, 2)
	for i := 0; i < 1; i++ {
		for _, s := range samples {
			i := int16(s * 16348)
			b[0] = byte(i & 0xff)
			b[1] = byte((i >> 8) & 0xff)
			_, err := f.Write(b)
			if err != nil {
				log.Fatalf("Error writing: %s", err.Error())
			}
		}
	}

	f.Close()
}

func writeToFileComplex(path string, samples []complex64) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error opening file: %s", err.Error())
	}
	w := bufio.NewWriter(f)

	for _, c := range samples {
		binary.WriteComplex64LE(w, c)
	}

	f.Close()
}

func readFromFile(path string) (samples []float64) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %s", err.Error())
	}

	samples = make([]float64, 0)
	b := make([]byte, 2)
	for {
		n, _ := f.Read(b)
		if n < 2 {
			break
		}
		i := int16(b[1]) << 8 | int16(b[0])
		samples = append(samples, float64(i))
	}
	f.Close()

	return samples
}

func readComplexFile(path string, read_sample binary.ReadSampleFunc) (
	samples []complex64) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %s", err.Error())
	}
	r := bufio.NewReader(f)

	samples = make([]complex64, 0)
	for {
		c, err := read_sample(r)
		if err != nil {
			break
		}
		samples = append(samples, c)
	}
	f.Close()

	return samples
}

func computeTable(n int, centre_hz, Δt float64) (r []complex64) {
	const π = math.Pi
	ω1 := - 2 * π * centre_hz
	//σ := 2 * π * 100.0

	r = make([]complex64, n)
	for i, _ := range r {
		t := float64(i) * Δt

		// Notch filter.
		r[i] = complex64(cmplx.Exp(complex(0.0, ω1 * t)))

		/*
		 // Gaussian filter.
		for f := centre_hz - 1000.0; f < centre_hz + 1000.0; f += 1.0 {
			ω0 := 2 * π * f
			g := math.Exp( - (ω0 - ω1)*(ω0 - ω1) / (2*σ*σ))
			r[i] += complex(g, 0) *  cmplx.Exp(complex(0, ω0*t))
		}
		 */
	}

	return r
}

var table_mark, table_space []complex64

func correlateBit(samples []float64, sample_rate,
	markHz, spaceHz float64) bool {
	const π = math.Pi
	Δt := 1.0 / sample_rate

	// TODO(tstranex): This is not thread safe. We should create a
	// correlator struct to store the tables.
	if len(table_mark) != len(samples) || len(table_space) != len(samples) {
		table_mark = computeTable(len(samples), markHz, Δt)
		table_space = computeTable(len(samples), spaceHz, Δt)
	}

	var F_mark complex64
	var F_space complex64
	for i, s := range samples {
		c := complex64(complex(s, 0))
		F_mark += c * table_mark[i]
		F_space += c * table_space[i]
	}

	return cmplx.Abs(complex128(F_mark)) > cmplx.Abs(complex128(F_space))
}

func correlateBitComplex(samples []complex64, sample_rate,
	markHz, spaceHz float64) bool {
	const π = math.Pi
	Δt := 1.0 / sample_rate

	// TODO(tstranex): This is not thread safe. We should create a
	// correlator struct to store the tables.
	if len(table_mark) != len(samples) || len(table_space) != len(samples) {
		table_mark = computeTable(len(samples), markHz, Δt)
		table_space = computeTable(len(samples), spaceHz, Δt)
	}

	var F_mark complex64
	var F_space complex64
	for i, c := range samples {
		F_mark += c * table_mark[i]
		F_space += c * table_space[i]
	}

	return cmplx.Abs(complex128(F_mark)) > cmplx.Abs(complex128(F_space))
}

func DemodulateAFSK1200(samples []float64, sample_rate float64) (
	packets [][]byte) {
	samples_per_bit := sample_rate / baud
	clen := int(samples_per_bit)

	// For clock recovery / bit synchronization, we simply try many
	// different bit streams. At least one of them will be synchronized.
 	bits := make([][]bool, clen)

	excess := 0.0
	for i := clen; i < len(samples); {
		for j := i; j < i + clen; j++ {
			bit := correlateBit(
				samples[j-clen:j], sample_rate, markHz, spaceHz)
			bits[j-i] = append(bits[j-i], bit)
		}
		skip := samples_per_bit + excess
		i += int(skip)
		excess = skip - float64(int(skip))
	}

	// We may detect the same packet multiple times.
	// Store them in a map to remove duplicates.
	packet_set := make(map[string][]byte)

	for i := 0; i < len(bits); i++ {
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

func DemodulateAFSK1200Complex(samples []complex64, sample_rate float64) (
	packets [][]byte) {
	const markHz = 0
	const spaceHz = 1000
	const baud = 1200
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
		}
		skip := samples_per_bit + excess
		i += int(skip)
		excess = skip - float64(int(skip))
	}

	// We may detect the same packet multiple times.
	// Store them in a map to remove duplicates.
	packet_set := make(map[string][]byte)

	for i := 0; i < len(bits); i++ {
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

func main3() {
	rate := 266650.0

	samples := readComplexFile(
		"/Users/tstranex/tmp/1362145622_corrected",
		binary.ReadComplex64LE)
	fmt.Printf("read samples\n")
	
	packets := DemodulateAFSK1200Complex(samples, rate)
	for _, p := range packets {
		//fmt.Printf("Packet %d:\n", i)
		h := strings.ToUpper(hex.EncodeToString(p))
		for i := 0; i < len(h); i += 2 {
			fmt.Printf("%s ", h[i:i+2])
		}
		fmt.Printf("\n")
	}
}

func main() {
	rate := 266910.0
	carrier := 49.0 / 256.0 * rate

	samples := readComplexFile(
		"/Users/tstranex/tmp/strand1_20130304.cut.bin", binary.ReadSampleUINT8)
	fmt.Printf("read samples\n")
	
	packets := DemodulateG3RUHIQ(samples, rate, carrier)
	for _, p := range packets {
		//fmt.Printf("Packet %d:\n", i)
		h := strings.ToUpper(hex.EncodeToString(p))
		for i := 0; i < len(h); i += 2 {
			fmt.Printf("%s ", h[i:i+2])
		}
		fmt.Printf("\n")
	}
}

func main2() {
	rate := 22050.0

	//samples := ModulateAFSK1200(EncodeHDLC([]byte("Hello world!")), rate)
	samples := ModulateG3RUH(EncodeHDLC([]byte("Hello world!")), rate)
	writeToFile("test.raw", samples)

	/*
	samples := readFromFile("2105714597790099704.raw")
	for rate := 22100.0; rate < 22150.0; rate += 1.0 {
		fmt.Printf("rate: %f\n", rate)
		for i, p := range DemodulateAFSK1200(samples, rate) {
			fmt.Printf("Packet %d: %s\n", i, string(p))
			fmt.Printf("%v\n", p)
		}
	}
	 */

	// ~/Downloads/sox-14.4.1/sox -t raw -c 1 -r 22050 -L -b 16 -e signed-integer test.raw test.wav
}