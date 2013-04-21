// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package mux

import "carpcomm/db"

func AuthenticateStation(sdb *db.StationDB, id, secret string) (bool, error) {
	s, err := sdb.Lookup(id)
	if err != nil {
		return false, err
	}
	if s == nil {
		return false, nil
	}
	return *s.Id == id && *s.Secret == secret, nil
}