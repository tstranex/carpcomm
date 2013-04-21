// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package mux

import "net"
import "net/http"
import "net/url"
import "crypto/tls"
import "bufio"
import "io"
import "encoding/json"
import "log"
import "errors"
import "time"
import "strings"

const keepAlivePingInterval = 20 * time.Minute
const stationRPCTimeout = 1 * time.Minute

func doStationRPC(conn net.Conn, r *http.Request) (*Response, error) {
	conn.SetDeadline(time.Now().Add(stationRPCTimeout))

	err := r.Write(conn)
	if err != nil {
		log.Printf("Error while writing: %s", err.Error())
		return nil, err
	}

	http_resp, err := http.ReadResponse(bufio.NewReader(conn), r)
	if err != nil {
		log.Printf("Error reading response: %s", err.Error())
		return nil, err
	}

	var resp Response
	resp.code = http_resp.StatusCode
	resp.data = []byte{}

	buf := make([]byte, 256)
	for {
		n, err := http_resp.Body.Read(buf)
		if n > 0 {
			// FIXME: We need a max length limit.
			resp.data = append(resp.data, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error while reading body: %s", err.Error())
			break
		}
	}

	http_resp.Body.Close()

	return &resp, nil
}

type identity struct {
	Version string
	Client string
	Station_id string
	Secret string
}

func callIdentify(conn net.Conn) (*identity, error) {
	r, err := http.NewRequest("GET", "/Identify", nil)
	if err != nil {
		log.Printf("Error calling Identify: %s", err.Error())
		return nil, err
	}
	resp, err := doStationRPC(conn, r)
	if err != nil {
		log.Printf("RPC error: %s", err.Error())
		return nil, err
	}
	if resp.code != 200 {
		log.Printf("Bad Identify status code from station: %v",
			resp.code)
		return nil, errors.New("Bad Identify status code")
	}

	var id identity
	if err = json.Unmarshal(resp.data, &id); err != nil {
		log.Printf("Bad Identify json response from station: %s",
			err.Error())
		return nil, errors.New("Bad Identify json response")
	}

	// Avoid logging the secret.
	d := strings.Replace((string)(resp.data), id.Secret, "SECRET", -1)
	log.Printf("Identify info: %s", d)

	return &id, nil
}

func callDisconnect(conn net.Conn, reason string) {
	log.Printf("Disconnecting station because: %s", reason)
	u := url.URL{}
	u.Path = "/Disconnect"
	u.RawQuery = url.Values{"reason": {reason}}.Encode()
	r, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Printf("Error constructing disconnect url: %s",
			err.Error())
		return
	}
	doStationRPC(conn, r)
}

func callPing(conn net.Conn) bool {
	u := url.URL{}
	u.Path = "/Ping"
	r, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Printf("Error ping url: %s", err.Error())
		return false
	}
	resp, err := doStationRPC(conn, r)
	if err != nil {
		log.Printf("RPC error: %s", err.Error())
		return false
	}
	if resp.code != 200 {
		log.Printf("Ping error code: %d", resp.code)
		return false
	}
	return true
}

func handleStation(c *Coordinator, conn net.Conn) {
	defer conn.Close()
	log.Printf("Station connected.")

	id, err := callIdentify(conn)
	if err != nil {
		log.Printf("Error calling Identify: %s", err.Error())
		return
	}
	log.Printf("Station identified: %s, %s, %s\n",
		id.Station_id, id.Version, id.Client)

	ok, err := AuthenticateStation(c.sdb, id.Station_id, id.Secret)
	if err != nil {
		log.Printf("Authentication error: %s", err.Error())
		callDisconnect(conn, "internal server error")
		return
	}
	if !ok {
		log.Printf("Authentication denied for station.")
		callDisconnect(conn, "authentication denied")
		return
	}

	input := make(chan Request)
	c.stationConnected(id.Station_id, input)

	keepalive := time.NewTicker(keepAlivePingInterval)
	defer keepalive.Stop()

	for {
		select {
		case r := <-input:
			if r.request == nil {
				// nil request means we should disconnect
				callDisconnect(conn, "disconnected by server")
				c.stationDisconnected(id.Station_id)
				r.response <- nil
				return
			}

			resp, err := doStationRPC(conn, r.request)
			r.response <- resp
			if err != nil {
				log.Printf(
					"Error during station RPC: %s",
					err.Error())
				c.stationDisconnected(id.Station_id)
				return
			}

		case <- keepalive.C:
			if !callPing(conn) {
				callDisconnect(conn, "ping time out")
				c.stationDisconnected(id.Station_id)
				log.Printf("Station connection timed out.")
				return
			}
		}
	}
}

func ListenAndServe(c *Coordinator, cert_file, private_key, port string) {
	cert, err := tls.LoadX509KeyPair(cert_file, private_key)
	if err != nil {
		log.Fatalf("Error loading certs: %s", err.Error())
	}

	var config tls.Config
	config.Certificates = []tls.Certificate{cert}

	ln, err := tls.Listen("tcp", port, &config)
	if err != nil {
		log.Fatalf("Failed to start listening on port %s: %s",
			port, err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf(
				"Error accepting connection: %s", err.Error())
			continue
		}
		go handleStation(c, conn)
	}
}
