package main

import "carpcomm/telemetry"
import "flag"
import "fmt"
import "code.google.com/p/goprotobuf/proto"
import "encoding/hex"

func decodeHexFrame(satellite_id, data string) {
	frame, err := hex.DecodeString(data)
	if err != nil {
		fmt.Printf("Not a hex frame: %s\n", err.Error())
		return
	}

	pl, err := telemetry.DecodeFrame(satellite_id, frame, 0)
	if err != nil {
		fmt.Printf("DecodeFrame error: %s\n", err.Error())
		return
	}
	for _, p := range pl {
		fmt.Printf("%s\n", proto.MarshalTextString(&p))
	}
}

func decodeMorse(satellite_id, data string) {
	pl, err := telemetry.DecodeMorse(satellite_id, data, 0)
	if err != nil {
		fmt.Printf("DecodeMorse error: %s\n", err.Error())
		return
	}
	for _, p := range pl {
		fmt.Printf("%s\n", proto.MarshalTextString(&p))
	}
}

func main() {
	flag.Parse()
	satellite_id := flag.Args()[0]
	data := flag.Args()[1]

	decodeHexFrame(satellite_id, data)
	decodeMorse(satellite_id, data)
}