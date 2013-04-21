// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "net/http"
import "log"
import "flag"
import "carpcomm/db"
import "carpcomm/pb"
import "code.google.com/p/goprotobuf/proto"
import "io/ioutil"

var ranking_list = flag.String(
	"ranking_list",
	"data/rankings.RankingList",
	"Path to RankingList proto file")


func loadRankingList() *pb.RankingList {
	buf, err := ioutil.ReadFile(*ranking_list)
	if err != nil {
		log.Panicf("Error loading rankings: %s", err.Error())
		return nil
	}
	rankings := &pb.RankingList{}
	err = proto.Unmarshal(buf, rankings)
	if err != nil {
		log.Panicf("Error loading rankings: %s", err.Error())
		return nil
	}
	log.Printf("Loaded rankings from %s", *ranking_list)
	return rankings
}

var globalRankings *pb.RankingList = nil

func GlobalRankings() *pb.RankingList {
	if globalRankings != nil {
		return globalRankings
	}
	globalRankings = loadRankingList()
	return globalRankings
}


type RankingView struct {
	Rank int
	User userView
	IQCount, MorseCount, FrameCount int
	Score int
}

func fillRankingView(userdb *db.UserDB, rank int, r *pb.Ranking) (
	rv RankingView) {
	rv.Rank = rank + 1
	rv.Score = (int)(*r.Score)
	for _, c := range r.Counts {
		switch *c.Format {
		case pb.Contact_Blob_IQ: rv.IQCount = (int)(*c.Count)
		case pb.Contact_Blob_MORSE: rv.MorseCount = (int)(*c.Count)
		case pb.Contact_Blob_FRAME: rv.FrameCount = (int)(*c.Count)
		}
	}
	rv.User = LookupUserView(userdb, *r.UserId)
	return rv
}

var rankingTemplate = NewDebuggableTemplate(
	nil,
	"ranking.html",
	"src/carpcomm/fe/templates/ranking.html",
	"src/carpcomm/fe/templates/page.html")

func rankingHandler(
	cdb *db.ContactDB, userdb *db.UserDB,
	w http.ResponseWriter, r *http.Request, user userView) {

	items := make([]RankingView, len(GlobalRankings().Ranking))
	for i, r := range GlobalRankings().Ranking {
		items[i] = fillRankingView(userdb, i, r)
	}

	c := NewRenderContext(user, items)
	err := rankingTemplate.Get().ExecuteTemplate(w, "ranking.html", c)
	if err != nil {
		log.Printf("Error rendering ranking: %s", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}


func AddRankingHttpHandlers(httpmux *http.ServeMux, s *Sessions,
	cdb *db.ContactDB, userdb *db.UserDB) {
	HandleFuncLoginOptional(httpmux, "/ranking", s,
		func(w http.ResponseWriter, r *http.Request, user userView) {
		rankingHandler(cdb, userdb, w, r, user)
	})
}
