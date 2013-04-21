// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "net/http"
import "net/url"
import "log"
import "encoding/json"
import "encoding/base64"
import "io/ioutil"
import "carpcomm/db"
import "carpcomm/pb"
import "carpcomm/streamer/contacts"
import "strconv"
import "errors"
import "fmt"

type PostPacketRequest struct {
	StationId string `json:"station_id"`
	StationSecret string `json:"station_secret"`
	Timestamp int64 `json:"timestamp"`
	SatelliteId string `json:"satellite_id"`
	Format string `json:"format"`
	FrameBase64 string `json:"frame_base64"`
}

func postPacketHandler(
	sdb *db.StationDB, cdb *db.ContactDB,
	w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Printf("postPacketHandler: Error reading post body: %s",
			err.Error())
		http.Error(w, "Error reading post body",
			http.StatusInternalServerError)
		return
	}

	var req PostPacketRequest
	if err := json.Unmarshal(data, &req); err != nil {
		log.Printf("postPacketHandler: JSON decode error: %s",
			err.Error())
		http.Error(w, "Error decoding JSON data.",
			http.StatusBadRequest)
		return
	}

	log.Printf("request: %v", req)

	frame, err := base64.StdEncoding.DecodeString(req.FrameBase64)
	if err != nil {
		log.Printf("postPacketHandler: base64 decode error: %s",
			err.Error())
		http.Error(w, "Error decoding base64 frame.",
			http.StatusBadRequest)
		return
	}

	if req.StationId == "" {
		http.Error(w, "Missing station_id", http.StatusBadRequest)
		return
	}
	station, err := sdb.Lookup(req.StationId)
	if err != nil {
		log.Printf("Error looking up station: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	contact, poperr := contacts.PopulateContact(
		req.SatelliteId,
		req.Timestamp,
		req.Format,
		frame,
		"",
		req.StationSecret,
		station)
	if poperr != nil {
		poperr.HttpError(w)
		return
	}

	log.Printf("Storing contact: %s", contact)

	err = cdb.Store(contact)
	if err != nil {
		log.Printf("Error storing contact: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// Do we need to set content-length?
}


type GetLatestPacketsRequest struct {
	StationId string
	StationSecret string
	SatelliteId string
	Limit int
}

func parseGetLatestPacketsRequest(values url.Values) (
	req GetLatestPacketsRequest, err error) {
	limit, err := strconv.Atoi(values.Get("limit"))
	if err != nil {
		return req, err
	}
	req.Limit = limit

	req.StationId = values.Get("station_id")
	if req.StationId == "" {
		return req, errors.New("Missing station_id")
	}

	req.StationSecret = values.Get("station_secret")
	if req.StationSecret == "" {
		return req, errors.New("Missing station_secret")
	}

	req.SatelliteId = values.Get("satellite_id")
	if req.SatelliteId == "" {
		return req, errors.New("Missing satellite_id")
	}

	return req, nil
}

type Packet struct {
	Timestamp int64 `json:"timestamp"`
	FrameBase64 string `json:"frame_base64"`
}
type GetLatestPacketsResponse []Packet

func getLatestPacketsHandler(
	sdb *db.StationDB, cdb *db.ContactDB,
	w http.ResponseWriter, r *http.Request) {

	log.Printf("Request: %s", r.URL.String())

	req, err := parseGetLatestPacketsRequest(r.URL.Query())
	if err != nil {
		log.Printf("getLatestPacketsHandler: " +
			"parse request error: %s", err.Error())
		http.Error(w, "Error parsing request.",
			http.StatusBadRequest)
		return
	}

	sat := db.GlobalSatelliteDB().Map[req.SatelliteId]
	if sat == nil {
		http.Error(w, "Unknown satellite_id.", http.StatusBadRequest)
		return
	}

	station, err := sdb.Lookup(req.StationId)
	if err != nil {
		log.Printf("Error looking up station: %s", err.Error())
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	if station == nil {
		log.Printf("Error looking up station.")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	// Authenticate the station.
	if station.Secret == nil || *station.Secret != req.StationSecret {
		log.Printf("Authentication failed.")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	// Make sure that the user is authorized for the satellite.
	found_good_id := false
	for _, station_id := range sat.AuthorizedStationId {
		if station_id == *station.Id {
			found_good_id = true
		}
	}
	if !found_good_id {
		log.Printf("Authentication failed.")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	contacts, err := cdb.SearchBySatelliteId(req.SatelliteId, req.Limit)
	if err != nil {
		log.Printf("getLatestPacketsHandler: " +
			"SearchBySatelliteId error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	packets := make(GetLatestPacketsResponse, 0)

	for _, c := range contacts {
		if c.StartTimestamp == nil {
			continue
		}
		timestamp := *c.StartTimestamp
		for _, b := range c.Blob {
			if len(packets) >= req.Limit {
				continue
			}

			if b.Format == nil ||
				*b.Format != pb.Contact_Blob_FRAME {
				continue
			}

			var p Packet
			p.Timestamp = timestamp
			p.FrameBase64 = base64.StdEncoding.EncodeToString(
				b.InlineData)
			packets = append(packets, p)
		}
	}

	json_body, err := json.Marshal(packets)
	if err != nil {
		log.Printf("json Marshal error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(json_body)))
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(json_body)
	if err != nil {
		log.Printf("Error writing response: %s", err.Error())
	}
}

func AddPacketHttpHandlers(mux *http.ServeMux,
	contactdb *db.ContactDB, 
	stationdb *db.StationDB) {

	mux.HandleFunc("/PostPacket",
		func(w http.ResponseWriter, r *http.Request) {
		postPacketHandler(stationdb, contactdb, w, r)
	})
	mux.HandleFunc("/GetLatestPackets",
		func(w http.ResponseWriter, r *http.Request) {
		getLatestPacketsHandler(stationdb, contactdb, w, r)
	})
}
