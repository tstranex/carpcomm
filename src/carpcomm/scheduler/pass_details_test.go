// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package scheduler

import "testing"
import "time"

func pointsEqual(a, b []SatPoint) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestPassDetails(t *testing.T) {
	points, err := passDetails(
		"pass_details.py",
		time.Date(2012, 10, 9, 20, 36, 19, 0.0, time.UTC),
		time.Minute,
		47.4,
		8.5,
		400.0,
		"1998-067CQ\n1 38854U 98067CP  12283.07336473  .00054398  00000-0  89879-3 0    14\n2 38854  51.6473 275.4897 0014651 145.1372 215.0744 15.51582271   671",
		20.0)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	expected_result := []SatPoint{
		{1.349814979e+09, 269.626116, 35.799211, 671746.75, -4939.208496, 46.989871, 1.734424, 403526.46875, true},
		{1.349814999e+09, 276.739676, 44.442965, 574441.5625, -4121.604492, 47.514109, 3.396919, 403230.8125, true},
		{1.349815019e+09, 290.267057, 54.878576, 499263.90625, -2836.593506, 48.011316, 5.093884, 402940.5625, true}}

	if !pointsEqual(points, expected_result) {
		t.Errorf("Unexpected result: %v\nexpected: %v",
			points, expected_result)
	}
}