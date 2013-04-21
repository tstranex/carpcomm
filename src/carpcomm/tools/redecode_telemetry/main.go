package main

import "carpcomm/pb"
import "carpcomm/db"
import "carpcomm/telemetry"
import "code.google.com/p/goprotobuf/proto"

import "log"
import "flag"

var db_prefix = flag.String("db_prefix", "", "Database table prefix")
var contact_id = flag.String("contact_id", "", "Contact id")
var satellite_id = flag.String("satellite_id", "", "Satellite id")

func RedecodeContact(contactdb *db.ContactDB, c *pb.Contact) {
	if c == nil {
		log.Fatalf("Contact not found")
	}

	log.Printf("Original contact:\n%s\n", proto.MarshalTextString(c))

	new_blobs := make([]*pb.Contact_Blob, 0)
	var freeform []byte
	for _, b := range c.Blob {
		if b.Format != nil && (
			*b.Format == pb.Contact_Blob_DATUM ||
			*b.Format == pb.Contact_Blob_FRAME) {
			// Strip out datums and frames.
			continue
		}
		new_blobs = append(new_blobs, b)

		if b.Format != nil && *b.Format == pb.Contact_Blob_FREEFORM {
			if freeform != nil {
				log.Fatalf("Contact contains multiple FREEFORM blobs.")
			}
			freeform = b.InlineData
		}
	}

	if freeform == nil {
		return;
	}

	data, frames := telemetry.DecodeFreeform(
		*c.SatelliteId, freeform, *c.StartTimestamp)
	for i, _ := range frames {
		b := new(pb.Contact_Blob)
		b.Format = pb.Contact_Blob_FRAME.Enum()
		b.InlineData = frames[i]
		new_blobs = append(new_blobs, b)
	}
	for i, _ := range data {
		b := new(pb.Contact_Blob)
		b.Format = pb.Contact_Blob_DATUM.Enum()
		b.Datum = &data[i]
		new_blobs = append(new_blobs, b)
	}

	c.Blob = new_blobs

	log.Printf("New contact:\n%s\n", proto.MarshalTextString(c))

	err := contactdb.Store(c)
	if err != nil {
		log.Fatalf("Error storing contact: %s", err.Error())
	}
}

func main() {
	flag.Parse()

	domain, err := db.NewDomain(*db_prefix)
	if err != nil {
		log.Fatalf("Database error: %s", err.Error())
	}

	contactdb := domain.NewContactDB()
	if err := contactdb.Create(); err != nil {
		log.Fatalf("Error creating contact table: %s", err.Error())
	}

	if *contact_id != "" {
		c, err := contactdb.Lookup(*contact_id)
		if err != nil {
			log.Fatalf("Error looking up contact: %s", err.Error())
		}
		RedecodeContact(contactdb, c)

	} else if *satellite_id != "" {
		contacts, err := contactdb.SearchBySatelliteId(
			*satellite_id, 1000)
		if err != nil {
			log.Fatalf("Error looking up contacts: %s", err.Error())
		}
		for _, c := range contacts {
			RedecodeContact(contactdb, c)
		}

	} else {
		log.Fatalf("Missing option")
	}
}