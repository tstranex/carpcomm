// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package mux


type StationCountArgs struct {
}

type StationCountResult struct {
	Count int
}


type StationListArgs struct {
}

type StationListResult struct {
	StationIds []string
}


type StationStatusArgs struct {
	StationId string
}

type StationStatusResult struct {
	IsConnected bool
}


type StationCallArgs struct {
	StationId string
	URL string
}

type StationCallResult struct {
	StatusCode int
	Data []byte
}