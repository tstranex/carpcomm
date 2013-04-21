// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package scheduler

import "log"
import "net/rpc"
import "time"
import "flag"
import "fmt"
import "sort"
import "carpcomm/db"
import "carpcomm/mux"
import "carpcomm/pb"
import "carpcomm/util"
import "carpcomm/util/timestamp"
import "carpcomm/streamer/contacts"

var streamer_address = flag.String(
	"streamer_address",
	"192.168.1.48:5050",
	"Streamer server address")

var api_server_address = flag.String(
	"api_server_address",
	"192.168.1.48:5051",
	"API server address")

const maxPredictionDuration = 10 * time.Hour
const controlDisabledDuration = 10 * time.Minute
const errorWaitDuration = 10  * time.Minute

func GetStreamURL(contact_id string) string {
	return fmt.Sprintf(
		"http://%s/%s", *streamer_address, contact_id)
}

func GetAPIServer() (string, string) {
	host, port, _ := util.SplitHostAndPort(*api_server_address)
	return host, port
}


func getNextPass(station *pb.Station) (pass Prediction, err error) {
	passes, err := PassPredictions(station)
	if err != nil {
		return pass, err
	}

	// Twiddle the predictions:

	twiddled := make(PredictionList, 0)

	// 1. Remove satellites with tracking disabled.
	for _, p := range passes {
		if p.Satellite.DisableTracking != nil &&
			*p.Satellite.DisableTracking == true {
			continue
		}
		twiddled = append(twiddled, p)
	}

	// 2. Impose maximum pass length.
	const maxPassLength = 5 * 60
	for _, p := range twiddled {
		d := p.EndTimestamp - p.StartTimestamp
		if d <= maxPassLength {
			continue
		}
		Δ := d - maxPassLength
		p.StartTimestamp += Δ/2
		p.EndTimestamp -= Δ/2
	}

	// Sort again since twiddling may affect the order.
	sort.Sort(twiddled)

	if len(twiddled) > 0 {
		return twiddled[0], nil
	}
	return pass, nil
}


// blocking
func capturePass(
	contactdb *db.ContactDB,
	mux_client *rpc.Client,
	station pb.Station,
	pass Prediction) error {

	log_label := *station.Id
	satellite_id := *pass.Satellite.Id
	duration := Duration(pass.EndTimestamp - pass.StartTimestamp)
	freq_hz := (int64)(*pass.CompatibleMode.channel.FrequencyHz)
	lateness := (float64)(time.Now().Unix()) - pass.StartTimestamp

	const kMotorProgramResolution = 10
	points, err := PassDetails(
		time.Now(),
		duration,
		*station.Lat,
		*station.Lng,
		*station.Elevation,
		*pass.Satellite.Tle,
		kMotorProgramResolution)
	if err != nil {
		log.Printf("%s: Error getting pass details: %s",
			log_label, err.Error())
		return err
	}
	motor_program := make([]mux.MotorCoordinate, len(points))
	for i, p := range points {
		motor_program[i].Timestamp = p.Timestamp
		motor_program[i].AzimuthDegrees = p.AzimuthDegrees
		motor_program[i].AltitudeDegrees = p.AltitudeDegrees
	}

	log.Printf("%s: Starting capture: satellite %s, "+
		"duration: %s, freq_hz: %d, lateness: %f [s]",
		log_label, satellite_id, duration, freq_hz, lateness)
	log.Printf("%s: motor_program: %v", log_label, motor_program)

	// FIXME: future: acquire lock

	log.Printf("%s: 1", log_label)
	contact, err := contacts.NewContact(
		&station, *station.Userid, &satellite_id)
	if err != nil {
		log.Printf("%s: Error creating contact: %s",
			log_label, err.Error())
		return err
	}
	log_label = *station.Id + "/" + *contact.Id

	log.Printf("%s: Contact id: %s", log_label, *contact.Id)

	log.Printf("%s: 2", log_label)
	err = mux.StationReceiverSetFrequency(mux_client, *station.Id, freq_hz)
	if err != nil {
		log.Printf("%s: StationReceiverSetFrequency failed: %s",
			log_label, err.Error())
		return err
	}

	log.Printf("%s: 3", log_label)
	err = contactdb.Store(contact)
	if err != nil {
		log.Printf("%s: Error storing contact: %s",
			log_label, err.Error())
		return err
	}

	log.Printf("%s: 4", log_label)
	stream_url := GetStreamURL(*contact.Id)
	err = mux.StationReceiverStart(mux_client, *station.Id, stream_url)
	if err != nil {
		log.Printf("%s: StationReceiverStart failed: %s",
			log_label, err.Error())
		//return err
	}

	// Stop the TNC in case it's already started.
	err = mux.StationTNCStop(mux_client, *station.Id)
	if err != nil {
		log.Printf("%s: StationTNCStop failed: %s",
			log_label, err.Error())
	}

	err = mux.StationTNCStart(
		mux_client, *station.Id, *api_server_address, satellite_id)
	if err != nil {
		log.Printf("%s: StationTNCStart failed: %s",
			log_label, err.Error())
	}

	log.Printf("%s: 5 StationMotorStart len=%d",
		log_label, len(motor_program))
	err = mux.StationMotorStart(mux_client, *station.Id, motor_program)
	if err != nil {
		log.Printf("%s: StationMotorStart failed: %s",
			log_label, err.Error())
	}

	log.Printf("%s: 6", log_label)
	time.Sleep(duration)

	log.Printf("%s: 7", log_label)

	err = mux.StationTNCStop(mux_client, *station.Id)
	if err != nil {
		log.Printf("%s: StationTNCStop failed: %s",
			log_label, err.Error())
	}

	err = mux.StationReceiverStop(mux_client, *station.Id)
	if err != nil {
		log.Printf("%s: StationReceiverStop failed: %s",
			log_label, err.Error())
	}

	err = mux.StationMotorStop(mux_client, *station.Id)
	if err != nil {
		log.Printf("%s: StationMotorStop failed: %s",
			log_label, err.Error())
	}

	// FIXME: send end timestamp and update contact
	now := time.Now().Unix()
	contact.EndTimestamp = &now
	

	// FIXME: future: release lock

	return nil
}

