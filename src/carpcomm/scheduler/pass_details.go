// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package scheduler

import "os/exec"
import "fmt"
import "time"

type SatPoint struct {
	Timestamp float64
	AzimuthDegrees float64
	AltitudeDegrees float64
	Range float64  // [m]
	RangeVelocity float64  // [m/s]
	LatitudeDegrees float64
	LongitudeDegrees float64
	Elevation float64  // [m]
	IsEclipsed bool
}


const passDetailsBinaryPath = "src/carpcomm/scheduler/pass_details.py"

func passDetails(
	binary_path string,
	begin_time time.Time,
	duration time.Duration,
	latitude_degrees,
	longitude_degrees,
	elevation_metres float64,
	tle string,
	resolution_seconds float64) ([]SatPoint, error) {

	begin_timestamp := 1e-9 * float64(begin_time.UnixNano())

	c := exec.Command("python", binary_path)
	in_pipe, err := c.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer in_pipe.Close()
	out_pipe, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer out_pipe.Close()
	if err := c.Start(); err != nil {
		return nil, err
	}
	defer c.Wait()

	fmt.Fprintf(in_pipe, "%f\n%f\n%f\n%f\n%f\n%s\n%f\n",
		begin_timestamp,
		duration.Seconds(),
		latitude_degrees,
		longitude_degrees,
		elevation_metres,
		tle,
		resolution_seconds)
	in_pipe.Close()

	var n int
	_, err = fmt.Fscanln(out_pipe, &n)
	if err != nil {
		return nil, err
	}

	r := make([]SatPoint, n)
	for i := 0; i < n; i++ {
		p := &r[i]
		fmt.Fscanln(out_pipe,
			&p.Timestamp,
			&p.AzimuthDegrees,
			&p.AltitudeDegrees,
			&p.Range,
			&p.RangeVelocity,
			&p.LatitudeDegrees,
			&p.LongitudeDegrees,
			&p.Elevation,
			&p.IsEclipsed)
	}
	return r, nil
}

func PassDetails(begin_time time.Time,
	duration time.Duration,
	latitude_degrees,
	longitude_degrees,
	elevation_metres float64,
	tle string,
	resolution_seconds float64) ([]SatPoint, error) {

	return passDetails(passDetailsBinaryPath,
		begin_time,
		duration,
		latitude_degrees,
		longitude_degrees,
		elevation_metres,
		tle,
		resolution_seconds)
}