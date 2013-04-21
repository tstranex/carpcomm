// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package util

import "unicode"
import "strings"
import "errors"
import "fmt"

func StripWhitespace(s string) (r string) {
	for _, c := range s {
		if !unicode.IsSpace(c) {
			r += (string)(c)
		}
	}
	return r
}

func StripNonHex(s string) (r string) {
	s = strings.ToUpper(s)
	for _, c := range s {
		if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') {
			r += (string)(c)
		}
	}
	return r
}

// server.com:123 -> server.com, 123
func SplitHostAndPort(server string) (string, string, error) {
	parts := strings.SplitN(server, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New(
			fmt.Sprintf("Invalid host:port: %s", server))
	}
	return parts[0], parts[1], nil
}
