package main

import "carpcomm/db"
import "carpcomm/mux"
import "net/http"
import "net/rpc"
import "log"
import "flag"

var cert_file = flag.String(
	"cert_file",
	"/Users/tstranex/carp/carp/frontend/cert.pem",
	"SSL server certificate file")
var private_key_file = flag.String(
	"private_key_file",
	"/Users/tstranex/carp/carp/frontend/private_key.pem",
	"SSL private key file")

var rpc_port = flag.String(
	"rpc_port", ":1235", "Internal RPC port")
var external_port = flag.String(
	"external_port", ":1234", "External mux port")
var db_prefix = flag.String("db_prefix", "r1-", "Database table prefix")

func main() {
	flag.Parse()
	log.Printf("Starting multiplexer.")

	domain, err := db.NewDomain(*db_prefix)
	if err != nil {
		log.Fatalf("Database error: %s", err.Error())
	}
	stationdb := domain.NewStationDB()
	c := mux.NewCoordinator(stationdb)

	go mux.ListenAndServe(c, *cert_file, *private_key_file, *external_port)

	rpc.Register(c)
	rpc.HandleHTTP()
	http.ListenAndServe(*rpc_port, nil)
}