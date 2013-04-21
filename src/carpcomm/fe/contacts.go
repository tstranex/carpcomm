// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import (
	"carpcomm/db"
	"carpcomm/pb"
	"carpcomm/scheduler"
	"log"
	"net/http"
	"html/template"
	"time"
	"encoding/hex"
)

const kStationContactsURLPrefix = "/station/contacts"

func stationContactsUrl(station_id string) string {
	return stationUrl(kStationContactsURLPrefix, station_id)
}

type contactsView struct {
	Title string
	C []*pb.Contact
}

func isFrameBlob(b pb.Contact_Blob) bool {
	if b.Format == nil {
		return false
	}
	return *b.Format == pb.Contact_Blob_FRAME
}

func isMorseBlob(b pb.Contact_Blob) bool {
	if b.Format == nil {
		return false
	}
	return *b.Format == pb.Contact_Blob_MORSE
}

func shouldShowMorse(c pb.Contact) bool {
	// Only show morse messages if there was at least something
	// intelligible.
	// TODO: Consider not adding unintelligable morse messages.
	for _, b := range c.Blob {
		if b.Format != nil && *b.Format == pb.Contact_Blob_DATUM {
			return true
		}
	}
	return false
}

func HasIQBlob(c *pb.Contact) bool {
	for _, b := range c.Blob {
		if b.Format != nil && *b.Format == pb.Contact_Blob_IQ {
			return true
		}
	}
	return false
}

func renderTimestamp(timestamp *int64) string {
	if timestamp == nil {
		return ""
	}
	return time.Unix(*timestamp, 0).UTC().Format(
		ContactTimeFormat)
}

func renderMorse(data []byte) string {
	return (string)(data)
}

func renderPacket(data []byte) string {
	return hex.Dump(data)
}

func satelliteShortName(satellite_id string) string {
	sat := db.GlobalSatelliteDB().Map[satellite_id]
	if sat == nil {
		return ""
	}
	return RenderSatelliteShortName(sat.Name)
}

var contactsFuncMap template.FuncMap = template.FuncMap{
	"HasIQBlob": HasIQBlob,
	"IsFrameBlob": isFrameBlob,
	"IsMorseBlob": isMorseBlob,
	"ShouldShowMorse": shouldShowMorse,
	"RenderMorse": renderMorse,
	"RenderPacket": renderPacket,
	"RenderTimestamp": renderTimestamp,
	"SatelliteViewURL": satelliteViewURL,
	"SatelliteShortName": satelliteShortName}

var contactsTemplate = NewDebuggableTemplate(
	contactsFuncMap,
	"contacts.html",
	"src/carpcomm/fe/templates/contacts.html",
	"src/carpcomm/fe/templates/contact_list.html",
	"src/carpcomm/fe/templates/page.html")

func stationContactsHandler(
	cdb *db.ContactDB, sdb *db.StationDB,
	w http.ResponseWriter, r *http.Request, user userView) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "'id' param missing", http.StatusBadRequest)
		return
	}

	s, err := sdb.Lookup(id)
	if err != nil {
		log.Printf("Sation DB lookup error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if s == nil {
		http.NotFound(w, r)
		return
	}

	if s.Userid == nil || user.Id != *s.Userid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	contacts, err := cdb.SearchByStationId(id, 100)
	if err != nil {
		log.Printf("SearchByStationId error: %s", err.Error());
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var cv contactsView
	cv.Title = *s.Name
	cv.C = contacts

	c := NewRenderContext(user, cv)
	err = contactsTemplate.Get().ExecuteTemplate(w, "contacts.html", c)
	if err != nil {
		log.Printf("Error rendering contacts view: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func contactIQHandler(w http.ResponseWriter, r *http.Request, user userView) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "'id' param missing", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, scheduler.GetStreamURL(id), http.StatusFound)
}

func contactSpectrogramHandler(
	w http.ResponseWriter, r *http.Request, user userView) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "'id' param missing", http.StatusBadRequest)
		return
	}
	url := scheduler.GetStreamURL(id) + ".png"
	http.Redirect(w, r, url, http.StatusFound)
}

func AddContactsHttpHandlers(
	httpmux *http.ServeMux, s *Sessions,
	contactdb *db.ContactDB, stationdb *db.StationDB) {
	HandleFuncLoginRequired(httpmux, kStationContactsURLPrefix, s,
		func(w http.ResponseWriter, r *http.Request, user userView) {
		stationContactsHandler(contactdb, stationdb, w, r, user)
	})
	HandleFuncLoginOptional(httpmux, "/contact/iq", s,
		contactIQHandler)
	HandleFuncLoginOptional(httpmux, "/contact/spectrogram", s,
		contactSpectrogramHandler)
}
