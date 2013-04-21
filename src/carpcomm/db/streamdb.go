// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "fmt"
import "carpcomm/pb"
import "reflect"

type ContactDB struct {
	table *SDBTable
}

const kContactColumn = "pb.Contact"
const kContactKeyStationId = "station_id"
const kContactKeyUserId = "user_id"
const kContactKeySatelliteId = "satellite_id"
const kContactKeyTimestamp = "timestamp"

func NewContactDB(table *SDBTable) *ContactDB {
	return &ContactDB{table}
}

func emptyIfUnknown(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (db *ContactDB) Store(contact *pb.Contact) error {
	values, err := encodeItem(kContactColumn, contact)
	if err != nil {
		return err
	}

	// TODO: We should not populate unknown attributes at all.
	values[kContactKeyStationId] = emptyIfUnknown(contact.StationId)
	values[kContactKeyUserId] = emptyIfUnknown(contact.UserId)
	values[kContactKeySatelliteId] = emptyIfUnknown(contact.SatelliteId)
	values[kContactKeyTimestamp] = fmt.Sprintf(
		"%016x", *contact.StartTimestamp)

	return db.table.put(*contact.Id, values)
}

func (db *ContactDB) Lookup(id string) (*pb.Contact, error) {
	s := &pb.Contact{}
	found, err := db.table.getProto(id, kContactColumn, s)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return s, nil
}

// Results are sorted by timestamp (newest first).
func (db *ContactDB) searchByKey(column, key string, limit int) ([]*pb.Contact, error) {
	query := fmt.Sprintf(
		"select * from `%s` where `%s` = '%s' and `%s` is not null order by `%s` desc limit %d",
		db.table.domain.Name,
		column,
		key,
		kContactKeyTimestamp,
		kContactKeyTimestamp,
		limit)
	result, err := db.table.search(
		query, kContactColumn, reflect.TypeOf(pb.Contact{}))
	if err != nil {
		return nil, err
	}
	conv := make([]*pb.Contact, len(result))
	for i, v := range(result) {
		conv[i] = v.(*pb.Contact)
	}
	return conv, nil
}

// Results are sorted by timestamp (newest first).
func (db *ContactDB) SearchBySatelliteId(satellite_id string, limit int) (
	[]*pb.Contact, error) {
	return db.searchByKey(kContactKeySatelliteId, satellite_id, limit)
}

// Results are sorted by timestamp (newest first).
func (db *ContactDB) SearchByStationId(station_id string, limit int) (
	[]*pb.Contact, error) {
	return db.searchByKey(kContactKeyStationId, station_id, limit)
}

// Results are sorted by timestamp (newest first).
func (db *ContactDB) SearchByUserId(user_id string, limit int) (
	[]*pb.Contact, error) {
	return db.searchByKey(kContactKeyUserId, user_id, limit)
}

func (db *ContactDB) GetAll() ([]*pb.Contact, error) {
	result, err := db.table.getAll(kContactColumn,
		reflect.TypeOf(pb.Contact{}))
	if err != nil {
		return nil, err
	}
	conv := make([]*pb.Contact, len(result))
	for i, v := range(result) {
		conv[i] = v.(*pb.Contact)
	}
	return conv, nil
}

func (db *ContactDB) Create() error {
	return db.table.create()
}
