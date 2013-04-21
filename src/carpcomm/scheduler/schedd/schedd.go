package main

import (
	"carpcomm/db"
	"carpcomm/scheduler"
	"flag"
	"log"
	"net/rpc"
)

var mux_address = flag.String("mux_address", ":1235", "Mux address")
var db_prefix = flag.String("db_prefix", "r1-", "Database table prefix")

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
	stationdb := domain.NewStationDB()
	contactdb := domain.NewContactDB()

	scheduler.ScheduleForever(stationdb, contactdb, mux)
}