// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

const HDLC_FLAG = 0x7e

// Compute a CRC-16-CCITT checksum as used for AX.25 HDLC.
func ChecksumCRC16HDLC(data []byte) uint16 {
	// The CRC algorithm used for AX.25 is poorly documented.
	// I managed to glean the algorithm from various sources and test it
	// against existing implementations.

	// References:
	// - http://en.wikipedia.org/wiki/Computation_of_CRC
	// - http://www.billnewhall.com/TechDepot/AX25CRC/CRC_for_AX25.pdf
	// - hdlc.c in multimon

	// Polynomial representation of CRC-16-CCITT:
	//   x^16 + x^12 + x^5 + x^0
	// = 0001 0000 0010 0001 (binary)
	// = 0x1021 (hex)
	polynomial := 0x1021

	// The CRC used by AX.25 has some differences to the generic
	// long division:
	// - The shift register is initialized with 0xffff.
	// - The final remainder is inverted.
	// - Bytes are processed LSB first.

	shift_reg := 0xffff  // Initial value (see above).
	for _, b := range data {
		for i := 0; i < 8; i++ {
			// Bytes are transmitted LSB first.
			bit := int(b & 1)
			b = b >> 1

			// Generic long devision.
			shift_reg = shift_reg ^ (bit << 15)
			shift_reg = shift_reg << 1
			if (shift_reg & 0x10000 > 0) {
				shift_reg = shift_reg ^ polynomial
			}
			shift_reg = shift_reg & 0xffff
		}
	}
	shift_reg = shift_reg ^ 0xffff  // Invert result (see above).

	// Reverse bits. This seems to be required by AX.25. However, I don't
	// know why. Possibly it should be done in EncodeHDLC instead.
	reversed := 0
	for i := 0; i < 16; i++ {
		reversed = reversed << 1
		reversed = reversed | (shift_reg & 1)
		shift_reg = shift_reg >> 1
	}

	return uint16(reversed)
}

func byteToBitsLSBFirst(b byte) []bool {
	r := make([]bool, 8)
	for i := 0; i < 8; i++ {
		r[i] = b & 1 > 0
		b = b >> 1
	}
	return r
}

func bitsToByteLSBFirst(bits []bool) (r byte) {
	for i := uint(0); i < 8; i++ {
		if bits[i] {
			r = r | (1 << i)
		}
	}
	return r
}

// Convert a byte buffer into a serial bit stream using HDLC encoding.
// HDLC includes begin and end frame delimiters and a CRC.
func EncodeHDLC(data []byte) (r []bool) {
	crc := ChecksumCRC16HDLC(data)
	data = append(data, byte(crc & 0xff))
	data = append(data, byte((crc >> 8) & 0xff))

	r = append(r, byteToBitsLSBFirst(HDLC_FLAG)...)

	num_ones := 0
	for _, b := range data {
		for _, bit := range byteToBitsLSBFirst(b) {
			r = append(r, bit)
			if bit {
				// Bit stuffing.
				num_ones++
				if num_ones == 5 {
					r = append(r, false)
					num_ones = 0
				}
			} else {
				num_ones = 0
			}

		}
	}

	r = append(r, byteToBitsLSBFirst(HDLC_FLAG)...)

	return r
}

func decodeHDLCBitFrame(bits []bool) []byte {
	if len(bits) < 23 {
		return nil
	}
	bits = bits[:len(bits)-7]  // Remove the ending HDLC flag.
	if len(bits) % 8 != 0 {
		return nil
	}

	data := make([]byte, len(bits)/8)
	for i := 0; i < len(data); i++ {
		data[i] = bitsToByteLSBFirst(bits[i*8:(i+1)*8])
	}

	payload := data[0:len(data)-2]

	crc := ChecksumCRC16HDLC(payload)
	crc_ok := (data[len(data)-2] == byte(crc & 0xff) &&
		data[len(data)-1] == byte((crc >> 8) & 0xff))
	if !crc_ok {
		return nil
	}

	return payload
}

func DecodeHDLC(bits []bool) (r [][]byte) {
	r = make([][]byte, 0)

	stream := 0
	unstuffed := make([]bool, 0)
	num_ones := 0
	for _, b := range bits {
		stream = (stream << 1) & 0xff
		if b {
			stream = stream | 1
		}

		if num_ones != 5 {
			unstuffed = append(unstuffed, b)
		}
		if b {
			num_ones++
		} else {
			num_ones = 0
		}

		if stream == HDLC_FLAG {
			packet := decodeHDLCBitFrame(unstuffed)
			if packet != nil {
				r = append(r, packet)
			}
			unstuffed = make([]bool, 0)
		}
	}

	return r
}