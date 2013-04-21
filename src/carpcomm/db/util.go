// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "crypto/rand"
import "io"
import "fmt"

func BytesToInt63(bytes []byte) int64 {
	var r uint64 = 0
	for _, b := range(bytes) {
		r = (r<<8) | (uint64)(b)
	}
	r = r >> 1
	return (int64)(r)
}

func CryptoRandInt63() (int64, error) {
	bytes := make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		return 0, err
	}
	return BytesToInt63(bytes), nil
}

func CryptoRandId() (string, error) {
	id, err := CryptoRandInt63()
	return fmt.Sprintf("%d", id), err
}
