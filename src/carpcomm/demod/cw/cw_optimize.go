package cw

import "fmt"
import "os"
import "log"
import "io"
import "os/exec"
import "bufio"
import "strings"

import "carpcomm/pb"
import "carpcomm/util/binary"

func choose_threshold(a []float64) float64 {
	return 2.0
/*
	s := 0.0
	max := 0.0
	for _, v := range(a) {
		s += v
		if v > max {
			max = v
		}
	}
	mean := s / float64(len(a))
	return 0.5*(mean + max)
*/
}

const (
	DOT = 0
	DASH = 1
	MARK_SPACE = 2
	CHAR_SPACE = 3
	EMPTY = 4

	infty = 100000000
)


func loadMorseTable(path string) (m map[string]string) {
	f, err := os.Open(path)
	if err != nil {
		log.Panicf("Failed to load morse code table from %s: %s",
			path, err.Error())
		return nil
	}

	m = make(map[string]string)
	for {
		var c, code string
		_, err := fmt.Fscanln(f, &c, &code)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panicf("Error reading %s: %s", path, err.Error())
			return nil
		}
		m[code] = c
	}
	return m
}

var _morse_table map[string]string
func LookupMorse(code string) (string, bool) {
	if _morse_table == nil {
		_morse_table = loadMorseTable("src/carpcomm/demod/cw/morse.txt")
	}

	c, ok := _morse_table[code]
	if !ok {
		return "?<" + code + ">", false
	}
	return c, true
}


func err(expected bool, v, th float64) int {
	if (v > th) == expected {
		return 0
	}
	return 1
}

func errs(expected bool, v []float64, th float64) int {
	e := 0
	for _, vv := range v {
		e += err(expected, vv, th)
	}
	return e
}

func min2(values []int, i1, i2 int) (int, int) {
	if values[i1] < values[i2] {
		return values[i1], i1
	}
	return values[i2], i2
}

func min3(values []int, i1, i2, i3 int) (int, int) {
	_, i := min2(values, i1, i2)
	return min2(values, i, i3)
}

func decode_cw(power[] float64, dot_len int) (words []string) {
	dash_len := 3*dot_len
	mark_space_len := dot_len
	char_space_len := dash_len
	word_spacing := 2*char_space_len

	N := len(power)
	th := choose_threshold(power)

	/*
	for i, p := range power {
		if p > th {
			fmt.Printf("*")
		} else {
			fmt.Printf(" ")
		}
		if i % 80 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
	 */

	cost := make([][]int, N+1)
	choice := make([][]int, N+1)
	for i := 0; i <= N; i++ {
		cost[i] = make([]int, 5)
		choice[i] = make([]int, 5)
	}
	cost[N][DOT] = 0
	cost[N][DASH] = 0
	cost[N][MARK_SPACE] = 0
	cost[N][CHAR_SPACE] = 0
	cost[N][EMPTY] = 0

	for i := N-1; i >= 0; i-- {
		// empty
		// followed by either empty, dot or dash
		cost[i][EMPTY], choice[i][EMPTY] = min3(
			cost[i+1], EMPTY, DOT, DASH)
		cost[i][EMPTY] += err(false, power[i], th)

		// char_space
		// followed by either dot, dash, empty
		if i+char_space_len <= N {
			cost[i][CHAR_SPACE], choice[i][CHAR_SPACE] = min3(
				cost[i+char_space_len], DOT, DASH, EMPTY)
			cost[i][CHAR_SPACE] += errs(
				false, power[i:i+char_space_len], th)
		} else {
			cost[i][CHAR_SPACE] = infty
			choice[i][CHAR_SPACE] = infty
		}

		// mark_space
		// followed by either dot or dash
		if i+mark_space_len <= N {
			cost[i][MARK_SPACE], choice[i][MARK_SPACE] = min2(
				cost[i+mark_space_len], DOT, DASH)
			cost[i][MARK_SPACE] += errs(
				false, power[i:i+mark_space_len], th)
		} else {
			cost[i][MARK_SPACE] = infty
			choice[i][MARK_SPACE] = infty
		}

		// dot
		// followed by either mark_space, char_space
		if i+dot_len <= N {
			cost[i][DOT], choice[i][DOT] = min2(
				cost[i+dot_len], MARK_SPACE, CHAR_SPACE)
			cost[i][DOT] += errs(true, power[i:i+dot_len], th)
		} else {
			cost[i][DOT] = infty
			choice[i][DOT] = infty
		}

		// dash
		// followed by either mark_space, char_space
		if i+dash_len <= N {
			cost[i][DASH], choice[i][DASH] = min2(
				cost[i+dash_len], MARK_SPACE, CHAR_SPACE)
			cost[i][DASH] += errs(true, power[i:i+dash_len], th)
		} else {
			cost[i][DASH] = infty
			choice[i][DASH] = infty
		}
	}

	i := 0
	c := EMPTY
	empties := 0
	code := ""
	word := ""
	words = nil
	for {
		if i >= N {
			break
		}
		nc := choice[i][c]

		switch c {
		case EMPTY:
			i++
			empties++
			if nc != EMPTY {
				if empties > word_spacing && len(word) > 0 {
					words = append(words, word)
					word = ""
				}
				empties = 0
			}
		case MARK_SPACE:
			i += mark_space_len
		case CHAR_SPACE:
			i += char_space_len
			//log.Printf("code: %s", code)
			s, ok := LookupMorse(code)
			if ok {
				word += s
			}
			code = ""
		case DOT:
			i += dot_len
			code += "."
		case DASH:
			i += dash_len
			code += "-"
		default:
			log.Printf("panic!\n")
			log.Printf("i = %d\ncost = %d\nc = %d\nnc = %d\n",
				i, cost[i][c], c, nc)
		}

		c = nc
	}

	if len(word) > 0 {
		words = append(words, word)
	}

	return words
}


const cwFilterPath = "src/carpcomm/demod/cw/cw_filter.py"

func DecodeCW(path string,
	sample_rate float64,
	sample_type pb.IQParams_Type,
	cw_params *pb.CWParams) (
	[]pb.Contact_Blob, error) {

	const decimation = 1024 //512

	// 1. CW Filter
	filtered_path := fmt.Sprintf("%s_filtered", path)
	c := exec.Command("python", cwFilterPath,
		path,
		sample_type.String(),
		fmt.Sprintf("%f", sample_rate),
		fmt.Sprintf("%d", decimation),
		filtered_path)
	err := c.Run()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filtered_path)
	if err != nil {
		log.Printf("Error opening %s: %s", filtered_path, err.Error())
		return nil, err
	}
	r := bufio.NewReader(file)
	filtered := make([]float64, 0)
	for {
		f, err := binary.ReadFloat64LE(r)
		if err != nil {
			break
		}
		filtered = append(filtered, f)
	}
	file.Close()

	// Delete temporary file.
	err = os.Remove(filtered_path)
	if err != nil {
		fmt.Printf("Error deleting file: %s", err.Error())
		return nil, err
	}


	// 2. Optimize

	frame_duration_s := float64(decimation) / sample_rate

	log.Printf("%s\n", cw_params)

	dot_len := int(*cw_params.DotDurationS / frame_duration_s + 0.5)

	log.Printf("duration: %f\ndot_len: %d\n", frame_duration_s, dot_len)

	words := decode_cw(filtered, dot_len)
	text := strings.Join(words, " ")

	if text == "" {
		return nil, nil
	}

	blobs := make([]pb.Contact_Blob, 1)
	blobs[0].Format = pb.Contact_Blob_MORSE.Enum()
	blobs[0].InlineData = ([]byte)(text)
	return blobs, nil
}