// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "log"
import "fmt"
import "os"
import "sort"
import "time"

type fileList []os.FileInfo

func (fi fileList) Len() int {
	return len(fi)
}

func (fi fileList) Less(i, j int) bool {
	return fi[i].ModTime().After(fi[j].ModTime())
}

func (fi fileList) Swap(i, j int) {
	fi[i], fi[j] = fi[j], fi[i]
}

func garbageCollectDir(tmp_dir string, threshold_mb int) {
	f, err := os.Open(tmp_dir)
	if err != nil {
		log.Printf("Garbage collect error: %s", err.Error())
		return
	}

	fis, err := f.Readdir(0)
	if err != nil {
		log.Printf("Garbage collect error: %s", err.Error())
		return
	}

	sort.Sort((fileList)(fis))

	var total_size int64
	threshold := int64(threshold_mb) * 1024 * 1024
	for _, fi := range fis {
		total_size += fi.Size()
		path := fmt.Sprintf("%s/%s", *stream_tmp_dir, fi.Name())
		if total_size > threshold {
			log.Printf("Garbage collecting file: %s", path)
			err = os.Remove(path)
			if err != nil {
				log.Printf(
					"Error deleting file: %s", err.Error())
			}
		}
	}
}

func garbageCollectLoop(
	tmp_dir string, threshold_mb int, interval time.Duration) {
	log.Printf("Starting garbage collection loop for %s", tmp_dir)
	for {
		garbageCollectDir(tmp_dir, threshold_mb)
		time.Sleep(interval)
	}
}