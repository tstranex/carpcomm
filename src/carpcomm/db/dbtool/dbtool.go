package main

import "code.google.com/p/goprotobuf/proto"

import "carpcomm/pb"
import "carpcomm/db"

import "log"
import "io"
import "flag"
import "fmt"
import "os"
import "io/ioutil"

var backup_dir = flag.String("backup_dir", "", "Backup directory")
var db_prefix = flag.String("db_prefix", "", "Database table prefix")
var table = flag.String("table", "", "Table name")
var id = flag.String("id", "", "Record id")

func RestoreUserTable(input *db.RecordReader, output *db.UserDB) error {
	for {
		rec, err := input.ReadRecord()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error reading record: %s", err.Error())
			return err
		}

		u := &pb.User{}
		err = proto.Unmarshal(rec, u)
		if err != nil {
			log.Printf("Error parsing record: %s", err.Error())
			return err
		}

		err = output.Store(u)
		if err != nil {
			log.Printf("Error writing record: %s", err.Error())
			return err
		}

		fmt.Printf(".")
	}
	fmt.Printf("\n")

	log.Printf("User table restored.")
	return nil
}

func RestoreCommentsTable(input *db.RecordReader, output *db.CommentDB) error {
	for {
		rec, err := input.ReadRecord()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error reading record: %s", err.Error())
			return err
		}

		c := &pb.Comment{}
		err = proto.Unmarshal(rec, c)
		if err != nil {
			log.Printf("Error parsing record: %s", err.Error())
			return err
		}

		err = output.Store(c)
		if err != nil {
			log.Printf("Error writing record: %s", err.Error())
			return err
		}

		fmt.Printf(".")
	}
	fmt.Printf("\n")

	log.Printf("Comments table restored.")
	return nil
}

func RestoreContactsTable(input *db.RecordReader, output *db.ContactDB) error {
	for {
		rec, err := input.ReadRecord()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error reading record: %s", err.Error())
			return err
		}

		c := &pb.Contact{}
		err = proto.Unmarshal(rec, c)
		if err != nil {
			log.Printf("Error parsing record: %s", err.Error())
			return err
		}

		err = output.Store(c)
		if err != nil {
			log.Printf("Error writing record: %s", err.Error())
			return err
		}

		fmt.Printf(".")
	}
	fmt.Printf("\n")

	log.Printf("Contacts table restored.")
	return nil
}

func RestoreStationsTable(input *db.RecordReader, output *db.StationDB) error {
	for {
		rec, err := input.ReadRecord()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error reading record: %s", err.Error())
			return err
		}

		s := &pb.Station{}
		err = proto.Unmarshal(rec, s)
		if err != nil {
			log.Printf("Error parsing record: %s", err.Error())
			return err
		}

		err = output.Store(s)
		if err != nil {
			log.Printf("Error writing record: %s", err.Error())
			return err
		}

		fmt.Printf(".")
	}
	fmt.Printf("\n")

	log.Printf("Stations table restored.")
	return nil
}

func getRecordReader(backup_dir string, name string) *db.RecordReader {
	rr, err := db.NewRecordReaderForFile(backup_dir + "/" + name)
	if err != nil {
		log.Panicf("Error opening record file: %s", err.Error())
	}
	return rr
}

func restore(domain* db.Domain) {
	user_rr := getRecordReader(*backup_dir, "test-users.User.rec")
	station_rr := getRecordReader(*backup_dir, "test-stations.Station.rec")
	comment_rr := getRecordReader(*backup_dir, "test-comments.Comment.rec")
	contact_rr := getRecordReader(*backup_dir, "test-contacts.Contact.rec")

	userdb := domain.NewUserDB()
	if err := userdb.Create(); err != nil {
		log.Fatalf("Error creating user table: %s", err.Error())
	}
	stationdb := domain.NewStationDB()
	if err := stationdb.Create(); err != nil {
		log.Fatalf("Error creating station table: %s", err.Error())
	}
	commentdb := domain.NewCommentDB()
	if err := commentdb.Create(); err != nil {
		log.Fatalf("Error creating comment table: %s", err.Error())
	}
	contactdb := domain.NewContactDB()
	if err := contactdb.Create(); err != nil {
		log.Fatalf("Error creating contact table: %s", err.Error())
	}
	
	if err := RestoreUserTable(user_rr, userdb); err != nil {
		log.Fatalf("Error restoring user table: %s", err.Error())
	}
	if err := RestoreStationsTable(station_rr, stationdb); err != nil {
		log.Fatalf("Error restoring station table: %s", err.Error())
	}
	if err := RestoreCommentsTable(comment_rr, commentdb); err != nil {
		log.Fatalf("Error restoring comments table: %s", err.Error())
	}
	if err := RestoreContactsTable(contact_rr, contactdb); err != nil {
		log.Fatalf("Error restoring contacts table: %s", err.Error())
	}
}

func lookup(domain* db.Domain) {
	if *table == "station" {
		stationdb := domain.NewStationDB()
		s, err := stationdb.Lookup(*id)
		if err != nil {
			log.Fatalf("Error looking up station: %s", err.Error());
		}
		proto.MarshalText(os.Stdout, s)
	} else if *table == "contact" {
		contactdb := domain.NewContactDB()
		s, err := contactdb.Lookup(*id)
		if err != nil {
			log.Fatalf("Error looking up contact: %s", err.Error());
		}
		proto.MarshalText(os.Stdout, s)
	} else {
		log.Fatalf("Unknown table: %s", *table)
	}
}

func store(domain* db.Domain) {
	if *table == "station" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Error reading proto: %s", err.Error());
		}

		var s pb.Station
		err = proto.UnmarshalText((string)(bytes), &s)
		if err != nil {
			log.Fatalf("Error parsing proto: %s", err.Error());
		}

		stationdb := domain.NewStationDB()
		err = stationdb.Store(&s)
		if err != nil {
			log.Fatalf("Error storing station: %s", err.Error());
		}

		log.Printf("Stored station.")
	} else if *table == "contact" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Error reading proto: %s", err.Error());
		}

		var s pb.Contact
		err = proto.UnmarshalText((string)(bytes), &s)
		if err != nil {
			log.Fatalf("Error parsing proto: %s", err.Error());
		}

		contactdb := domain.NewContactDB()
		err = contactdb.Store(&s)
		if err != nil {
			log.Fatalf("Error storing contact: %s", err.Error());
		}

		log.Printf("Stored contact.")
	} else {
		log.Fatalf("Unknown table: %s", *table)
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatalf("Missing command.")
	}

	domain, err := db.NewDomain(*db_prefix)
	if err != nil {
		log.Fatalf("Database error: %s", err.Error())
	}

	cmd := flag.Args()[0]
	if cmd == "restore" {
		restore(domain)
	} else if cmd == "lookup" {
		lookup(domain)
	} else if cmd == "store" {
		store(domain)
	} else {
		log.Fatalf("Unknown command: %s", cmd)
	}
}