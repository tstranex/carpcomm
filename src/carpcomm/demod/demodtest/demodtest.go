// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "flag"
import "fmt"

import "carpcomm/demod"
import "carpcomm/pb"

var input_file = flag.String("input_file", "", "")
var satellite_id = flag.String("satellite_id", "", "")
var sample_rate = flag.Float64("sample_rate", 266910, "")
var format = flag.String("format", "UINT8", "")

func main() {
	flag.Parse()
	t := (pb.IQParams_Type)(pb.IQParams_Type_value[*format])
	blobs, err := demod.DecodeFromIQ(
		*satellite_id, *input_file, *sample_rate, t)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	for _, b := range blobs {
		fmt.Printf("%s\n", b)
	}
}
