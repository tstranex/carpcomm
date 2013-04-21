package binary

import "io"
import "math"
import "fmt"
import "errors"
import "carpcomm/pb"

func ReadFloat64LE(r io.Reader) (float64, error) {
	b := make([]byte, 8)
	n, err := r.Read(b)
	if n != 8 {
		return 0, errors.New(fmt.Sprintf(
			"Too few bytes read: %s", err.Error()))
	}

	// litte-endian
	var u uint64
	for i := 7; i >= 0; i-- {
		u = (u << 8) | uint64(b[i])
	}
	return math.Float64frombits(u), nil
}

func ReadComplex64LE(r io.Reader) (complex64, error) {
	b := make([]byte, 8)
	n, err := r.Read(b)
	if n != 8 {
		return 0, errors.New(fmt.Sprintf(
			"Too few bytes read: %s", err.Error()))
	}

	// litte-endian
	var rp uint32 = uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 |
		uint32(b[3])<<24
	var ip uint32 = uint32(b[4]) | uint32(b[5])<<8 | uint32(b[6])<<16 |
		uint32(b[7])<<24

	return complex(math.Float32frombits(rp), math.Float32frombits(ip)), nil
}

func WriteComplex64LE(w io.Writer, c complex64) error {
	b := make([]byte, 8)

	// little-endian

	re := math.Float32bits(real(c))
	for i := 0; i < 4; i++ {
		b[i] = byte(re & 255)
		re = re >> 8
	}

	im := math.Float32bits(imag(c))
	for i := 4; i < 8; i++ {
		b[i] = byte(im & 255)
		im = im >> 8
	}

	n, err := w.Write(b)
	if n != 8 {
		return errors.New(fmt.Sprintf(
			"Too few bytes written: %d, %s", n, err.Error()))
	}
	
	return nil
}

func ReadSampleUINT8(r io.Reader) (complex64, error) {
	b := make([]byte, 2)
	n, err := r.Read(b)
	if n != 2 {
		return 0, errors.New(fmt.Sprintf(
			"Too few bytes read: %s", err.Error()))
	}
	re := (float32(b[0]) - 127.0) * 0.008
	im := (float32(b[1]) - 127.0) * 0.008
	return complex(re, im), nil
}

func ReadSampleSINT16(r io.Reader) (complex64, error) {
	b := make([]byte, 4)
	n, err := r.Read(b)
	if n != 4 {
		return 0, errors.New(fmt.Sprintf(
			"Too few bytes read: %s", err.Error()))
	}
	re := float32(int16(b[0]) | (int16(b[1]) << 8)) / 32768.0
	im := float32(int16(b[2]) | (int16(b[3]) << 8)) / 32768.0
	return complex(re, im), nil
}

type ReadSampleFunc func(r io.Reader) (complex64, error)

func GetReadSampleFunc(t pb.IQParams_Type) ReadSampleFunc {
	if t == pb.IQParams_UINT8 {
		return ReadSampleUINT8
	} else if t == pb.IQParams_SINT16 {
		return ReadSampleSINT16
	} else if t == pb.IQParams_FLOAT32 {
		return ReadComplex64LE
	}
	return nil
}
