// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "code.google.com/p/goprotobuf/proto"
import "carpcomm/pb"
import "time"
import "reflect"
import "log"

type UserDB struct {
	table *SDBTable
}

const kUserColumn = "pb.User"
const kUserKeyGoogleKey = "google_key"

func NewUserDB(table *SDBTable) *UserDB {
	return &UserDB{table}
}

func (db *UserDB) lookupByGoogleKey(google_key string) (*pb.User, error) {
	users, err := db.table.lookupByKeyValue(
		kUserColumn,
		kUserKeyGoogleKey,
		google_key,
		reflect.TypeOf(pb.User{}))
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}

	if len(users) != 1 {
		log.Printf("Warning: multiple entries found for user " +
			"with google_key=%s", google_key)
	}

	u := (users[0]).(*pb.User)
	return u, nil
}

func (db *UserDB) UserLogin(verifiedEmail,
	displayName, photoUrl, authority string) (*pb.User, error) {

	google_key := verifiedEmail
	u, err := db.lookupByGoogleKey(google_key)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	if u == nil {
		// New user
		u = &pb.User{}
		u.GoogleKey = proto.String(google_key)
		u.Created = &now
		
		id, err := CryptoRandId()
		if err != nil {
			return nil, err
		}
		u.Id = proto.String(id)
	}

	u.Email = proto.String(verifiedEmail)
	u.DisplayName = proto.String(displayName)
	u.PhotoUrl = proto.String(photoUrl)
	u.Authority = proto.String(authority)
	u.LastLogin = &now

	// Update record
	return u, db.Store(u)
}

func (db *UserDB) Store(u *pb.User) error {
	values, err := encodeItem(kUserColumn, u)
	if err != nil {
		return err
	}

	if u.GoogleKey != nil {
		values[kUserKeyGoogleKey] = *u.GoogleKey
	}

	return db.table.put(*u.Id, values)
}

// Returns nil, nil if id was not found.
func (db *UserDB) Lookup(id string) (*pb.User, error) {
	u := &pb.User{}
	found, err := db.table.getProto(id, kUserColumn, u)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return u, nil
}

// Create the database table.
func (db *UserDB) Create() error {
	return db.table.create()
}
