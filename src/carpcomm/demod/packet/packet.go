package packet

import "carpcomm/pb"
import "carpcomm/demod/doppler"
import "bufio"
import "strings"
import "encoding/hex"
import "fmt"
import "log"
import "os/exec"
import "os"

const dopplerAnalysisPath = "src/carpcomm/demod/packet/doppler_analysis.py"
const afsk1200LSBPath = "src/carpcomm/demod/packet/afsk1200_lsb.py"
const nbfm9600Path = "src/carpcomm/demod/packet/nbfm9600.py"
const multimonPath = "bin/multimon"

// FIXME: add format param
func DecodePackets(path string,
	sample_rate_hz float64,
	sample_type pb.IQParams_Type,
	c pb.Channel) (
	blobs []pb.Contact_Blob, err error) {

	if c.DopplerStrategy == nil {
		return
	}
	doppler_strategy := c.DopplerStrategy.String()

	if c.Modulation == nil || c.Baud == nil {
		return
	}
	var demod_script string
	var multimon_type string
	if *c.Modulation == pb.Channel_LSB_BFSK &&
		*c.Baud == 1200 {
		demod_script = afsk1200LSBPath
		multimon_type = "AFSK1200"
	} else if *c.Modulation == pb.Channel_FM_GMSK &&
		*c.Baud == 9600 {
		demod_script = nbfm9600Path
		multimon_type = "FSK9600"
	} else {
		return nil, nil
	}

	log.Printf("Running doppler analysis")
	doppler_path := fmt.Sprintf("%s_doppler", path)
	c_analysis := exec.Command("python", dopplerAnalysisPath,
		path, sample_type.String(), doppler_path, doppler_strategy)
	err = c_analysis.Run()
	if err != nil {
		log.Printf("Error running %s: %s",
			dopplerAnalysisPath, err.Error())
		return nil, err
	}

	log.Printf("Running doppler correction")
	corrected_path := fmt.Sprintf("%s_corrected", path)
	err = doppler.ApplyDopplerCorrections(
		path, sample_type, doppler_path, corrected_path)
	if err != nil {
		log.Printf(
			"Error applying doppler corrections: %s", err.Error())
		return nil, err
	}

	log.Printf("Running demodulation")
	wav_path := fmt.Sprintf("%s.wav", path)
	c_demod := exec.Command("python", demod_script,
		corrected_path, fmt.Sprintf("%f", sample_rate_hz), wav_path)
	err = c_demod.Run()
	if err != nil {
		log.Printf("Error running %s: %s", demod_script, err.Error())
		return nil, err
	}

	log.Printf("Running multimon")
	c_decode := exec.Command(multimonPath,
		"-a", multimon_type,
		"-t", "wav",
		wav_path)
	out_pipe, err := c_decode.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := c_decode.Start(); err != nil {
		return nil, err
	}

	r := bufio.NewReader(out_pipe)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		const packetPrefix = "HEXPACKET: "
		if !strings.HasPrefix(line, packetPrefix) {
			continue
		}
		log.Printf("%s", line)

		line = line[len(packetPrefix):len(line)-1]
		frame, err := hex.DecodeString(line)
		if err != nil {
			log.Printf("Error decoding hex packet: %s", err.Error())
		}
		if len(frame) > 0 {
			var blob pb.Contact_Blob
			blob.Format = pb.Contact_Blob_FRAME.Enum()
			blob.InlineData = frame
			blobs = append(blobs, blob)
		}
	}
	out_pipe.Close()
	c_decode.Wait()

	log.Printf("Decoded %d packets.", len(blobs))

	// Delete temporary files.
	err = os.Remove(corrected_path)
	if err != nil {
		log.Printf("Error deleting file: %s", err.Error())
		return blobs, err
	}

	return blobs, nil
}