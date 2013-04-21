// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package scheduler

import "testing"
import "time"

func predictionsEqual(a, b []Prediction) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}

var timeWithTwoPasses = time.Date(2012, 6, 4, 8, 53, 18, 263081100, time.UTC)
var timeWithOneHighAltitudePass = time.Date(
	2012, 6, 11, 7, 33, 18, 263081100, time.UTC)

func testPrediction(begin_time time.Time) ([]Prediction, error) {
	return predict(
		"predict.py",
		begin_time,
		12*time.Hour,
		47.4,
		8.5,
		400.0,
		20.0,
		270.0,
		90.0,
		"SWISSCUBE               \n1 35932U 09051B   12110.66765508  .00000638  00000-0  15500-3 0  5172\n2 35932  98.3348 213.8703 0006768 284.4795  75.6141 14.52927878136365")
}

func TestPredictTwoPasses(t *testing.T) {
	p, err := testPrediction(timeWithTwoPasses)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	expected_result := []Prediction{
		{nil, 1.338810064e+09, 1.338810216e+09,
			41.776832, 90.248273, 36.405011,
			CompatibleMode{}},
		{nil, 1.338815964e+09, 1.338816173e+09,
			339.350162, 269.772146, 30.092478,
			CompatibleMode{}}}

	if !predictionsEqual(p, expected_result) {
		t.Errorf("Unexpected result: %v\nexpected: %v",
			p, expected_result)
	}
}

func TestPredictOneHighAltitudePass(t *testing.T) {
	p, err := testPrediction(timeWithOneHighAltitudePass)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	expected_result := []Prediction{
		{nil, 1.339416896e+09, 1.339417092e+09,
			19.304741, 90.143757, 75.378439,
			CompatibleMode{}}}

	if !predictionsEqual(p, expected_result) {
		t.Errorf("Unexpected result: %v\nexpected: %v",
			p, expected_result)
	}
}

func TestPredictRegression1(t *testing.T) {
	p, err := predict(
		"predict.py",
		time.Date(2012, 6, 19, 20, 7, 52, 0, time.UTC),
		12*time.Hour,
		40.134000,
		-74.038000,
		10.000000,
		10.000000,
		0.000000,
		360.000000,
		"NOAA 18 [+]\n1 28654U 05018A   12152.90553951  .00000282  00000-0  15477-3 0  1987\n2 28654  99.0455 110.8604 0013833 197.4644 162.6526 14.11656240362269")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	expected_result := []Prediction{
		{nil, 1.340136472e+09, 1.340136652e+09,
			296.244893, 324.210229, 24.568984,
			CompatibleMode{}},
		{nil, 1.340174661e+09, 1.340175044e+09,
			55.216701, 127.326583, 15.555426,
			CompatibleMode{}}}

	if !predictionsEqual(p, expected_result) {
		t.Errorf("Unexpected result: %v\nexpected: %v",
			p, expected_result)
	}
}

// 2012-06-02 02:16: 45120960 ns/op
func BenchmarkPredict(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testPrediction(timeWithTwoPasses)
	}
}