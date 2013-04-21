package contacts

import "carpcomm/db"
import "carpcomm/pb"
import "carpcomm/telemetry"
import "net/http"
import "time"
import "log"
import "errors"
import "fmt"
import "code.google.com/p/goprotobuf/proto"


// satellite_id can be nil if it's unknown
// station can be null and user_id empty for anonymous contacts
func NewContact(station *pb.Station, user_id string, satellite_id *string) (
	*pb.Contact, error) {
	s := &pb.Contact{}

	id, err := db.CryptoRandId()
	if err != nil {
		return nil, err
	}
	s.Id = proto.String(id)

	if satellite_id != nil && *satellite_id != "" {
		s.SatelliteId = satellite_id
	}

	if user_id != "" {
		s.UserId = proto.String(user_id)
	}

	if station != nil {
		s.StationId = station.Id
		s.Lat = station.Lat
		s.Lng = station.Lng
		s.Elevation = station.Elevation
	}

	now := time.Now().Unix()
	s.StartTimestamp = &now

	return s, nil
}

// satellite_id can be empty if unknown
func StartNewConsoleContact(stationdb *db.StationDB, contactdb *db.ContactDB,
	station_id, user_id, satellite_id string) (
	id string, err error) {

	station, err := stationdb.Lookup(station_id)
	if err != nil {
		return "", err
	}
	if station == nil {
		return "", errors.New(
			fmt.Sprintf("Unknown station: %s", station_id))
	}

	s, err := NewContact(station, user_id, &satellite_id)
	if err != nil {
		return "", err
	}

	if err = contactdb.Store(s); err != nil {
		return "", err
	}

	return *s.Id, nil
}

type PopulateContactError struct {
	HttpStatusCode int
	PublicMessage string
}

func NewPopulateContactError(code int, public string) (
	*PopulateContactError) {
	return &PopulateContactError{code, public}
}

func (e *PopulateContactError) HttpError(w http.ResponseWriter) {
	http.Error(w, e.PublicMessage, e.HttpStatusCode)
}

func (e *PopulateContactError) Error() string {
	return e.PublicMessage
}

// If station is nil, an anonymous contact will be created.
// Otherwise, an authenticated contact will be created.
func PopulateContact(
	satellite_id string,
	timestamp int64,
	format_s string,
	data []byte,
	authenticated_user_id string,
	station_secret string,
	station *pb.Station) (*pb.Contact, *PopulateContactError) {

	if satellite_id == "" {
		return nil, NewPopulateContactError(
			http.StatusBadRequest, "Missing id")
	}

	if len(data) == 0 {
		return nil, NewPopulateContactError(
			http.StatusBadRequest, "Empty frame")
	} else if len(data) > 16384 {
		return nil, NewPopulateContactError(
			http.StatusBadRequest, "Exceeded size limit")
	}

	if time.Unix(timestamp, 0).After(time.Now()) {
		return nil, NewPopulateContactError(
			http.StatusBadRequest, "Timestamp is in the future.")
	}

	intformat, ok := pb.Contact_Blob_Format_value[format_s]
	if !ok {
		return nil, NewPopulateContactError(
			http.StatusBadRequest, "Can't parse format.")
	}
	format := (pb.Contact_Blob_Format)(intformat)

	// Make sure the satellite exists.
	sat := db.GlobalSatelliteDB().Map[satellite_id]
	if sat == nil {
		return nil, NewPopulateContactError(
			http.StatusBadRequest, "Unknown satellite")
	}

	// We get a lot of invalid frames produced by noise on KISS serial
	// lines. Filter them out early to avoid storing too much garbage in
	// the database.
	if format == pb.Contact_Blob_FRAME && !IsValidFrame(
		satellite_id, data) {
		return nil, NewPopulateContactError(
			http.StatusBadRequest, "Invalid data frame.")
	}

	var user_id string
	if station != nil {
		// Make an authenticated contact.
		// Ensure that either the station secret is correct or the user
		// owns the station.
		if station_secret != "" {
			if *station.Secret != station_secret {
				return nil, NewPopulateContactError(
					http.StatusUnauthorized, "")
			}
		} else {
			if *station.Userid != authenticated_user_id {
				return nil, NewPopulateContactError(
					http.StatusUnauthorized, "")
			}
		}
		user_id = *station.Userid
	} else {
		// Make an anonymous contact.
		user_id = ""
	}

	contact, err := NewContact(station, user_id, &satellite_id)
	if err != nil {
		log.Printf("Error creating contact: %s", err.Error())
		return nil, NewPopulateContactError(
			http.StatusInternalServerError, "")
	}

	contact.StartTimestamp = &timestamp

	blob := &pb.Contact_Blob{}
	blob.Format = format.Enum()
	blob.InlineData = data

	contact.Blob = []*pb.Contact_Blob{blob}

	blobs, err := DecodeBlobs(
		satellite_id, timestamp, []pb.Contact_Blob{*blob})
	if err != nil {
		log.Printf("Error decoding blobs: %s", err.Error())
	}
	contact.Blob = append(contact.Blob, blobs...)

	return contact, nil
}

// Try to decode morse messages and data frames.
// This should really be done asynchronously by the telemetry pipeline.
func DecodeBlobs(
	satellite_id string, timestamp int64, blobs []pb.Contact_Blob) (
	result []*pb.Contact_Blob, err error) {

	var datums []pb.TelemetryDatum
	var new_frames [][]byte
	for _, b := range blobs {
		if b.Format == nil {
			log.Printf("Missing blob format.")
			continue
		}

		if *b.Format ==  pb.Contact_Blob_FREEFORM {
			d, f := telemetry.DecodeFreeform(
				satellite_id, b.InlineData, timestamp)
			datums = append(datums, d...)
			new_frames = append(new_frames, f...)
		} else if *b.Format == pb.Contact_Blob_MORSE {
			d, _ := telemetry.DecodeMorse(
				satellite_id,
				(string)(b.InlineData),
				timestamp)
			datums = append(datums, d...)
		} else if *b.Format == pb.Contact_Blob_FRAME {
			d, _ := telemetry.DecodeFrame(
				satellite_id, b.InlineData, timestamp)
			datums = append(datums, d...)
		} else {
			log.Printf("Unknown blob format: %v", *b.Format)
			return nil, NewPopulateContactError(
				http.StatusBadRequest, "Unknown format.")
		}
	}

	for _, f := range new_frames {
		b := &pb.Contact_Blob{}
		b.Format = pb.Contact_Blob_FRAME.Enum()
		b.InlineData = []byte(f)
		result = append(result, b)
	}

	for _, d := range datums {
		b := &pb.Contact_Blob{}
		b.Format = pb.Contact_Blob_DATUM.Enum()
		b.Datum = &pb.TelemetryDatum{}
		*b.Datum = d
		result = append(result, b)
	}

	return result, nil
}