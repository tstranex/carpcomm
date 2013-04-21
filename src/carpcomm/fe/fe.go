// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import (
	"carpcomm/db"
	"flag"
	"log"
	"net/http"
	"net/rpc"
)

var mux_address = flag.String("mux_address", ":1235", "Mux address")
var port = flag.String("port", ":8080", "External port")
var db_prefix = flag.String("db_prefix", "r1-", "Database table prefix")
var debug_templates = flag.Bool("debug_templates",
	false, "Enable template debugging")
var debug_auth = flag.Bool("debug_auth", false, "Enable auth on localhost")

type RenderContext struct {
	User userView
	Body   interface{} // per-specific data
}

func NewRenderContext(user userView, body interface{}) *RenderContext {
	return &RenderContext{user, body}
}

func main() {
	flag.Parse()

	mux, err := rpc.DialHTTP("tcp", *mux_address)
	if err != nil {
		log.Fatalf("Mux dial error: %s", err.Error())
	}

	domain, err := db.NewDomain(*db_prefix)
	if err != nil {
		log.Fatalf("Database error: %s", err.Error())
	}

	userdb := domain.NewUserDB()
	stationdb := domain.NewStationDB()
	contactdb := domain.NewContactDB()
	commentdb := domain.NewCommentDB()

	s := NewSessions()

	AddLoginHttpHandlers(s, userdb)
	AddStationHttpHandlers(
		http.DefaultServeMux, s, stationdb, mux, userdb, contactdb)
	AddConsoleHttpHandlers(http.DefaultServeMux,
		s, stationdb, contactdb, mux)
	AddHomeHttpHandlers(http.DefaultServeMux, s, stationdb, mux)
	AddSatelliteHttpHandlers(http.DefaultServeMux, s,
		contactdb, userdb, stationdb, commentdb)
	AddRankingHttpHandlers(http.DefaultServeMux, s,
		contactdb, userdb)
	AddCommentsHttpHandlers(http.DefaultServeMux, s, commentdb)
	AddContactsHttpHandlers(http.DefaultServeMux, s, contactdb, stationdb)
	AddUserHttpHandlers(http.DefaultServeMux, s,
		stationdb, userdb, contactdb)

	log.Printf("fe started.")

	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Printf("ListenAndServe error: %s", err.Error());
	}
}
