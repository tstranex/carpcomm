// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "carpcomm/pb"
import "strings"

func decodeMorseEachWord(decoded_morse string, timestamp int64,
	f func(decoded_morse string, timestamp int64) (
	[]pb.TelemetryDatum, error)) (
	data []pb.TelemetryDatum, err error) {
	for _, w := range strings.Split(decoded_morse, " ") {
		d, e := f(w, timestamp)
		data = append(data, d...)
		err = e
	}
	if data == nil {
		return nil, err
	}
	return data, nil
}
