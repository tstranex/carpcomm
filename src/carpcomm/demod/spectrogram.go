// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package demod

import "carpcomm/pb"
import "log"
import "fmt"
import "os/exec"
import "time"

const spectrogramPath = "src/carpcomm/demod/spectrogram.py"

func SpectrogramTitle(c pb.Contact, iq_params pb.IQParams) string {
	var title string

	if c.StartTimestamp != nil {
		title += time.Unix(*c.StartTimestamp, 0).UTC().String() + " "
	}

	if c.SatelliteId != nil {
		title += *c.SatelliteId + " "
	}

	title += fmt.Sprintf("rate=%d type=%s id=%s",
		*iq_params.SampleRate, iq_params.Type.String(), *c.Id)
	return title
}

func Spectrogram(iq_path string,
	iq_params pb.IQParams,
	out_path string,
	title string) error {
	log.Printf("Creating spectrogram")

	c := exec.Command("python", spectrogramPath,
		iq_path,
		fmt.Sprintf("%d", *iq_params.SampleRate),
		iq_params.Type.String(),
		out_path,
		title)
	err := c.Run()
	if err != nil {
		log.Printf("Error running %s: %s", spectrogramPath, err.Error())
		return err
	}

	return nil
}