// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2012 Carpcomm GmbH

// API for the Carpcomm network.

package api

import "crypto/tls"
import "net/http"
import "net/url"
import "io/ioutil"
import "encoding/base64"
import "encoding/json"
import "fmt"
import "errors"

const defaultApiHost = "api.carpcomm.com:5051"

type APIClient struct {
	host string
	client http.Client
	station_id, station_secret string
}

func NewAPIClient(station_id, station_secret string) (*APIClient, error) {
	if station_id == "" {
		return nil, errors.New("station_id is empty")
	}
	if station_secret == "" {
		return nil, errors.New("station_secret is empty")
	}

	// FIXME(tstranex): Change InsecureSkipVerify to false once the server
	// certificate is registered.
	tr := &http.Transport{
	TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	var c APIClient
	c.host = defaultApiHost
	c.client = http.Client{Transport: tr}
	c.station_id = station_id
	c.station_secret = station_secret

	return &c, nil
}


type Packet struct {
	Timestamp int64
	Frame []byte
}

type packet struct {
	Timestamp int64 `json:"timestamp"`
	FrameBase64 string `json:"frame_base64"`
}

// Fetch packets for the given satellite from the server.
// At most 'limit' packets will be returned.
func (c *APIClient) GetLatestPackets(satellite_id string, limit int) (
	[]Packet, error) {

	v := make(url.Values)
	v.Set("station_id", c.station_id)
	v.Set("station_secret", c.station_secret)
	v.Set("satellite_id", satellite_id)
	v.Set("limit", fmt.Sprintf("%d", limit))

	var u url.URL
	u.Scheme = "https"
	u.Host = c.host
	u.Path = "/GetLatestPackets"
	u.RawQuery = v.Encode()

	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("HTTP error: %s", resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var d []packet
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}

	result := make([]Packet, len(d))
	for i, p := range d {
		f, err := base64.StdEncoding.DecodeString(p.FrameBase64)
		if err != nil {
			return nil, err
		}

		result[i].Timestamp = p.Timestamp
		result[i].Frame = f
	}

	return result, nil
}