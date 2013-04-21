// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "code.google.com/p/goprotobuf/proto"
import "carpcomm/db"
import "carpcomm/pb"
import "net/http"
import "log"
import "html/template"
import "strings"

const userURLPrefix = "/user/"


type profileView struct {
	User *pb.User
	IsOwner bool
	HeardSatellites []*pb.Satellite
	Stations []*pb.Station
}

var userViewTemplate = NewDebuggableTemplate(
	template.FuncMap{
	        "RenderSatelliteName": RenderSatelliteName,
	        "SatelliteViewURL": satelliteViewURL,
        },
	"user.html",
	"src/carpcomm/fe/templates/user.html",
	"src/carpcomm/fe/templates/page.html")

func renderUserProfile(
	cdb *db.ContactDB, stationdb *db.StationDB,
	w http.ResponseWriter, r *http.Request, user userView, u *pb.User) {
	// TODO: It would be better if we could restrict to contacts which
	// have telemetry.
	contacts, err := cdb.SearchByUserId(*u.Id, 100)
	if err != nil {
		log.Printf("cdb.SearchByUserId error: %s", err.Error())
		// Continue since this isn't a critical error.
	}
	heard_satellite_ids := make(map[string]bool)
	for _, c := range contacts {
		if c.SatelliteId == nil {
			continue
		}
		for _, b := range c.Blob {
			if b.Format != nil &&
				*b.Format == pb.Contact_Blob_DATUM {
				heard_satellite_ids[*c.SatelliteId] = true
				break
			}
		}
	}

	var pv profileView
	pv.User = u
	pv.IsOwner = (*u.Id == user.Id)
	pv.HeardSatellites = make([]*pb.Satellite, 0)
	for satellite_id, _ := range heard_satellite_ids {
		pv.HeardSatellites = append(pv.HeardSatellites,
			db.GlobalSatelliteDB().Map[satellite_id])
	}

	stations, err := stationdb.UserStations(*u.Id)
	if err != nil {
		log.Printf("Error getting user stations: %s", err.Error())
		// Continue rendering since it's not a critial error.
	}
	pv.Stations = make([]*pb.Station, 0)
	for _, s := range stations {
		if s.Lat != nil && s.Lng != nil {
			pv.Stations = append(pv.Stations, s)
		}
	}

	c := NewRenderContext(user, pv)
	err = userViewTemplate.Get().ExecuteTemplate(w, "user.html", c)
	if err != nil {
		log.Printf("Error rendering user view: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func changeUserProfile(userdb *db.UserDB,
	w http.ResponseWriter, r *http.Request,
	user userView, u *pb.User) bool {
	log.Printf("hello changeUserProfile")

	if *u.Id != user.Id {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}

	callsign := r.Form.Get("callsign")
	if len(callsign) > 10 {
		callsign = callsign[:10]
	}
	callsign = strings.ToUpper(callsign)
	if callsign != "" {
		u.Callsign = proto.String(callsign)
	} else {
		u.Callsign = nil
	}

	if err := userdb.Store(u); err != nil {
		log.Printf("Error storing user: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return false
	}

	return true
}

func userViewHandler(
	cdb *db.ContactDB, userdb *db.UserDB, stationdb *db.StationDB,
	w http.ResponseWriter, r *http.Request, user userView) {

	if len(r.URL.Path) < len(userURLPrefix) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id := r.URL.Path[len(userURLPrefix):]

	u, err := userdb.Lookup(id)
	if err != nil {
		log.Printf("Error looking up user: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if u == nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" && !changeUserProfile(userdb, w, r, user, u) {
		return
	}
	renderUserProfile(cdb, stationdb, w, r, user, u)
}


func AddUserHttpHandlers(
	httpmux *http.ServeMux, s *Sessions,
	stationdb *db.StationDB,
	userdb *db.UserDB,
	contactdb *db.ContactDB) {
	HandleFuncLoginOptional(httpmux, userURLPrefix, s,
		func(w http.ResponseWriter, r *http.Request, user userView) {
		userViewHandler(contactdb, userdb, stationdb, w, r, user)
	})
}
