// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "net/http"
import "net/url"
import "html/template"
import "log"
import "carpcomm/db"
import "carpcomm/pb"
import fe_telemetry "carpcomm/fe/telemetry"
import "carpcomm/streamer/contacts"
import "carpcomm/scheduler"
import "time"
import "strings"
import "strconv"
import "image/png"

const satelliteURLPrefix = "/satellite/"
const satelliteListUrl = "/satellite/list"
const satelliteOrbitUrl = "/satellite/orbit"

func orbitHandler(w http.ResponseWriter, r *http.Request, user userView) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "'id' param missing", http.StatusBadRequest)
		return
	}

	sat := db.GlobalSatelliteDB().Map[id]
	if sat == nil {
		http.NotFound(w, r)
		return
	}

	if sat.Tle == nil || *sat.Tle == "" {
		http.NotFound(w, r)
		return
	}

	points, err := scheduler.PassDetails(
		time.Now(),
		5 * time.Hour,
		0.0, 0.0, 0.0,  // Observer is irrelevant.
		*sat.Tle,
		25.0)
	if err != nil {
		log.Printf("Error getting PassDetails: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	png.Encode(w, orbitMap(points))
}

func satelliteViewURL(id string) string {
	return satelliteURLPrefix + id
}

func satelliteObjectId(id string) string {
	return satelliteViewURL(id)
}

func validateSatelliteObjectId(object_id string) (
	ok bool, redirect_url string) {
	if !strings.HasPrefix(object_id, satelliteURLPrefix) {
		return false, ""
	}
	return true, object_id
}

type contactView struct {
	Timestamp string
	User userView
}

const ContactTimeFormat = "2006-01-02 15:04 MST"

func fillContactView(c pb.Contact, userdb *db.UserDB) *contactView {
	if c.UserId == nil {
		return nil
	}
	var cv contactView
	cv.User = LookupUserView(userdb, *c.UserId)
	cv.Timestamp = time.Unix(*c.StartTimestamp, 0).UTC().Format(
		ContactTimeFormat)
	return &cv
}

type satelliteViewContext struct {
	S *pb.Satellite
	TelemetryHead []fe_telemetry.LabelValue
	TelemetryTail [][]fe_telemetry.LabelValue
	LatestContact *contactView
	Comments CommentsView
	Stations []*pb.Station
}

func GetURLHost(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Printf("Error in GetURLHost for %s: %s",
			rawurl, err.Error())
	}
	return u.Host
}

func RenderSatelliteName(names []*pb.TextWithLang) string {
	if len(names) == 0 {
		log.Printf("No names given to RenderSatelliteName.");
	}
	// TODO: Handle name localization properly. For now, just concatenate
	// all of them.
	result := ""
	for _, name := range names {
		result =  result + *name.Text + " "
	}
	return result[:len(result)-1]
}

func RenderSatelliteShortName(names []*pb.TextWithLang) string {
	if len(names) == 0 {
		log.Printf("No names given to RenderSatelliteName.");
	}
	// TODO: Handle name localization properly. For now, just take the
	// first.
	return *names[0].Text
}

var satelliteViewTemplate = NewDebuggableTemplate(
	template.FuncMap{
	        "EngNotationHz": fe_telemetry.EngNotationHz,
	        "EngNotationWatt": fe_telemetry.EngNotationWatt,
	        "GetURLHost": GetURLHost,
	        "RenderSatelliteName": RenderSatelliteName,
        },
	"satellite.html",
	"src/carpcomm/fe/templates/satellite.html",
	"src/carpcomm/fe/templates/comments.html",
	"src/carpcomm/fe/templates/page.html")

