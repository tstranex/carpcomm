// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "net/http"
import "log"
import "fmt"
import "os"
import "io"
import "strconv"
import "time"
import "carpcomm/db"
import "carpcomm/pb"
import "flag"
import "code.google.com/p/goprotobuf/proto"
import "strings"

var cert_file = flag.String(
	"cert_file",
	"/Users/tstranex/carp/carp/frontend/cert.pem",
	"SSL server certificate file")
var private_key_file = flag.String(
	"private_key_file",
	"/Users/tstranex/carp/carp/frontend/private_key.pem",
	"SSL private key file")

var stream_tmp_dir = flag.String(
	"stream_tmp_dir",
	"/tmp/streamer",
	"Directory for temporary stream files")
var port = flag.String("port", ":5050", "External port")
var tls_port = flag.String("tls_port", ":5051", "External port with TLS")
var db_prefix = flag.String("db_prefix", "r1-", "Database table prefix")
var gc_threshold_mb = flag.Int(
	"gc_threshold_mb", 3000, "Garbage collection threshold in MB")

type Handler struct {
	contactdb *db.ContactDB
	queue IQProcessingQueue
}

func NewHandler(contactdb *db.ContactDB, queue IQProcessingQueue) *Handler {
	return &Handler{contactdb, queue}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s", r.Method, r.URL.Path)

	if r.Method == "PUT" {
		h.Put(w, r)
	} else if r.Method == "GET" {
		h.Get(w, r)
	} else {
		http.Error(w, "Expected PUT method.", http.StatusBadRequest)
		return
	}
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]

	var iq_params pb.IQParams

	rate_s := r.URL.Query().Get("rate")
	if rate_s == "" {
		rate_s = "250977"
	}
	rate, err := strconv.Atoi(rate_s)
	if err != nil || rate <= 0 {
		http.Error(w, "Invalid 'rate' param.", http.StatusBadRequest)
		return
	}
	iq_params.SampleRate = proto.Int32((int32)(rate))

	type_s := r.URL.Query().Get("type")
	if type_s == "" {
		type_s = "UINT8"
	}
	inttype, ok := pb.IQParams_Type_value[type_s]
	if !ok {
		http.Error(w, "Invalid 'type' param.", http.StatusBadRequest)
		return
	}
	iq_params.Type = (pb.IQParams_Type)(inttype).Enum()
	
	c, err := h.contactdb.Lookup(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if c == nil {
		http.NotFound(w, r)
		return
	}

	// TODO(tstranex): If IQ data has already been uploaded, we should
	// prevent it from uploaded again.

	// TODO(tstranex): Use s3 multipart uploading instead. That will allow
	// us to get rid of the local file and the size limit.

	begin_time := time.Now()
	local_path := fmt.Sprintf("%s/%s", *stream_tmp_dir, id)
	file, err := os.Create(local_path)
	if err != nil {
		log.Printf("Error opening %s: %v", local_path, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	const size_limit = 256 * 1024 * 1024  // [bytes]
	written, _ := io.Copy(file, io.LimitReader(r.Body, size_limit))
	file.Close()

	d := time.Now().Sub(begin_time)
	upload_rate := float64(written) / d.Seconds() / 1024.0

	log.Printf("%s: Read %d bytes in %s: %f kB/s",
		id, written, d.String(), upload_rate)
	w.WriteHeader(http.StatusNoContent)

	// TODO: upload to s3

	iq_blob := &(pb.Contact_Blob{})
	iq_blob.Format = pb.Contact_Blob_IQ.Enum()
	iq_blob.IqParams = &iq_params
	c.Blob = append(c.Blob, iq_blob)

	// TODO: Need to be careful about locking the ContactDB record.
	// Currently we are the only process that modifies an existing ContactDB
	// record. In the future we may need to be more careful.
	if err := h.contactdb.Store(c); err != nil {
		log.Printf("Error storing contact: %s", err.Error())
		return
	}

	h.queue <- id
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]

	// FIXME: add authentication and verify id

	// For some reason, mime type detection assigns text/plain for IQ data,
	// which causes the browser to display it instead of download it.
	// So set the content-type explicitly.
	if !strings.Contains(id, ".") {
		w.Header().Add("Content-Type", "application/octet-stream")
	}

	local_path := fmt.Sprintf("%s/%s", *stream_tmp_dir, id)
	http.ServeFile(w, r, local_path)
}

func listenAndServeUploader(contactdb *db.ContactDB, queue IQProcessingQueue) {
	h := NewHandler(contactdb, queue)
	err := http.ListenAndServe(*port, h)
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
		return
	}
}

func main() {
	flag.Parse()

	domain, err := db.NewDomain(*db_prefix)
	if err != nil {
		log.Fatalf("Database error: %s", err.Error())
	}
	contactdb := domain.NewContactDB()
	stationdb := domain.NewStationDB()

	queue := make(IQProcessingQueue)
	go ProcessIQQueue(queue, contactdb)
	go listenAndServeUploader(contactdb, queue)
	go garbageCollectLoop(*stream_tmp_dir, *gc_threshold_mb, time.Minute)

	AddPacketHttpHandlers(http.DefaultServeMux, contactdb, stationdb)

	log.Printf("Starting streamer server")
	err = http.ListenAndServeTLS(
		*tls_port, *cert_file, *private_key_file, nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
		return
	}
}