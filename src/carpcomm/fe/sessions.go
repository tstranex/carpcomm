// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "code.google.com/p/goauth2/oauth"
import "net/http"
import "net/url"
import "log"
import "io/ioutil"
import "crypto/rand"
import "math/big"
import "encoding/json"
import "sync"
import "html/template"

import "carpcomm/db"

//import "carpcomm/pb"

const sidCookie = "sid"


type userView struct {
	Id, Name, PhotoUrl string
}

type Sessions struct {
	sidToUser map[string]userView
	lock        sync.Mutex
}

func NewSessions() (s *Sessions) {
	s = new(Sessions)
	s.sidToUser = make(map[string]userView)
	return s
}

func (s *Sessions) Begin(userid string, displayName, photoUrl *string) (
	*http.Cookie, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return nil, err
	}
	sid := n.String()

	var uv userView
	uv.Id = userid
	if displayName != nil {
		uv.Name = *displayName
	} else {
		uv.Name = ""
	}
	if photoUrl != nil {
		uv.PhotoUrl = *photoUrl
	} else {
		uv.PhotoUrl = "/images/unknown_avatar.png"
	}

	s.lock.Lock()
	s.sidToUser[sid] = uv
	s.lock.Unlock()

	cookie := http.Cookie{
		Name:  sidCookie,
		Value: sid,
	}
	return &cookie, nil
}

func (s *Sessions) GetUser(r *http.Request) userView {
	cookie, err := r.Cookie(sidCookie)
	if err != nil {
		return userView{"", "", ""}
	}
	sid := cookie.Value

	s.lock.Lock()
	uv, ok := s.sidToUser[sid]
	s.lock.Unlock()

	if !ok {
		return userView{"", "", ""}
	}
	return uv
}

func (s *Sessions) End(r *http.Request) *http.Cookie {
	new_cookie := &http.Cookie{
		Name:   sidCookie,
		Value:  "",
		MaxAge: -1,
	}

	cookie, err := r.Cookie(sidCookie)
	if err != nil {
		return new_cookie
	}
	sid := cookie.Value

	s.lock.Lock()
	delete(s.sidToUser, sid)
	s.lock.Unlock()

	return new_cookie
}

const GoogleOAuth2 = "https://accounts.google.com/o/oauth2/auth"
const googleAuthScope = "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"
const googleTokenURL = "https://accounts.google.com/o/oauth2/token"

var debugAuthConfig = oauth.Config{
	ClientId:     "FIXME_PUT_GOOGLE_CLIENT_ID_HERE",
	ClientSecret: "FIXME_PUT_GOOGLE_CLIENT_SECRET_HERE",
	Scope:       googleAuthScope,
	AuthURL:     GoogleOAuth2,
	TokenURL:    googleTokenURL,
	RedirectURL: "http://localhost:8000/authcb_google_oauth2",
}

var prodAuthConfig = oauth.Config{
        ClientId:     "FIXME_PUT_GOOGLE_CLIENT_ID_HERE",
        ClientSecret: "FIXME_PUT_GOOGLE_CLIENT_SECRET_HERE",
	Scope:        googleAuthScope,
        AuthURL:      GoogleOAuth2,
        TokenURL:     googleTokenURL,
        RedirectURL:  "http://carpcomm.com/authcb_google_oauth2",
}

func getConfig() *oauth.Config {
	if *debug_auth {
		log.Print("Using debug auth config")
		return &debugAuthConfig
	}
	return &prodAuthConfig
}

var loginTemplate = template.Must(template.ParseFiles(
	"src/carpcomm/fe/templates/login.html",
	"src/carpcomm/fe/templates/page.html"))

func LoginHandler(s *Sessions, userdb *db.UserDB,
	w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL.String())
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/"
	}

	uv := s.GetUser(r)
	if uv.Id != "" {
		http.Redirect(w, r, redirect, http.StatusFound)
		return
	}

	login_url := getConfig().AuthCodeURL(redirect)
	c := NewRenderContext(uv, login_url)
	err := loginTemplate.ExecuteTemplate(w, "login.html", c)
	if err != nil {
		log.Printf("Error rendering login page: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	//http.Redirect(w, r, login_url, http.StatusFound)
}

type oAuthUserInfo struct {
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
	Timezone       string
	Gender         string
}

func GoogleLoginCallbackHandler(s *Sessions, userdb *db.UserDB,
	w http.ResponseWriter, r *http.Request) {

	log.Print("hello GoogleLoginCallbackHandler")

	code := r.URL.Query().Get("code")
	redirect := r.URL.Query().Get("state")
	log.Printf("redirect: %s", redirect)

	t := &oauth.Transport{Config: getConfig()}
	token, err := t.Exchange(code)
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("GoogleLoginCallbackHandler got token: %s\n", token)

	resp, err := t.Client().Get(
		"https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		log.Print(err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Print(err)
		return
	}

	var userinfo oAuthUserInfo
	if err := json.Unmarshal(data, &userinfo); err != nil {
		log.Print(err)
		return
	}
	log.Printf("User logged in: %s", userinfo.Email)

	if userinfo.Verified_email == false || userinfo.Email == "" {
		log.Printf("email not verified")
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	user, err := userdb.UserLogin(
		userinfo.Email, userinfo.Name, userinfo.Picture, GoogleOAuth2)
	if err != nil {
		log.Printf("UserDB error: %s", err.Error())
		http.Error(w, "Unable to complete login",
			http.StatusInternalServerError)
		return
	}

	cookie, err := s.Begin(
		*user.Id, user.DisplayName, user.PhotoUrl)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("session id: %v", cookie)

	http.SetCookie(w, cookie)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func LogoutHandler(s *Sessions, userdb *db.UserDB,
	w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/"
	}

	cookie := s.End(r)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func makeHandler(s *Sessions, userdb *db.UserDB,
	f func(*Sessions, *db.UserDB,
		http.ResponseWriter, *http.Request),) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(s, userdb, w, r)
	}
}

func AddLoginHttpHandlers(s *Sessions, userdb *db.UserDB) {
	http.HandleFunc("/login", makeHandler(s, userdb, LoginHandler))
	http.HandleFunc("/authcb_google_oauth2",
		makeHandler(s, userdb, GoogleLoginCallbackHandler))
	http.HandleFunc("/logout", makeHandler(s, userdb, LogoutHandler))
}

type HandlerFuncWithUser func(
	w http.ResponseWriter, r *http.Request, user userView)

func HandleFuncLoginRequired(
	mux *http.ServeMux, path string, s *Sessions, f HandlerFuncWithUser) {
	wrapped := func(w http.ResponseWriter, r *http.Request) {
		uv := s.GetUser(r)
		log.Printf("current userid: %s", uv.Id)
		if uv.Id != "" {
		//if true {
			f(w, r, uv)
		} else {
			var url url.URL
			url.Path = "/login"
			q := url.Query()
			q.Set("redirect", r.RequestURI)
			url.RawQuery = q.Encode()
			log.Print(url.String())
			http.Redirect(w, r, url.String(), http.StatusFound)
		}
	}
	mux.HandleFunc(path, wrapped)
}

func HandleFuncLoginOptional(
	mux *http.ServeMux, path string, s *Sessions, f HandlerFuncWithUser) {
	wrapped := func(w http.ResponseWriter, r *http.Request) {
		f(w, r, s.GetUser(r))
	}
	mux.HandleFunc(path, wrapped)
}
