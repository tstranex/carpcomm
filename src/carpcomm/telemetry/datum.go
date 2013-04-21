// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package telemetry

import "carpcomm/pb"
import "code.google.com/p/goprotobuf/proto"
import "testing"
import "math"
import "time"

func NewBoolDatum(key string, timestamp int64, b bool) (d pb.TelemetryDatum) {
	d.Key = proto.String(key)
	d.Timestamp = proto.Int64(timestamp)
	d.Boolean = proto.Bool(b)
	return d
}

func NewInt64Datum(key string, timestamp int64, i int64) (d pb.TelemetryDatum) {
	d.Key = proto.String(key)
	d.Timestamp = proto.Int64(timestamp)
	d.Int64 = proto.Int64(i)
	return d
}

func NewDoubleDatum(key string, timestamp int64, v float64) (
	d pb.TelemetryDatum) {
	d.Key = proto.String(key)
	d.Timestamp = proto.Int64(timestamp)
	d.Double = proto.Float64(v)
	return d
}

func NewIntervalDatum(key string, timestamp int64, minmax [2]float64) (
	d pb.TelemetryDatum) {
	d.Key = proto.String(key)
	d.Timestamp = proto.Int64(timestamp)
	d.IntervalMin = proto.Float64(minmax[0])
	d.IntervalMax = proto.Float64(minmax[1])
	return d
}

func NewTimestampDatum(key string, timestamp int64, t time.Time) (
	d pb.TelemetryDatum) {
	d.Key = proto.String(key)
	d.Timestamp = proto.Int64(timestamp)
	d.UnixTimestamp = proto.Int64(t.Unix())
	return d
}



const testEps = 1e-5

func ExpectBoolDatum(t *testing.T, d pb.TelemetryDatum,
	key string, b bool) {
	if *d.Key != key {
		t.Errorf("Expected key %s, found %s", key, *d.Key)
	}
	if d.Boolean == nil {
		t.Errorf("Expected boolean for %s", key)
		return
	}
	if *d.Boolean != b {
		t.Errorf("Found %t, expected %t for %s", *d.Boolean, b, key)
	}
}

func ExpectInt64Datum(t *testing.T, d pb.TelemetryDatum,
	key string, i int64) {
	if *d.Key != key {
		t.Errorf("Expected key %s, found %s", key, *d.Key)
	}
	if d.Int64 == nil {
		t.Errorf("Expected int64 for %s", key)
		return
	}
	if *d.Int64 != i {
		t.Errorf("Found %d, expected %d for %s", *d.Int64, i, key)
	}
}

func ExpectDoubleDatum(t *testing.T, d pb.TelemetryDatum,
	key string, v float64) {
	if *d.Key != key {
		t.Errorf("Expected key %s, found %s", key, *d.Key)
	}
	if d.Double == nil {
		t.Errorf("Expected double for %s", key)
		return
	}
	if math.Abs(*d.Double - v) > testEps {
		t.Errorf("Found %f, expected %f for %s", *d.Double, v, key)
	}
}

func ExpectTimestampDatum(t *testing.T, d pb.TelemetryDatum,
	key string, ts time.Time) {
	v := ts.Unix()
	if *d.Key != key {
		t.Errorf("Expected key %s, found %s", key, *d.Key)
	}
	if d.UnixTimestamp == nil {
		t.Errorf("Expected unix_timestamp for %s", key)
		return
	}
	if *d.UnixTimestamp != v {
		t.Errorf("Found %d, expected %d for %s",
			*d.UnixTimestamp, v, key)
	}
}

func ExpectIntervalDatum(t *testing.T, d pb.TelemetryDatum,
	key string, min, max float64) {
	if *d.Key != key {
		t.Errorf("Expected key %s, found %s", key, *d.Key)
	}
	if d.IntervalMin == nil || d.IntervalMax == nil {
		t.Errorf("Expected interval for %s", key)
		return
	}
	if math.Abs(*d.IntervalMin - min) > testEps {
		t.Errorf("Found interval_min %f, expected %f for %s",
			*d.IntervalMin, min, key)
	}
	if math.Abs(*d.IntervalMax - max) > testEps {
		t.Errorf("Found interval_max %f, expected %f for %s",
			*d.IntervalMax, max, key)
	}
}
