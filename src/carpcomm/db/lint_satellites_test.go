// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "code.google.com/p/goprotobuf/proto"
import "testing"
import "io/ioutil"
import "carpcomm/pb"
import "strings"
import "regexp"
import "net/url"

const asciiSatellites = "satellites.txt"

const satelliteIdPattern = "^[a-z0-9]+$"
const datumKeyPattern = "^d:[a-z0-9]+:[_a-z0-9]+$"

func LintName(t *testing.T, context string, n pb.TextWithLang) {
	if n.Lang == nil || *n.Lang == "" {
		t.Errorf("%s: missing lang", context)
	}
	if n.Text == nil || *n.Text == "" {
		t.Errorf("%s: missing text", context)
	}
}

func LintChannel(t *testing.T, context string, c pb.Channel) {
	if c.FrequencyHz == nil || *c.FrequencyHz <= 0.0 {
		t.Errorf("%s: missing frequency_hz", context)
	}
}

func LintSchema(t *testing.T, id string, schema pb.TelemetrySchema) {
	if len(schema.Datum) == 0 {
		t.Errorf("%s: Empty schema.", id)
	}

	// Schema keys should begin with "d:<id>:".
	prefix := "d:" + id + ":"
	keys := make(map[string]bool)
	for _, d := range schema.Datum {
		if !strings.HasPrefix(*d.Key, prefix) {
			t.Errorf("%s: Invalid prefix: %s", id, *d.Key)
		}
		key_valid, err := regexp.MatchString(datumKeyPattern, *d.Key)
		if err != nil {
			t.Error(err)
			return
		}
		if !key_valid {
			t.Errorf("%s: Key is invalid: %s", id, *d.Key)
		}
		if keys[*d.Key] {
			t.Errorf("%s: Duplicate key: %s", id, *d.Key)
		}
		keys[*d.Key] = true

		if len(d.Name) == 0 {
			t.Errorf("%s: No datum names.", *d.Key)
		}
		for _, n := range d.Name {
			LintName(t, *d.Key, *n)
		}
	}
}

func LintURL(t *testing.T, id string, rawurl *string) {
	if rawurl == nil {
		t.Errorf("%s: missing url", id)
		return
	}

	_, err := url.Parse(*rawurl)
	if err != nil {
		t.Errorf("%s: invalid url", id)
		return
	}
}

func LintSatellite(t *testing.T, sat pb.Satellite) {
	if sat.Id == nil {
		t.Error("Missing satellite id.")
		return
	}
	id := *sat.Id
	id_ok, err := regexp.MatchString(satelliteIdPattern, id)
	if err != nil {
		t.Error(err)
		return
	}
	if !id_ok {
		t.Errorf("%s: id is invalid", id)
	}

	// It must have at least one name.
	if len(sat.Name) == 0 {
		t.Errorf("%s: should have at least one name", id)
	}
	// Language tags must be filled in.
	for _, n := range sat.Name {
		LintName(t, id + "/name", *n)
	}

	// TODO: Missing country code should be a warning (e.g. ISS).

	LintURL(t, id + "/website", sat.Website)

	for _, c := range sat.Channels {
		LintChannel(t, id + "/channels", *c)
	}

	if sat.Schema != nil {
		LintSchema(t, id, *sat.Schema)
	}
}

func TestLintSatelliteList(t *testing.T) {
	buf, err := ioutil.ReadFile(asciiSatellites)
	if err != nil {
		t.Errorf("Error reading satellites: %s", err.Error())
		return
	}
	sl := &pb.SatelliteList{}
	err = proto.UnmarshalText((string)(buf), sl)
	if err != nil {
		t.Errorf("Error unmarshalling satellites: %s", err.Error())
		return
	}

	if len(sl.Satellite) == 0 {
		t.Error("No satellites found.")
		return
	}

	for _, sat := range sl.Satellite {
		LintSatellite(t, *sat)
	}
}
