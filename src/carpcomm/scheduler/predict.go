// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package scheduler

import "carpcomm/db"
import "carpcomm/pb"
import "os/exec"
import "fmt"
import "sort"
import "time"
import "log"

type CompatibleMode struct {
	channel pb.Channel
	limits pb.AzElLimits
}

type Prediction struct {
	Satellite *pb.Satellite
	StartTimestamp, EndTimestamp float64
	StartAzimuthDegrees float64
	EndAzimuthDegrees float64
        MaxAltitudeDegrees float64
	CompatibleMode
}

func (a Prediction) Equals(b Prediction) bool {
	return a.Satellite == b.Satellite &&
		a.StartTimestamp == b.StartTimestamp &&
		a.EndTimestamp == b.EndTimestamp &&
		a.StartAzimuthDegrees == b.StartAzimuthDegrees &&
		a.EndAzimuthDegrees == b.EndAzimuthDegrees &&
		a.MaxAltitudeDegrees == b.MaxAltitudeDegrees
}

type PredictionList []Prediction

func (p PredictionList) Len() int {
	return len(p)
}
func (p PredictionList) Less(i, j int) bool {
	return p[i].StartTimestamp < p[j].StartTimestamp
}
func (p PredictionList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}


func Duration(seconds float64) time.Duration {
	return (time.Duration)(seconds * (float64)(time.Second))
}


const predictBinaryPath = "src/carpcomm/scheduler/predict.py"

func predict(
	predict_binary_path string,
	begin_time time.Time,
	duration time.Duration,
	latitude_degrees,
	longitude_degrees,
	elevation_metres,
	min_altitude_degrees,
	min_azimuth_degrees,
	max_azimuth_degrees float64,
	tle string) ([]Prediction, error) {

	begin_timestamp := 1e-9 * float64(begin_time.UnixNano())

	c := exec.Command("python", predict_binary_path)
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

	fmt.Fprintf(in_pipe, "%f\n%f\n%f\n%f\n%f\n%f\n%f\n%f\n%s\n",
		begin_timestamp,
		duration.Seconds(),
		latitude_degrees,
		longitude_degrees,
		elevation_metres,
		min_altitude_degrees,
		min_azimuth_degrees,
		max_azimuth_degrees,
		tle)
	in_pipe.Close()

	var n int
	_, err = fmt.Fscanln(out_pipe, &n)
	if err != nil {
		return nil, err
	}

	r := make([]Prediction, n)
	for i := 0; i < n; i++ {
		p := &r[i]
		fmt.Fscanln(out_pipe,
			&p.StartTimestamp,
			&p.EndTimestamp,
			&p.StartAzimuthDegrees,
			&p.EndAzimuthDegrees,
			&p.MaxAltitudeDegrees)
	}
	return r, nil
}

func Is2mBand(hz float64) bool {
	return 100e6 <= hz && hz < 200e6
}
func Is70cmBand(hz float64) bool {
	return 400e6 <= hz && hz < 500e6
}

func isValidAzElLimits(l *pb.AzElLimits) bool {
	if l == nil ||
		l.MinElevationDegrees == nil ||
		l.MinAzimuthDegrees == nil ||
		l.MaxAzimuthDegrees == nil {
		return false
	}
	return true
}

func CompatibleModes(station *pb.Station, sat *pb.Satellite) []CompatibleMode {
	if station.Capabilities == nil {
		return nil
	}

	var modes []CompatibleMode

	for _, c := range sat.Channels {
		if c.Downlink == nil || *c.Downlink == false {
			continue
		}
		if Is2mBand(*c.FrequencyHz) {
			if isValidAzElLimits(station.Capabilities.VhfLimits) {
				modes = append(modes, CompatibleMode{
					*c, *station.Capabilities.VhfLimits})
			}
		} else if Is70cmBand(*c.FrequencyHz) {
			if isValidAzElLimits(station.Capabilities.UhfLimits) {
				modes = append(modes, CompatibleMode{
					*c, *station.Capabilities.UhfLimits})
			}
		}
	}

	return modes
}

func PassPredictions(station *pb.Station) (PredictionList, error) {
	db := db.GlobalSatelliteDB()

	if station.Lat == nil || station.Lng == nil {
		return nil, nil
	}

	lat := *station.Lat
	lng := *station.Lng
	elevation := 0.0
	if station.Elevation != nil {
		elevation = *station.Elevation
	}

	begin_time := time.Now()

	var predictErr error = nil

	var all_passes PredictionList
	for _, sat := range db.List {
		if sat.Tle == nil {
			continue
		}

		modes := CompatibleModes(station, sat)
		for _, mode := range modes {

			passes, err := predict(
				predictBinaryPath,
				begin_time,
				18*time.Hour,
				lat, lng, elevation,
				*mode.limits.MinElevationDegrees,
				*mode.limits.MinAzimuthDegrees,
				*mode.limits.MaxAzimuthDegrees,
				*sat.Tle)
			if err != nil {
				predictErr = err
				log.Printf("Prediction error: %s", err.Error())
				// This may be a temporary error so don't stop
				// immediately.
			}

			for i, _ := range passes {
				passes[i].Satellite = sat
				passes[i].CompatibleMode = mode
			}
			all_passes = append(all_passes, passes...)
		}
	}
	sort.Sort(all_passes)

	return all_passes, predictErr
}
