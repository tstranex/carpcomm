package convert

import "os"
import "log"
import "bufio"
import "carpcomm/util/binary"

func ConvertRTLSDRToComplex64(in_path, out_path string) error {
	in_file, err := os.Open(in_path)
	if err != nil {
		log.Printf("Error opening input file: %s", err.Error())
		return err
	}
	out_file, err := os.Create(out_path)
	if err != nil {
		log.Printf("Error opening output file: %s", err.Error())
		return err
	}

	r := bufio.NewReader(in_file)
	w := bufio.NewWriter(out_file)

	for {
		c, err := binary.ReadSampleUINT8(r)
		if err != nil {
			break
		}

		err = binary.WriteComplex64LE(w, c)
		if err != nil {
			return err
		}
	}

	in_file.Close()
	out_file.Close()

	return nil
}