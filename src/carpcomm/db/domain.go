// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "launchpad.net/goamz/aws"

import "log"

// A group of tables in a single region that make up the full database.
type Domain struct {
	db_prefix string
	auth aws.Auth
}

func NewDomain(db_prefix string) (*Domain, error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Printf("AWS auth error: %s", err.Error())
		return nil, err
	}

	var d Domain
	d.db_prefix = db_prefix
	d.auth = auth
	return &d, nil
}

func (d *Domain) NewUserDB() *UserDB {
	return NewUserDB(NewSDBTable(&d.auth, &aws.USEast, d.db_prefix+"users"))
}

func (d *Domain) NewStationDB() *StationDB {
	return NewStationDB(
		NewSDBTable(&d.auth, &aws.USEast, d.db_prefix+"stations"))
}

func (d *Domain) NewContactDB() *ContactDB {
	return NewContactDB(
		NewSDBTable(&d.auth, &aws.USEast, d.db_prefix+"contacts"))
}

func (d *Domain) NewCommentDB() *CommentDB {
	return NewCommentDB(
		NewSDBTable(&d.auth, &aws.USEast, d.db_prefix+"comments"))
}