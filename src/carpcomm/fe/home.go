// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import (
	"carpcomm/db"
	"carpcomm/mux"
	"carpcomm/pb"
	"log"
	"net/http"
	"net/rpc"
)

type homeContext struct {
	NumStations       int
	NumOnlineStations int
	Stations          []*pb.Station
}

var homeTemplate = NewDebuggableTemplate(
	nil,
	"home.html",
	"src/carpcomm/fe/templates/home.html",
	"src/carpcomm/fe/templates/page.html")

func homeHandler(stationdb *db.StationDB, m *rpc.Client,
	w http.ResponseWriter, r *http.Request, user userView) {
	log.Printf("userid: %s", user.Id)

	hc := homeContext{}

	n, err := stationdb.NumStations()
	if err != nil {
		n = -1
	}
	hc.NumStations = n

	hc.Stations, err = stationdb.UserStations(user.Id)
	if err != nil {
		log.Printf("Error getting user stations: %s", err.Error())
		// Continue rendering since it's not a critial error.
	}

	var args mux.StationCountArgs
	var count mux.StationCountResult
	err = m.Call("Coordinator.StationCount", args, &count)
	if err != nil {
		count.Count = -1
	}
	hc.NumOnlineStations = count.Count

	c := NewRenderContext(user, hc)
	err = homeTemplate.Get().ExecuteTemplate(w, "home.html", c)
	if err != nil {
		log.Printf("Error rendering home page: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func AddHomeHttpHandlers(httpmux *http.ServeMux, s *Sessions,
	stationdb *db.StationDB, mux *rpc.Client) {
	HandleFuncLoginRequired(httpmux, "/home", s,
		func(w http.ResponseWriter, r *http.Request, user userView) {
			homeHandler(stationdb, mux, w, r, user)
		})
}
