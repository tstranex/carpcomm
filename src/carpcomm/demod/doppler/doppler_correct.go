package doppler

import "fmt"
import "os"
import "log"
import "io"
import "math"
import "math/cmplx"
import "errors"
import "bufio"
import "carpcomm/util/binary"
import "carpcomm/pb"

type dopplerPair struct {
	sample_num int
	delta_frac float64
}

func readDopplerPair(r io.Reader) (p dopplerPair, err error) {
	n, err := fmt.Fscanf(r, "%d %f", &p.sample_num, &p.delta_frac)
	if n != 2 {
		return dopplerPair{}, errors.New(fmt.Sprintf(
			"Too few values read: %n, %s", n, err.Error()))
	}
	return p, nil
}

func ApplyDopplerCorrections(
	signal_path string,
	sample_type pb.IQParams_Type,
	doppler_path, output_path string) (error) {

	read_sample := binary.GetReadSampleFunc(sample_type)
	if read_sample == nil {
		log.Printf("Invalid sample type")
		return errors.New("Invalid sample type")
	}

	signal_file, err := os.Open(signal_path)
	if err != nil {
		log.Printf("Error opening signal file: %s", err.Error())
		return err
	}
	doppler_file, err := os.Open(doppler_path)
	if err != nil {
		log.Printf("Error opening doppler file: %s", err.Error())
		return err
	}
	output_file, err := os.Create(output_path)
	if err != nil {
		log.Printf("Error opening output file: %s", err.Error())
		return err
	}

	// Buffered io gives a speedup of 6x!
	r := bufio.NewReader(signal_file)
	w := bufio.NewWriter(output_file)

	n := 0

	last_doppler, err := readDopplerPair(doppler_file)
	if err != nil {
		log.Printf("Doppler read error: %s", err.Error())
		return err
	}
	next_doppler, err := readDopplerPair(doppler_file)
	if err != nil {
		log.Printf("Doppler read error: %s", err.Error())
		return err
	}
	has_more_dopplers := true

	for {
		c, err := read_sample(r)
		if err != nil {
			break
		}

		if n > next_doppler.sample_num && has_more_dopplers {
			// We need a new doppler pair.
			d, err := readDopplerPair(doppler_file)
			if err != nil {
				// We've run out of dopplers. Do nothing.
				has_more_dopplers = false
			} else {
				last_doppler = next_doppler
				next_doppler = d
			}
		}

		// exp(-i 2πΔf t)
		frac := last_doppler.delta_frac
		corrector := cmplx.Exp(complex(0.0, -2*math.Pi*frac*float64(n)))
		c = c * complex64(corrector)

		err = binary.WriteComplex64LE(w, c)
		if err != nil {
			log.Printf(
				"Error writing output sample: %s", err.Error())
			return err
		}

		n++
	}

	log.Printf("Doppler corrected %d samples.\n", n)
	signal_file.Close()
	doppler_file.Close()
	output_file.Close()

	return nil
}
