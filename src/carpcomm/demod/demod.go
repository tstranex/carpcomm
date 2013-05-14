// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package demod

import "carpcomm/pb"
import "carpcomm/db"
import "carpcomm/demod/cw"
import "carpcomm/demod/packet"
import "log"
import "errors"
import "fmt"

func DecodeFromIQ(satellite_id, path string,
	sample_rate_hz float64, sample_type pb.IQParams_Type) (
	blobs []pb.Contact_Blob, err error) {
	sat := db.GlobalSatelliteDB().Map[satellite_id]
	if sat == nil {
		e := errors.New(
			fmt.Sprintf("Unknown satellite_id: %s", satellite_id))
		log.Print(e.Error())
		return nil, e
	}

	// 1. Morse decoding
	var cw_params *pb.CWParams
	for _, c := range sat.Channels {
		if c.Modulation != nil && *c.Modulation == pb.Channel_CW {
			cw_params = c.CwParams
			break
		}
	}
	if cw_params != nil {
		b, err := cw.DecodeCW(
			path, sample_rate_hz, sample_type, cw_params)
		blobs = append(blobs, b...)
		if err != nil {
			log.Printf("Error during DecodeCW: %s", err.Error())
			return blobs, err
		}
	}

	// 2. Frame decoding
	for _, c := range sat.Channels {
		if c.Modulation == nil {
			continue
		}
		if *c.Modulation == pb.Channel_CW {
			// CW is handled above.
			continue
		}
		b, err := packet.DecodePackets(
			path, sample_rate_hz, sample_type, *c)
		blobs = append(blobs, b...)
		if err != nil {
			log.Printf("Error during DecodePackets: %s",
				err.Error())
			return blobs, err
		}
	}

	return blobs, nil
}
