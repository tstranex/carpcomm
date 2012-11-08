// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2012 Carpcomm GmbH

// API for the Carpcomm network.
// BETA
//
// Examples:
//
// Initialize the client:
//
//   c, err := api.NewAPIClient("your_station_id", "your_station_secret")
//
// Fetch some packets:
//
//   packets, err := c.GetLatestPackets("your_satellite_id", 3)
//
// Post a packet:
//
//   var p api.Packet
//   p.Timestamp = time.Now().Unix()
//   p.Frame = make([]byte, 100)
//   err := c.PostPacket("your_satellite_id", p)
package api

import "crypto/tls"
import "net/http"
import "net/url"
import "io/ioutil"
import "encoding/base64"
import "encoding/json"
import "fmt"
import "errors"
import "bytes"

const defaultApiHost = "api.carpcomm.com:5051"
const jsonMimeType = "application/json"

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
	Timestamp int64  // Unix timestamp: UTC seconds since the epoch.
	Frame []byte
}


type postPacketRequest struct {
	StationId string `json:"station_id"`
	StationSecret string `json:"station_secret"`
	Timestamp int64 `json:"timestamp"`
	SatelliteId string `json:"satellite_id"`
	Format string `json:"format"`
	FrameBase64 string `json:"frame_base64"`
}


// Post a packet that was received from the given satellite.
func (c *APIClient) PostPacket(satellite_id string, p Packet) error {
	var req postPacketRequest
	req.StationId = c.station_id
	req.StationSecret = c.station_secret
	req.Timestamp = p.Timestamp
	req.SatelliteId = satellite_id
	req.Format = "FRAME"
	req.FrameBase64 = base64.StdEncoding.EncodeToString(p.Frame)
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	var u url.URL
	u.Scheme = "https"
	u.Host = c.host
	u.Path = "/PostPacket"

	fmt.Printf("url: %s\n", u.String())
	fmt.Printf("body: %s\n", body)

	resp, err := c.client.Post(
		u.String(), jsonMimeType, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		text, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf(
			"HTTP error: %s: %s", resp.Status, (string)(text)))
	}

	return nil
}


type packet struct {
	Timestamp int64 `json:"timestamp"`
	FrameBase64 string `json:"frame_base64"`
}

// Fetch packets for the given satellite from the server.
// At most 'limit' packets will be returned.
// Note that you can only receive packets from satellites that you are
// authorized to read.
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