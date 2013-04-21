// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "log"
import "fmt"
import "carpcomm/db"
import "carpcomm/demod"
import "carpcomm/pb"
import "carpcomm/streamer/contacts"

func processIQDataForSatellite(
	satellite_id string,
	local_path string,
	iq_params pb.IQParams,
	timestamp int64) (result []*pb.Contact_Blob) {

	log.Printf("Processing new IQ data %s for %s", local_path, satellite_id)

	blobs, err := demod.DecodeFromIQ(
		satellite_id, local_path,
		(float64)(*iq_params.SampleRate), *iq_params.Type)
	if err != nil {
		log.Printf("Error while processing IQ data: %s", err.Error())
		// Don't exit since some blobs may have been generated anyway.
	}
	log.Printf("Decoded %d blobs.", len(blobs))
	for _, b := range blobs {
		var cb pb.Contact_Blob = b
		result = append(result, &cb)
	}

	decoded_blobs, err := contacts.DecodeBlobs(
		satellite_id, timestamp, blobs)
	if err != nil {
		log.Printf("Error decoding blobs: %s", err.Error())
	}
	log.Printf("Decoded %d telemetry blobs.", len(decoded_blobs))
	result = append(result, decoded_blobs...)

	return result
}

// Consider moving this to a completely different worker binary.
func processNewIQData(contact_id string, contactdb *db.ContactDB) {
	log.Printf("%s: Processing IQ data", contact_id)

	contact, err := contactdb.Lookup(contact_id)
	if err != nil {
		log.Printf("%s: Error looking up contact: %s",
			contact_id, err.Error())
		return
	}
	if contact == nil {
		log.Printf("%s: Contact not found.", contact_id)
		return
	}

	// Get IQParams from the contact.
	var iq_params *pb.IQParams
	for _, b := range contact.Blob {
		if b.Format != nil && *b.Format == pb.Contact_Blob_IQ {
			if b.IqParams != nil {
				iq_params = b.IqParams
			}
		}
	}
	if iq_params == nil {
		log.Printf("%s: IQ blob missing", contact_id)
		return
	}

	local_path := fmt.Sprintf("%s/%s", *stream_tmp_dir, contact_id)

	png_path := fmt.Sprintf("%s.png", local_path)
	demod.Spectrogram(local_path, *iq_params, png_path,
		demod.SpectrogramTitle(*contact, *iq_params))

	if contact.SatelliteId != nil {
		blobs := processIQDataForSatellite(
			*contact.SatelliteId,
			local_path,
			*iq_params,
			*contact.StartTimestamp)
		contact.Blob = append(contact.Blob, blobs...)
	}

	log.Printf("%s: Storing updated contact: %v", contact_id, contact.Blob)
	
	// TODO: Need to be careful about locking the ContactDB record.
	// Currently we are the only process that modifies an existing ContactDB
	// record. In the future we may need to be more careful.
	if err := contactdb.Store(contact); err != nil {
		log.Printf("%s: Error storing contact: %s",
			contact_id, err.Error())
		return
	}
	log.Printf("%s: Wrote updated contact to db.", contact_id)
}

type IQProcessingQueue chan string

func ProcessIQQueue(queue IQProcessingQueue, contactdb *db.ContactDB) {
	log.Printf("Starting IQ processing queue")
	for {
		contact_id := <-queue
		processNewIQData(contact_id, contactdb)
	}
}
