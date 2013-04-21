// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package mux

import "net/http"
import "carpcomm/db"
import "sync"
import "log"
import "errors"

type Response struct {
	code int
	data []byte
}

type Request struct {
	request *http.Request  // nil means disconnect immediately
	response chan *Response
}

type Coordinator struct {
	stations map[string] chan Request
	stations_lock sync.RWMutex
	sdb *db.StationDB
}

func NewCoordinator(sdb *db.StationDB) *Coordinator {
	var c Coordinator
	c.stations = make(map[string] chan Request)
	c.sdb = sdb
	return &c
}

func (c *Coordinator) StationCount(
	args *StationCountArgs, reply *StationCountResult) error {
	c.stations_lock.RLock()
	reply.Count = len(c.stations)
	c.stations_lock.RUnlock()
	return nil
}

func (c *Coordinator) StationList(
	args *StationListArgs, result *StationListResult) error {
	result.StationIds = []string{}
	c.stations_lock.RLock()
	for k := range(c.stations) {
		result.StationIds = append(result.StationIds, k)
	}
	c.stations_lock.RUnlock()
	return nil
}

func (c *Coordinator) StationStatus(
	args *StationStatusArgs, result *StationStatusResult) error {
	c.stations_lock.RLock()
	result.IsConnected = (c.stations[args.StationId] != nil)
	c.stations_lock.RUnlock()
	return nil
}

func (c *Coordinator) StationCall(
	args *StationCallArgs, result *StationCallResult) error {

	r, err := http.NewRequest("GET", args.URL, nil)
	if err != nil {
		return err
	}

	resp := c.call(args.StationId, r)
	if resp == nil {
		return errors.New("Station RPC error")
	}

	result.StatusCode = resp.code
	result.Data = resp.data
	return nil
}

// returns nil for rpc errors
func (c *Coordinator) call(station_id string, r *http.Request) *Response {
	c.stations_lock.RLock()
	station, ok := c.stations[station_id]
	c.stations_lock.RUnlock()
	if !ok {
		log.Printf("No such station connected.")
		return nil
	}

	response := make(chan *Response)
	station <- Request{r, response}
	return <-response;
}

func (c *Coordinator) stationConnected(station_id string, input chan Request) {
	// If the station is already connected, disconnect it first.
	c.stations_lock.RLock()
	existing_input, already_connected := c.stations[station_id]
	c.stations_lock.RUnlock()
	if already_connected {
		log.Printf("Disconnecting duplicate existing station.")
		r := Request{nil, make(chan *Response)}
		existing_input <- r
		<- r.response
	}

	c.stations_lock.Lock()
	_, ok := c.stations[station_id]
	if ok {
		log.Printf("Station already connected. This shouldn't happen!")
	}
	c.stations[station_id] = input
	c.stations_lock.Unlock()
}

func (c *Coordinator) stationDisconnected(station_id string) {
	c.stations_lock.Lock()
	delete(c.stations, station_id)
	c.stations_lock.Unlock()
}
