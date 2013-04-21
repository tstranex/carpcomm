package timestamp

import "time"

func TimeToTimestampFloat(t time.Time) float64 {
	// TODO: handle nanoseconds
	return (float64)(t.Unix())
}

func TimestampFloatToTime(timestamp float64) time.Time {
	sec := (int64)(timestamp)
	nsec := (int64)(1e9 * (timestamp - (float64)(sec)))
	return time.Unix(sec, nsec)
}
