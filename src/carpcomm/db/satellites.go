// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "carpcomm/pb"
import "code.google.com/p/goprotobuf/proto"
import "log"
import "flag"
import "io/ioutil"
import "sort"

var satellite_list = flag.String(
	"satellite_list",
	"data/satellites.SatelliteList",
	"Path to SatelliteList proto file")

type satList []*pb.Satellite

func (sl satList) Len() int {
	return len(sl)
}
func (sl satList) Less(i, j int) bool {
	a, b := sl[i], sl[j]

	// 1. has photo
	if len(a.Photo) > 0 && len(b.Photo) == 0 {
		return true
	} else if len(a.Photo) == 0 && len(b.Photo) > 0 {
		return false
	}

	// 2. has schema
	if a.Schema != nil && b.Schema == nil {
		return true
	} else if a.Schema == nil && b.Schema != nil {
		return false
	}

	// 3. launch timestamp
	if a.LaunchTimestamp != nil && b.LaunchTimestamp != nil &&
	*a.LaunchTimestamp != *b.LaunchTimestamp {
		return *a.LaunchTimestamp > *b.LaunchTimestamp
	} else if a.LaunchTimestamp == nil && b.LaunchTimestamp != nil {
		return false
	} else if a.LaunchTimestamp != nil && b.LaunchTimestamp == nil {
		return true
	}

	// 4. tie breaker
	return *a.Id < *b.Id
}
func (sl satList) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

type SatelliteDB struct {
	List []*pb.Satellite
	Map map[string]*pb.Satellite
}

func LoadSatelliteDB(filename string) (db *SatelliteDB, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	sl := &pb.SatelliteList{}
	err = proto.Unmarshal(buf, sl)
	if err != nil {
		return nil, err
	}

	sort.Sort((satList)(sl.Satellite))

	db = &SatelliteDB{}
	db.List = sl.Satellite
	db.Map = make(map[string]*pb.Satellite)
	for _, s := range db.List {
		db.Map[*s.Id] = s
	}
	return db, nil
}

func loadGlobalDB() *SatelliteDB {
	db, err := LoadSatelliteDB(*satellite_list)
	if err != nil || db == nil {
		log.Panicf("Error loading satellites from %s: %s",
			*satellite_list, err.Error())
	}
	log.Printf("Loaded satellites from %s", *satellite_list)
	return db
}

var globalDB *SatelliteDB = nil


func GlobalSatelliteDB() *SatelliteDB {
	if globalDB != nil {
		return globalDB
	}
	globalDB = loadGlobalDB()
	return globalDB
}