// blocking
// FIXME: handle disconnections
func scheduleStation(stationdb *db.StationDB,
	contactdb *db.ContactDB,
	mux_client *rpc.Client,
	station_id string,
	shutdown_chan chan string) {

	defer func() { shutdown_chan <- station_id }()

	log_label := station_id
	log.Printf("%s: Scheduling starting.", log_label)

	for {
		// Lookup the station in the loop since the owner might
		// edit the parameters.
		station, err := stationdb.Lookup(station_id)
		if err != nil {
			log.Printf("%s: Station lookup error: %s",
				log_label, err.Error())
			continue
		}
		if station == nil {
			log.Printf("%s: Station doesn't exist.", log_label)
			return
		}

		next_pass, err := getNextPass(station)
		if err != nil {
			log.Printf("%s: Error getting pass predictions: %s",
				log_label, err.Error())
		}

		var delay time.Duration
		if next_pass.Satellite == nil {
			// No passes coming up soon. Just wait a bit and try
			// again.
			log.Printf("%s: no upcoming passes", log_label)
			delay = maxPredictionDuration
		} else if  station.SchedulerEnabled == nil ||
			*station.SchedulerEnabled == false {
			//log.Printf("%s: scheduler disabled", station_id)
			// TODO: We are effectively polling the SchedulerEnabled
			// bit. We really should instead have the frontend
			// send a notification to reschedule the station.
			delay = controlDisabledDuration
		} else {
			log.Printf("%s: next pass: %s",
				log_label, *next_pass.Satellite.Id)
			delay = timestamp.TimestampFloatToTime(
				next_pass.StartTimestamp).Sub(time.Now())
		}

		log.Printf("%s: waiting for %s", log_label, delay)

		// wait {timer, disconnect}
		// FIXME: handle disconnections during the sleep
		time.Sleep(delay)

		if next_pass.Satellite == nil {
			continue
		}

		// Check that station has enabled scheduling.
		// We have to look up the station again since it might have
		// changed while we were sleeping.
		station, err = stationdb.Lookup(station_id)
		if err != nil {
			log.Printf("%s: Station lookup error: %s",
				log_label, err.Error())
			continue
		}
		if station == nil {
			log.Printf("%s: Station doesn't exist.", log_label)
			return
		}
		if station.SchedulerEnabled == nil ||
			*station.SchedulerEnabled == false {
			continue
		}

		err = capturePass(contactdb, mux_client, *station, next_pass)
		if err != nil {
			// There was an error of some sort.
			// Wait a bit before trying again.
			time.Sleep(errorWaitDuration)
		}
	}
}

func ScheduleForever(stationdb *db.StationDB,
	contactdb *db.ContactDB,
	mux_client *rpc.Client) {
	log.Printf("Scheduler started")

	active_stations := make(map[string]bool)
	shutdown_chan := make(chan string)

	for {
		for {
			all_handled := false
			select {
			case id := <- shutdown_chan:
				log.Printf("Station no longer active: %s", id)
				active_stations[id] = false
				continue;
			default:
				all_handled = true
			}
			if all_handled {
				break
			}
		}

		var args mux.StationListArgs
		var online_stations mux.StationListResult
		err := mux_client.Call(
			"Coordinator.StationList", args, &online_stations)
		if err != nil {
			log.Printf("Error calling mux: %s", err.Error())
			continue
		}

		for _, id := range online_stations.StationIds {
			if !active_stations[id] {
				go scheduleStation(
					stationdb, contactdb, mux_client,
					id, shutdown_chan)
				active_stations[id] = true
			}
		}

		time.Sleep(time.Minute)
	}
}
