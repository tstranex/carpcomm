// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import (
	"carpcomm/db"
	"carpcomm/mux"
	"carpcomm/pb"
	"carpcomm/scheduler"
	"carpcomm/streamer/contacts"
	"html/template"
	"log"
	"net/http"
	"net/rpc"
	"net/url"
	"sort"
)


var consoleTemplate = NewDebuggableTemplate(
	template.FuncMap{
	        "StationContactsUrl": stationContactsUrl,
        },
	"station_console.html",
	"src/carpcomm/fe/templates/station_console.html",
	"src/carpcomm/fe/templates/page.html")

type satOption struct {
	Id, Name string
	FrequencyHz int64
}

type satOptionList []satOption

func (sol satOptionList) Len() int {
	return len(sol)
}
func (sol satOptionList) Less(i, j int) bool {
	a, b := sol[i], sol[j]
	return a.Name < b.Name
}
func (sol satOptionList) Swap(i, j int) {
	sol[i], sol[j] = sol[j], sol[i]
}

// TODO: We really should cache the result.
func fillSatOptions() (r satOptionList) {
	for _, sat := range db.GlobalSatelliteDB().List {
		var so satOption
		so.Id = *sat.Id
		so.Name = RenderSatelliteShortName(sat.Name)

		for _, c := range sat.Channels {
			if c.Downlink == nil || *c.Downlink == false {
				continue
			}
			so.FrequencyHz = (int64)(*c.FrequencyHz)
			break
		}

		r = append(r, so)
	}

	sort.Sort(r)

	return r
}


type consoleViewContext struct {
	S pb.Station
	Satellites []satOption
}


func consolePageHandler(sdb *db.StationDB, contactdb *db.ContactDB,
	mux *rpc.Client,
	w http.ResponseWriter, r *http.Request, user userView) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "'id' param missing", http.StatusBadRequest)
		return
	}

	station, err := sdb.Lookup(id)
	if err != nil {
		log.Printf("Error looking up station: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	auth, _ := canOperateStation(sdb, id, user.Id)
	if !auth {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var cv consoleViewContext
	cv.S = *station
	cv.Satellites = fillSatOptions()
	c := NewRenderContext(user, cv)
	err = consoleTemplate.Get().ExecuteTemplate(
		w, "station_console.html", c)
	if err != nil {
		log.Printf("Error rendering station console: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func callStation(m *rpc.Client, w http.ResponseWriter,
	station_id string, action string, params url.Values) {
	var args mux.StationCallArgs
	args.StationId = station_id
	u := url.URL{}
	u.Path = "/" + action
	if params != nil {
		u.RawQuery = params.Encode()
	}
	args.URL = u.String()

	var result mux.StationCallResult
	err := m.Call("Coordinator.StationCall", args, &result)
	if err != nil {
		log.Printf("StationCall error: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(result.StatusCode)
	w.Write(result.Data)
}

func canOperateStation(sdb *db.StationDB, station_id, userid string) (
	bool, error) {
	owner, err := sdb.GetStationUserId(station_id)
	if err != nil {
		log.Printf("Error looking up station: %s", err.Error())
		return false, err
	}
	return owner == userid, nil
}

func consoleJSONHandler(sdb *db.StationDB, contactdb *db.ContactDB,
	m *rpc.Client,
	w http.ResponseWriter, r *http.Request, user userView) {
	query := r.URL.Query()

	id := query.Get("id")
	if id == "" {
		http.Error(w, "'id' param missing", http.StatusBadRequest)
		return
	}
	action := query.Get("action")
	if action == "" {
		http.Error(w, "'action' param missing", http.StatusBadRequest)
		return
	}

	auth, _ := canOperateStation(sdb, id, user.Id)
	if !auth {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	// FIXME: authentication and scheduling/locking

	satellite_id := query.Get("satellite_id")
	if satellite_id != "" && 
		db.GlobalSatelliteDB().Map[satellite_id] == nil {
		http.Error(
			w, "invalid satellite_id",
			http.StatusBadRequest)
		return
	}

	switch action {
	case "ReceiverGetState":
		callStation(m, w, id, action, nil)
	case "ReceiverStart":
		log.Printf("ReceiverStart: satellite_id=%s", satellite_id)
		contact_id, err := contacts.StartNewConsoleContact(
			sdb, contactdb, id, user.Id, satellite_id)
		if err != nil {
			log.Printf("StartNewIQContact error: %s", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		v := url.Values{}
		v.Add("stream_url", scheduler.GetStreamURL(contact_id))
		callStation(m, w, id, action, v)
	case "ReceiverStop":
		callStation(m, w, id, action, nil)
	case "ReceiverWaterfallPNG":
		callStation(m, w, id, action, nil)
	case "ReceiverSetFrequency":
		hz := query.Get("hz")
		if hz == "" {
			http.Error(
				w, "'hz' param missing", http.StatusBadRequest)
			return
		}
		v := url.Values{}
		v.Add("hz", hz)
		callStation(m, w, id, action, v)
	case "TNCStart":
		log.Printf("TNCStart: satellite_id=%s", satellite_id)

		host, port := scheduler.GetAPIServer()
		if satellite_id == "" {
			host = ""
			port = "0"
		}

		v := url.Values{}
		v.Add("api_host", host)
		v.Add("api_port", port)
		v.Add("satellite_id", satellite_id)
		callStation(m, w, id, action, v)
	case "TNCStop":
		callStation(m, w, id, action, nil)
	case "TNCGetLatestFrames":
		callStation(m, w, id, action, nil)
	case "MotorGetState":
		callStation(m, w, id, action, nil)
	case "MotorStart":
		program := query.Get("program")
		if program == "" {
			http.Error(w, "'program' param missing",
				http.StatusBadRequest)
			return
		}
		v := url.Values{}
		v.Add("program", program)
		callStation(m, w, id, action, v)
	case "MotorStop":
		callStation(m, w, id, action, nil)
	default:
		http.NotFound(w, r)
	}
}

type ConsoleHandlerFunc func(
	*db.StationDB, *db.ContactDB, *rpc.Client,
	http.ResponseWriter, *http.Request, userView)

func consoleHandler(stationdb *db.StationDB, contactdb *db.ContactDB,
	mux *rpc.Client, f ConsoleHandlerFunc) HandlerFuncWithUser {
	return func(w http.ResponseWriter, r *http.Request, user userView) {
		f(stationdb, contactdb, mux, w, r, user)
	}
}

func AddConsoleHttpHandlers(httpmux *http.ServeMux,
	s *Sessions, stationdb *db.StationDB, contactdb *db.ContactDB,
	mux *rpc.Client) {
	HandleFuncLoginRequired(httpmux, "/station/console", s,
		consoleHandler(stationdb, contactdb, mux, consolePageHandler))
	HandleFuncLoginRequired(httpmux, "/station/call", s,
		consoleHandler(stationdb, contactdb, mux, consoleJSONHandler))
}
