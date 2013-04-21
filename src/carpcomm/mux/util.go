// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package mux

import "net/rpc"
import "net/url"
import "log"
import "errors"
import "fmt"
import "encoding/json"
import "time"
import "carpcomm/util"

const stationCallRPC = "Coordinator.StationCall"

const stationReceiverSetFrequency = "ReceiverSetFrequency"
const stationReceiverStart = "ReceiverStart"
const stationReceiverStop = "ReceiverStop"

const stationTNCStart = "TNCStart"
const stationTNCStop = "TNCStop"

const stationMotorStart = "MotorStart"
const stationMotorStop = "MotorStop"

const stationStatusCodeOk = 200


func CallStation(mux_client *rpc.Client,
	station_id string, action string, params url.Values) (
	StationCallResult, error) {

	var args StationCallArgs
	args.StationId = station_id
	u := url.URL{}
	u.Path = "/" + action
	if params != nil {
		u.RawQuery = params.Encode()
	}
	args.URL = u.String()

	var result StationCallResult
	err := mux_client.Call(stationCallRPC, args, &result)
	if err != nil {
		log.Printf("StationCall error: %s", err.Error())
		return result, err
	}

	return result, nil
}

func CallStationAndCheckStatus(mux_client *rpc.Client,
	station_id string, action string, params url.Values) error {
	result, err := CallStation(mux_client, station_id, action, params)
	if err != nil {
		return err
	}
	if result.StatusCode != stationStatusCodeOk {
		return errors.New(fmt.Sprintf(
			"Station returned error code: %d", result.StatusCode))
	}
	return nil
}


func StationReceiverSetFrequency(mux_client *rpc.Client,
	station_id string,
	freq_hz int64) error {

	params := url.Values{}
	params.Add("hz", fmt.Sprintf("%d", freq_hz))
	return CallStationAndCheckStatus(
		mux_client, station_id, stationReceiverSetFrequency, params)
}

func StationReceiverStart(
	mux_client *rpc.Client, station_id, stream_url string) error {

	params := url.Values{}
	params.Add("stream_url", stream_url)
	return CallStationAndCheckStatus(
		mux_client, station_id, stationReceiverStart, params)
}

func StationReceiverStop(mux_client *rpc.Client, station_id string) error {
	return CallStationAndCheckStatus(
		mux_client, station_id, stationReceiverStop, nil)
}

func StationTNCStart(
	mux_client *rpc.Client, station_id string,
	api_server string,
	satellite_id string) error {

	host, port, err := util.SplitHostAndPort(api_server)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("api_host", host)
	params.Add("api_port", port)
	params.Add("satellite_id", satellite_id)

	return CallStationAndCheckStatus(
		mux_client, station_id, stationTNCStart, params)
}

func StationTNCStop(mux_client *rpc.Client, station_id string) error {
	return CallStationAndCheckStatus(
		mux_client, station_id, stationTNCStop, nil)
}

type MotorCoordinate struct {
	Timestamp float64
	AzimuthDegrees float64
	AltitudeDegrees float64
}

func StationMotorStart(mux_client *rpc.Client, station_id string,
	program []MotorCoordinate) error {

	if len(program) == 0 {
		return errors.New("Empty motor program.")
	}
	coords := make([][3]float64, len(program))
	start_t := float64(time.Now().Unix())
	for i, c := range program {
		coords[i][0] = c.Timestamp - start_t
		coords[i][1] = c.AzimuthDegrees
		coords[i][2] = c.AltitudeDegrees
	}

	p, err := json.Marshal(coords)
	if err != nil {
		log.Printf("Error json marshalling motor program: %s",
			err.Error())
		return err
	}

	params := url.Values{}
	params.Add("program", (string)(p))

	return CallStationAndCheckStatus(
		mux_client, station_id, stationMotorStart, params)
}

func StationMotorStop(mux_client *rpc.Client, station_id string) error {
	return CallStationAndCheckStatus(
		mux_client, station_id, stationMotorStop, nil)
}
