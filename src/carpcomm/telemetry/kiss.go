// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

func DecodeKISS(data []byte) [][]byte {
	const FEND = 0xc0
	const FESC = 0xdb
	const TFEND = 0xdc
	const TFESC = 0xdd
	const SIZE_LIMIT = 8192

	frames := make([][]byte, 0)

	state := 0
	frame := make([]byte, 0)
	for _, b := range data {
		if b == FEND {
			if len(frame) > 1 {
				// Ignore empty frames.
				frames = append(frames, frame[1:])
			}
			state = 0
			frame = nil
			continue
		}

		if state == 0 {
			if b == FESC {
				state = 1
			} else {
				if len(frame) < SIZE_LIMIT {
					frame = append(frame, b)
				}
				state = 0
			}
		} else if state == 1 {
			// We're in escape mode.
			if b == TFEND {
				if len(frame) < SIZE_LIMIT {
					frame = append(frame, FEND)
				}
				state = 0
			} else if b == TFESC {
				if len(frame) < SIZE_LIMIT {
					frame = append(frame, FESC)
				}
				state = 0
			} else {
				// Error, ignore it.
				state = 0
			}
		} else {
			// Unknown state. This shouldn't happen.
			state = 0
		}
	}

	return frames
}