func satelliteViewHandler(
	cdb *db.ContactDB, userdb *db.UserDB, stationdb *db.StationDB,
	commentdb *db.CommentDB,
	w http.ResponseWriter, r *http.Request, user userView) {

	if len(r.URL.Path) < len(satelliteURLPrefix) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id := r.URL.Path[len(satelliteURLPrefix):]

	sat := db.GlobalSatelliteDB().Map[id]
	if sat == nil {
		http.NotFound(w, r)
		return
	}

	// TODO: It would be better if we could restrict to contacts which
	// have telemetry.
	contacts, err := cdb.SearchBySatelliteId(id, 100)
	if err != nil {
		log.Printf("cdb.SearchBySatelliteId error: %s", err.Error())
		// Continue since this isn't a critical error.
	}

	t := make([]pb.TelemetryDatum, 0)
	var latest_contact *contactView
	for _, c := range contacts {
		for _, b := range c.Blob {
			if b.Format != nil &&
				*b.Format == pb.Contact_Blob_DATUM {
				t = append(t, *b.Datum)
				if latest_contact == nil {
					latest_contact = fillContactView(
						*c, userdb)
				}
			}
		}
	}

	sv := satelliteViewContext{}
	sv.S = sat
	if sat.Schema != nil {
		t := fe_telemetry.RenderTelemetry(*sat.Schema, t, "en")
		if len(t) > 0 {
			sv.TelemetryHead = t[0]
		}
		if len(t) > 1 {
			sv.TelemetryTail = t[1:]
		}
	}
	sv.LatestContact = latest_contact

	sv.Comments, _ = LoadCommentsByObjectId(
		satelliteObjectId(id), commentdb, userdb)

	sv.Stations, err = stationdb.UserStations(user.Id)
	if err != nil {
		log.Printf("Error getting user stations: %s", err.Error())
		// Continue rendering since it's not a critial error.
	}

	c := NewRenderContext(user, &sv)
	err = satelliteViewTemplate.Get().ExecuteTemplate(
		w, "satellite.html", c)
	if err != nil {
		log.Printf("Error rendering satellite view: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

type DebuggableTemplate struct {
	funcs template.FuncMap
	name string
	files []string
	tpl *template.Template
}

func NewDebuggableTemplate(
	funcs template.FuncMap,
	name string,
	files ...string) *DebuggableTemplate {
	var dt DebuggableTemplate
	dt.funcs = funcs
	dt.name = name
	dt.files = files
	dt.tpl = dt.load()
	return &dt
}

func (dt *DebuggableTemplate) load() *template.Template {
	tpl := template.New(dt.name)
	tpl.Funcs(dt.funcs)
	return template.Must(tpl.ParseFiles(dt.files...))
}

func (dt *DebuggableTemplate) Get() *template.Template {
	if *debug_templates {
		log.Printf("Reloading template: %s", dt.name)
		return dt.load()
	}
	return dt.tpl
}

var satelliteListTemplate = NewDebuggableTemplate(
	template.FuncMap{
	        "RenderSatelliteName": RenderSatelliteName,
	        "SatelliteViewURL": satelliteViewURL,
        },
	"satellite_list",
	"src/carpcomm/fe/templates/satellite_list.html",
	"src/carpcomm/fe/templates/page.html")

func satelliteListHandler(
	w http.ResponseWriter, r *http.Request, user userView) {
	c := NewRenderContext(user, db.GlobalSatelliteDB().List)
	err := satelliteListTemplate.Get().ExecuteTemplate(
		w, "satellite_list.html", c)
	if err != nil {
		log.Printf("Error rendering satellite list: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}


type contactConfirmContext struct {
	SatelliteUrl, SatelliteName string
	Telemetry [][]fe_telemetry.LabelValue
	Data string
}

var contactConfirmTemplate = NewDebuggableTemplate(
	nil,
	"contact_confirm",
	"src/carpcomm/fe/templates/contact_confirm.html",
	"src/carpcomm/fe/templates/page.html")

func satellitePostContactHandler(
	sdb *db.StationDB, cdb *db.ContactDB,
	w http.ResponseWriter, r *http.Request,
	user userView) {

	if r.Method != "POST" {
		http.Redirect(w, r, satelliteListUrl, http.StatusFound)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("satellitePostContactHandler form: %v\n\n", r.Form)

	satellite_id := r.Form.Get("satellite_id")
	if satellite_id == "" {
		http.Error(w, "Missing satellite id", http.StatusBadRequest)
		return
	}
	sat := db.GlobalSatelliteDB().Map[satellite_id]
	if sat == nil {
		http.Error(w, "Unknown satellite id", http.StatusBadRequest)
		return
	}
	
	timestamp, err := strconv.ParseInt(r.Form.Get("timestamp"), 10, 64)
	if err != nil {
		http.Error(w, "Can't parse timestamp.", http.StatusBadRequest)
		return
	}

	data := r.Form.Get("data")
	frame := ([]byte)(data)

	var station *pb.Station

	// There are two options: logged-in or anonymous.
	if user.Id == "" {
		// Anonymous
		station = nil
	} else {
		// Logged-in user

		station_id := r.Form.Get("station_id")
		if station_id == "" {
			http.Error(
				w, "Missing station id", http.StatusBadRequest)
			return
		}

		station, err = sdb.Lookup(station_id)
		if err != nil {
			log.Printf("Error looking up station: %s", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	contact, poperr := contacts.PopulateContact(
		satellite_id,
		timestamp,
		"FREEFORM",
		frame,
		user.Id,
		"",
		station)

	if poperr != nil {
		poperr.HttpError(w)
		return
	}

	log.Printf("Contact: %s", contact)

	err = cdb.Store(contact)
	if err != nil {
		log.Printf("Error storing contact: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var cc contactConfirmContext
	cc.SatelliteUrl = satelliteViewURL(*sat.Id)
	cc.SatelliteName = RenderSatelliteName(sat.Name)
	cc.Data = data

	if sat.Schema != nil {
		t := make([]pb.TelemetryDatum, 0)
		for _, b := range contact.Blob {
			if b.Format != nil &&
				*b.Format == pb.Contact_Blob_DATUM {
				t = append(t, *b.Datum)
			}
		}

		cc.Telemetry = fe_telemetry.RenderTelemetry(
			*sat.Schema, t, "en")
	}

	err = contactConfirmTemplate.Get().ExecuteTemplate(
		w, "contact_confirm.html", NewRenderContext(user, cc))
	if err != nil {
		log.Printf(
			"Error rendering contact_confirm view: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func AddSatelliteHttpHandlers(httpmux *http.ServeMux, s *Sessions,
	cdb *db.ContactDB, userdb *db.UserDB,
	stationdb *db.StationDB, commentdb *db.CommentDB) {
	HandleFuncLoginOptional(httpmux, satelliteURLPrefix, s,
		func(w http.ResponseWriter, r *http.Request, user userView) {
		satelliteViewHandler(cdb, userdb, stationdb, commentdb,
			w, r, user)
	})
	HandleFuncLoginOptional(httpmux, satelliteOrbitUrl, s,
		orbitHandler)
	HandleFuncLoginOptional(httpmux, satelliteListUrl, s,
		satelliteListHandler)
	HandleFuncLoginOptional(httpmux, "/satellite/contact", s,
		func(w http.ResponseWriter, r *http.Request, user userView) {
		satellitePostContactHandler(stationdb, cdb, w, r, user)
	})
}
