// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "carpcomm/pb"
import "code.google.com/p/goprotobuf/proto"

import crand "crypto/rand"
import "io"
import "reflect"
import "encoding/base64"
import "time"
import "log"

func NewStation(userid string) (*pb.Station, error) {
	var s pb.Station

	s.Userid = proto.String(userid)

	id, err := CryptoRandId()
	if err != nil {
		return nil, err
	}
	s.Id = proto.String(id)

	secret := make([]byte, 20)
	_, err = io.ReadFull(crand.Reader, secret)
	if err != nil {
		return nil, err
	}
	s.Secret = proto.String(base64.StdEncoding.EncodeToString(secret))

	now := time.Now().Unix()
	s.Created = &now

	s.Name = proto.String("Unnamed Station")

	return &s, nil
}

/*
func DefaultCapabilities() *pb.Capabilities {
	return &pb.Capabilities{}
}
*/


type StationDB struct {
	table *SDBTable
	useridCache map[string]string
}

const kStationColumn = "pb.Station"

func NewStationDB(table *SDBTable) *StationDB {
	return &StationDB{
		table,
		make(map[string]string)}
}

func (db *StationDB) Store(s *pb.Station) error {
	return db.table.setProto(*s.Id, kStationColumn, s)
}

// Returns nil, nil if id was not found.
func (db *StationDB) Lookup(id string) (*pb.Station, error) {
	s := &pb.Station{}
	found, err := db.table.getProto(id, kStationColumn, s)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return s, nil
}

func (db *StationDB) GetStationUserId(id string) (string, error) {
	// Currently, the userid of a station shouldn't change so we can cache
	// it. This may change in the future.
	userid, found := db.useridCache[id]
	if found {
		return userid, nil
	}

	s, err := db.Lookup(id)
	if err != nil {
		return "", err
	}
	db.useridCache[id] = *s.Userid
	return *s.Userid, nil
}

// FIXME: Consider caching to speed this up.
func (db *StationDB) GetStationName(id string) (string, error) {
	s, err := db.Lookup(id)
	if err != nil {
		log.Printf("Error looking up station: %s", err.Error())
		return "", err
	}
	if s == nil || s.Name == nil {
		return "", nil
	}
	return *s.Name, nil
}

func (db *StationDB) Delete(id string) error {
	return db.table.delete(id, kStationColumn)
}

func (db *StationDB) AllStations() ([]*pb.Station, error) {
	result, err := db.table.getAll(kStationColumn,
		reflect.TypeOf(pb.Station{}))
	if err != nil {
		return nil, err
	}
	conv := make([]*pb.Station, len(result))
	for i, v := range(result) {
		conv[i] = v.(*pb.Station)
	}
	return conv, nil
}

func (db *StationDB) NumStations() (int, error) {
	stations, err := db.AllStations()
	if err != nil {
		return 0, err
	}
	return len(stations), nil
}

func (db *StationDB) UserStations(userid string) ([]*pb.Station, error) {
	if userid == "" {
		return nil, nil
	}

	// FIXME: use a search query

	result, err := db.table.getAll(kStationColumn,
		reflect.TypeOf(pb.Station{}))
	if err != nil {
		return nil, err
	}

	filtered := []*pb.Station{}
	for _, v := range(result) {
		s := v.(*pb.Station)
		if s.Userid != nil && *s.Userid == userid {
			filtered = append(filtered, s)
		}
	}
	return filtered, nil
}

func (db *StationDB) Create() error {
	return db.table.create()
}