// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2012 Carpcomm GmbH

// Program to test telemetry decoding with live packets from the network.
//
// Compile:
//   go install github.com/tstranex/carpcomm/telemetry/livetelem
// Run:
//   ./bin/livetelem --station_id=your_id --station_secret=your_secret

package main

import "github.com/tstranex/carpcomm/telemetry"
import "github.com/tstranex/carpcomm/api"
import "fmt"
import "flag"
import "log"

var station_id = flag.String("station_id", "", "Station id")
var station_secret = flag.String("station_secret", "", "Station secret")

func main() {
	// Set your satellite decoder here.
	satellite_id := "hrbe"
	decodeFrame := telemetry.DecodeFrame_hrbe

	flag.Parse()

	c, err := api.NewAPIClient(*station_id, *station_secret)
	if err != nil {
		log.Printf("NewAPIClient error: %s", err.Error())
		return
	}

	packets, err := c.GetLatestPackets(satellite_id, 3)
	if err != nil {
		log.Printf("GetLatestPackets error: %s", err.Error())
		return
	}

	for _, p := range packets {
		fmt.Printf("\nPacket: %v\n", p.Frame)
		data, err := decodeFrame(p.Frame, p.Timestamp)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
		for _, d := range data {
			fmt.Printf("TelemetryDatum: %v\n", &d)
		}
	}
}