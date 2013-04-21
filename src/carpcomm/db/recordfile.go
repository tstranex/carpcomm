// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

// Simple binary record file format.

package db

import "io"
import "errors"
import "log"
import "os"

const kRecordWriterV0Header = "RecordWriter000"

type RecordReader struct {
	r io.Reader
}

func NewRecordReader(r io.Reader) (*RecordReader, error) {
	header := make([]byte, len(kRecordWriterV0Header))
	_, err := io.ReadFull(r, header)
	if err != nil {
		return nil, err
	}
	if string(header) != kRecordWriterV0Header {
		return nil, errors.New("Invalid RecordWriter header.")
	}

	return &RecordReader{r}, nil
}

func NewRecordReaderForFile(path string) (*RecordReader, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("File open error: %s", err.Error())
		return nil, err
	}
	rr, err := NewRecordReader(f)
	if err != nil {
		log.Printf("RecordReader error: %s", err.Error())
		return nil, err
	}
	return rr, nil
}

func (rr *RecordReader) ReadRecord() ([]byte, error) {
	size := make([]byte, 4)
	_, err := io.ReadFull(rr.r, size)
	if err != nil {
		return nil, err
	}

	// Big-endian
	var n uint32
	for i := 0; i < 4; i++ {
		n = (n << 8) | (uint32)(size[i])
	}

	rec := make([]byte, n)
	_, err = io.ReadFull(rr.r, rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}