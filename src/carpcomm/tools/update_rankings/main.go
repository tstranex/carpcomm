package main

import "log"
import "carpcomm/db"
import "carpcomm/pb"
import "flag"
import "sort"
import "io/ioutil"
import "code.google.com/p/goprotobuf/proto"
import "os"

var input_file = flag.String("input_file", "", "Input filename")
var output_file = flag.String(
	"output_file", "data/rankings.RankingList", "Output filename")

type RankingList pb.RankingList

func (rl RankingList) Len() int {
	return len(rl.Ranking)
}

func (rl RankingList) Less(i, j int) bool {
	return *rl.Ranking[i].Score > *rl.Ranking[j].Score
}

func (rl RankingList) Swap(i, j int) {
	rl.Ranking[i], rl.Ranking[j] = rl.Ranking[j], rl.Ranking[i]
}


func main() {
	flag.Parse()

	f, err := os.Open(*input_file)
	if err != nil {
		log.Panicf("File open error: %s", err.Error())
	}
	rr, err := db.NewRecordReader(f)
	if err != nil {
		log.Panicf("RecordReader error: %s", err.Error())
	}

	items := make(map[string]*pb.Ranking)
	counts := make(map[string]map[pb.Contact_Blob_Format]int)
	n := 0
	for {
		b, err := rr.ReadRecord()
		if err != nil {
			break
		}
		n++

		c := &pb.Contact{}
		err = proto.Unmarshal(b, c)
		if err != nil {
			log.Panicf("Table error: %s", err.Error())
		}

		if c.UserId == nil {
			continue
		}

		userid := *c.UserId
		item, ok := items[userid]
		if !ok {
			item = &pb.Ranking{}
			item.UserId = &userid
			item.Score = proto.Int32(0)
			items[userid] = item

			counts[userid] = make(map[pb.Contact_Blob_Format]int)
		}

		for _, b := range c.Blob {
			*item.Score++
			counts[userid][*b.Format]++
		}
	}

	log.Printf("Read %d contacts.", n)

	ranked := &RankingList{}
	ranked.Ranking = make([]*pb.Ranking, len(items))
	i := 0
	for _, item := range items {
		ranked.Ranking[i] = item
		i++
	}
	sort.Sort(ranked)

	for _, item := range ranked.Ranking {
		for format, count := range counts[*item.UserId] {
			c := &pb.ContactCount{}
			c.Format = format.Enum()
			c.Count = proto.Int32((int32)(count))
			item.Counts = append(item.Counts, c)
		}
	}

	buf, err := proto.Marshal((*pb.RankingList)(ranked))
	if err != nil {
		log.Panicf("Marshal error: %s", err.Error())
	}
	ioutil.WriteFile(*output_file, buf, 0666)
}