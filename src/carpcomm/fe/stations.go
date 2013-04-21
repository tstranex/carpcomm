// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import (
	"carpcomm/db"
	"carpcomm/mux"
	"carpcomm/pb"
	"carpcomm/scheduler"
	"carpcomm/util/timestamp"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/rpc"
	"net/url"
	"reflect"
	"strconv"
	texttemplate "text/template"
	"time"
)

import "code.google.com/p/goprotobuf/proto"

func stationUrl(handler, id string) string {
	v := &url.Values{}
	v.Set("id", id)
	url := url.URL{}
	url.Path = handler
	url.RawQuery = v.Encode()
	return url.String()
}

var kmlTemplate = texttemplate.Must(texttemplate.ParseFiles(
	"src/carpcomm/fe/templates/stations.kml"))

func stationKMLHandler(sdb *db.StationDB,
	w http.ResponseWriter, r *http.Request) {
	stations, err := sdb.AllStations()
	if err != nil {
		log.Printf("Station DB AllStations error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = kmlTemplate.ExecuteTemplate(w, "stations.kml", stations)
	if err != nil {
		log.Printf("Error rendering station kml map: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func addStationHandler(sdb *db.StationDB,
	w http.ResponseWriter, r *http.Request, user userView) {
	station, err := db.NewStation(user.Id)
	if err != nil {
		log.Printf("Station DB NewStation error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err = sdb.Store(station); err != nil {
		log.Printf("Station DB Store error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	redirect := stationUrl("/station/edit", *station.Id)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func deleteStationHandler(sdb *db.StationDB,
	w http.ResponseWriter, r *http.Request, user userView) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "'id' param missing", http.StatusBadRequest)
		return
	}

	s, err := sdb.Lookup(id)
	if err != nil {
		log.Printf("Station DB lookup error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if s == nil || *s.Userid != user.Id {
		http.NotFound(w, r)
		return
	}

	err = sdb.Delete(id)
	if err != nil {
		log.Printf("Station DB delete error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Success
	http.Redirect(w, r, "/home", http.StatusFound)
}

var stationsFuncMap template.FuncMap = template.FuncMap{
	"derefbool": func(p interface{}) bool {
		v := reflect.ValueOf(p)
		if v.IsNil() {
			return false
		}
		return v.Elem().Bool()
	},
	"roundn": func(n int, f *float64) string {
		if f == nil {
			return ""
		}
		format := fmt.Sprintf("%%.%df", n)
		return fmt.Sprintf(format, *f)
	},
	"derefstr": func(p interface{}) string {
		v := reflect.ValueOf(p)
		if v.IsNil() {
			return ""
		}
		return fmt.Sprintf("%v", v.Elem())
	},
	"SatelliteViewURL": satelliteViewURL,
}

func unionFuncMap(f ...template.FuncMap) (result template.FuncMap) {
	result = make(template.FuncMap)
	for _, m := range f {
		for k, v := range m {
			// It would be nice if we could panic if we find
			// conflicting values. However, it's impossible to
			// function types.
			result[k] = v
		}
	}
	return result
}

type predictionView struct {
	SatelliteId                           string
	SatelliteName                         string
	StartTime, Duration                   string
	StartAzimuth, EndAzimuth, MaxAltitude float64
}

const TimeFormat = "02 Jan 15:04 MST"

func FillPredictionView(p scheduler.Prediction, pv *predictionView) {
	pv.SatelliteId = *p.Satellite.Id
	pv.SatelliteName = RenderSatelliteName(p.Satellite.Name)
	pv.StartTime = timestamp.TimestampFloatToTime(
		p.StartTimestamp).UTC().Format(TimeFormat)
	pv.Duration = scheduler.Duration(
		p.EndTimestamp - p.StartTimestamp).String()
	pv.StartAzimuth = p.StartAzimuthDegrees
	pv.EndAzimuth = p.EndAzimuthDegrees
	pv.MaxAltitude = p.MaxAltitudeDegrees
}

func LookupUserView(userdb *db.UserDB, userid string) (v userView) {
	v.Id = userid

	u, _ := userdb.Lookup(userid)
	if u == nil {
		return v
	}

	if u.DisplayName != nil {
		v.Name = *u.DisplayName
		if u.Callsign != nil {
			v.Name += " " + *u.Callsign
		}
	}
	if u.PhotoUrl != nil {
		v.PhotoUrl = *u.PhotoUrl
	}

	return v
}


type stationContext struct {
	S             *pb.Station
	IsOwner       bool
	IsOnline      bool
	CurrentTime   string
	NextPasses    []predictionView
	Operator      userView
	Contacts []*pb.Contact
}

func GetStationContext(s *pb.Station, userid string,
	m *rpc.Client, userdb *db.UserDB, cdb *db.ContactDB) (
	sc stationContext) {
	sc.S = s
	sc.IsOwner = (userid == *s.Userid)

	var args mux.StationStatusArgs
	args.StationId = *s.Id
	var status mux.StationStatusResult
	err := m.Call("Coordinator.StationStatus", args, &status)
	if err != nil {
		sc.IsOnline = false
	} else {
		sc.IsOnline = status.IsConnected
	}

	sc.CurrentTime = time.Now().UTC().Format(TimeFormat)
	passes, _ := scheduler.PassPredictions(s)
	passViews := make([]predictionView, 0)
	for _, pass := range passes {
		var pv predictionView
		FillPredictionView(pass, &pv)
		// Don't display different modes of the same satellite.
		if len(passViews) == 0 || passViews[len(passViews)-1] != pv {
			passViews = append(passViews, pv)
		}
	}
	sc.NextPasses = passViews

	sc.Operator = LookupUserView(userdb, *s.Userid)


	if sc.IsOwner {
		c, err := cdb.SearchByStationId(*s.Id, 5)
		if err != nil {
			log.Printf("SearchByStationId error: %s", err.Error());
			// This is not a fatal error.
		}
		sc.Contacts = c
	}

	return sc
}

var stationViewTemplate = NewDebuggableTemplate(
	unionFuncMap(stationsFuncMap, contactsFuncMap),
	"station.html",
	"src/carpcomm/fe/templates/station.html",
	"src/carpcomm/fe/templates/contact_list.html",
	"src/carpcomm/fe/templates/page.html")

func stationViewHandler(sdb *db.StationDB, m *rpc.Client, userdb *db.UserDB,
	contactdb *db.ContactDB,
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

	if s.Capabilities == nil {
		s.Capabilities = &pb.Capabilities{}
	}

	sc := GetStationContext(s, user.Id, m, userdb, contactdb)

	c := NewRenderContext(user, sc)
	err = stationViewTemplate.Get().ExecuteTemplate(w, "station.html", c)
	if err != nil {
		log.Printf("Error rendering station view: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

var stationEditTemplate = NewDebuggableTemplate(
	stationsFuncMap,
	"edit_station.html",
	"src/carpcomm/fe/templates/edit_station.html",
	"src/carpcomm/fe/templates/page.html")

func editStationGET(s *pb.Station,
	w http.ResponseWriter, r *http.Request, user userView) {
	c := NewRenderContext(user, s)
	err := stationEditTemplate.Get().ExecuteTemplate(
		w, "edit_station.html", c)
	if err != nil {
		log.Printf("Error rendering station edit page: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func setOptionalFloat(s string, v **float64) {
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		*v = proto.Float64(f)
	} else {
		*v = nil
	}
}

func editStationPOST(sdb *db.StationDB, s *pb.Station,
	w http.ResponseWriter, r *http.Request, user userView) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print(r.Form)

	s.Name = proto.String(r.Form.Get("name"))
	s.Notes = proto.String(r.Form.Get("notes"))
	setOptionalFloat(r.Form.Get("latitude"), &s.Lat)
	setOptionalFloat(r.Form.Get("longitude"), &s.Lng)
	setOptionalFloat(r.Form.Get("elevation"), &s.Elevation)

	if r.Form["has_vhf"] == nil {
		s.Capabilities.VhfLimits = nil
	} else {
		vhf := &pb.AzElLimits{}
		s.Capabilities.VhfLimits = vhf
		setOptionalFloat(r.Form.Get("vhf_min_azimuth"),
			&vhf.MinAzimuthDegrees)
		setOptionalFloat(r.Form.Get("vhf_max_azimuth"),
			&vhf.MaxAzimuthDegrees)
		setOptionalFloat(r.Form.Get("vhf_min_elevation"),
			&vhf.MinElevationDegrees)
		setOptionalFloat(r.Form.Get("vhf_max_elevation"),
			&vhf.MaxElevationDegrees)
	}

	if r.Form["has_uhf"] == nil {
		s.Capabilities.UhfLimits = nil
	} else {
		uhf := &pb.AzElLimits{}
		s.Capabilities.UhfLimits = uhf
		setOptionalFloat(r.Form.Get("uhf_min_azimuth"),
			&uhf.MinAzimuthDegrees)
		setOptionalFloat(r.Form.Get("uhf_max_azimuth"),
			&uhf.MaxAzimuthDegrees)
		setOptionalFloat(r.Form.Get("uhf_min_elevation"),
			&uhf.MinElevationDegrees)
		setOptionalFloat(r.Form.Get("uhf_max_elevation"),
			&uhf.MaxElevationDegrees)
	}

	if r.Form["scheduler_enabled"] == nil {
		s.SchedulerEnabled = nil
	} else {
		s.SchedulerEnabled = proto.Bool(true)
	}

	err := sdb.Store(s)
	if err != nil {
		log.Printf("Station DB store error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	redirect := stationUrl("/station", *s.Id)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func editStationHandler(sdb *db.StationDB,
	w http.ResponseWriter, r *http.Request, user userView) {

	log.Printf("method: %s", r.Method)

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "'id' param missing", http.StatusBadRequest)
		return
	}

	s, err := sdb.Lookup(id)
	if err != nil {
		log.Printf("Station DB lookup error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if s == nil || *s.Userid != user.Id {
		http.NotFound(w, r)
		return
	}

	if s.Capabilities == nil {
		s.Capabilities = &pb.Capabilities{}
	}

	if r.Method == "POST" {
		editStationPOST(sdb, s, w, r, user)
	} else {
		editStationGET(s, w, r, user)
	}
}

type StationHandlerFunc func(
	*db.StationDB, http.ResponseWriter, *http.Request, userView)

func stationHandler(stationdb *db.StationDB,
	f StationHandlerFunc) HandlerFuncWithUser {
	return func(w http.ResponseWriter, r *http.Request, user userView) {
		f(stationdb, w, r, user)
	}
}

func AddStationHttpHandlers(
	httpmux *http.ServeMux, s *Sessions,
	stationdb *db.StationDB,
	mux *rpc.Client,
	userdb *db.UserDB,
	contactdb *db.ContactDB) {
	HandleFuncLoginRequired(httpmux, "/station/add", s,
		stationHandler(stationdb, addStationHandler))
	HandleFuncLoginRequired(httpmux, "/station/delete", s,
		stationHandler(stationdb, deleteStationHandler))
	HandleFuncLoginRequired(httpmux, "/station/edit", s,
		stationHandler(stationdb, editStationHandler))

	HandleFuncLoginOptional(httpmux, "/station", s,
		func(w http.ResponseWriter, r *http.Request, user userView) {
		stationViewHandler(
			stationdb, mux, userdb, contactdb, w, r, user)
		})

	httpmux.HandleFunc("/station/map.kml",
		func(w http.ResponseWriter, r *http.Request) {
			stationKMLHandler(stationdb, w, r)
		})
}
