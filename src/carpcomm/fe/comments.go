// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "net/http"
import "log"
import "carpcomm/db"
import "carpcomm/pb"
import "time"

type CommentView struct {
	Text string
	Timestamp string
	User userView
}

type CommentsView struct {
	ObjectId string
	C []CommentView
}

const commentTimeFormat = "2006-01-02 15:04 MST"

func FillCommentView(c pb.Comment, userdb *db.UserDB, cv *CommentView) {
	cv.Text = *c.Text
	cv.Timestamp = time.Unix(*c.Timestamp, 0).UTC().Format(
		commentTimeFormat)
	if c.UserId != nil {
		cv.User = LookupUserView(userdb, *c.UserId)
	}
}

func LoadCommentsByObjectId(object_id string,
	commentdb *db.CommentDB, userdb *db.UserDB) (CommentsView, error) {
	var cv CommentsView
	cv.ObjectId = object_id
	comments, err := commentdb.SearchByObjectId(object_id)
	if err != nil {
		log.Printf("Error looking up comments for object id %s: %s",
			object_id, err.Error())
		return cv, err
	}
	cv.C = make([]CommentView, len(comments))
	for i, c := range comments {
		FillCommentView(*c, userdb, &cv.C[i])
	}
	return cv, nil
}

func validateObjectId(object_id string) (ok bool, redirect_url string) {
	if ok, redirect_url = validateSatelliteObjectId(object_id); ok {
		return true, redirect_url
	}
	return false, ""
}

func commentPostHandler(
	cdb *db.CommentDB,
	w http.ResponseWriter, r *http.Request,
	user userView) {

	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("commentPostHandler form: %v\n\n", r.Form)

	object_id := r.Form.Get("object_id")
	ok, redirect_url := validateObjectId(object_id)
	if !ok {
		http.Error(w, "Invalid object_id", http.StatusBadRequest)
		return
	}

	text := r.Form.Get("text")
	if len(text) == 0 {
		http.Error(w, "Comment text missing", http.StatusBadRequest)
		return
	}
	if len(text) > 256 {
		http.Error(w, "Exceeded size limit", http.StatusBadRequest)
		return
	}

	var userid *string
	if user.Id == "" {
		userid = nil
	} else {
		userid = &user.Id
	}
	comment, err := db.NewComment(object_id, text, userid)
	if err != nil {
		log.Printf("Error creating comment: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	log.Printf("Comment: %s", comment)

	err = cdb.Store(comment)
	if err != nil {
		log.Printf("Error storing comment: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirect_url, http.StatusFound)
}

func AddCommentsHttpHandlers(httpmux *http.ServeMux, s *Sessions,
	commentdb *db.CommentDB) {
	HandleFuncLoginRequired(httpmux, "/comments/post", s,
		func(w http.ResponseWriter, r *http.Request, user userView) {
		commentPostHandler(commentdb, w, r, user)
	})
}
