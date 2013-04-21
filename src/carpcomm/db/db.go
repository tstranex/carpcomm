// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "launchpad.net/goamz/aws"
import "launchpad.net/goamz/exp/sdb"

import "code.google.com/p/goprotobuf/proto"

import "encoding/base64"
import "fmt"
import "reflect"
import "errors"
import "strconv"
import "math/rand"
import "math"
import "log"
import "time"

const maxRetries = 5


type SDBTable struct {
	s *sdb.SDB
	domain *sdb.Domain
}

func NewSDBTable(auth *aws.Auth, region *aws.Region, domain string) *SDBTable {
	t := &SDBTable{}
	t.s = sdb.New(*auth, *region)
	t.domain = t.s.Domain(domain)
	return t
}

// Keep in sync with _DecodeItem in carpcomm/tools/table.py.
func decodeItem(attrs []sdb.Attr, column string, p proto.Message) (
	found bool, err error) {

	table := make(map[string]string)
	for _, v := range(attrs) {
		table[v.Name] = v.Value
	}

	column_count := table[column + ".v2"]
	encoded_data := ""
	if column_count == "" {
		// v1 format
		encoded_data = table[column]
	} else {
		// v2 format
		count, err := strconv.Atoi(column_count)
		if err != nil {
			return false, err
		}
		for i := 0; i < count; i++ {
			encoded_data += table[column + "." + strconv.Itoa(i)]
		}
	}

	if encoded_data == "" {
		return false, nil
	}

	data, err := base64.StdEncoding.DecodeString(encoded_data)
	if err != nil {
		return false, err
	}
	if err = proto.Unmarshal(data, p); err != nil {
		return false, err
	}
	return true, nil
}

func encodeItem(column string, p proto.Message) (map[string]string, error) {
	data, err := proto.Marshal(p)
	if err != nil {
		return nil, err
	}
	v := base64.StdEncoding.EncodeToString(data)

	values := make(map[string]string)

	// SDB has a limit of 1024 bytes per value.
	if len(v) <= 1024 {
		// It's small enough to fit in the v1 format.
		// TODO(tstranex): Remove this optimization. Rather always use
		// the v2 format to reduce complexity.
		values[column] = v
	} else {
		// We need to use v2 format to split it over several columns.
		count := 0
		for ; len(v) > 0; {
			s := ""
			if len(v) > 1024 {
				s = v[:1024]
			} else {
				s = v
			}
			values[column + "." + strconv.Itoa(count)] = s
			v = v[len(s):]
			count++
		}
		values[column + ".v2"] = strconv.Itoa(count)
	}

	return values, nil
}


// Returns false, nil if there is no value for the id and column.
func (table *SDBTable) getProto(id, column string, p proto.Message) (
	found bool, err error) {
	item := table.domain.Item(id)
	resp, err := item.Attrs(nil, true)
	if err != nil {
		return false, err
	}
	return decodeItem(resp.Attrs, column, p)
}

func (table *SDBTable) put(id string, values map[string]string) error {
	var attrs sdb.PutAttrs
	for k, v := range values {
		attrs.Replace(k, v)
	}

	item := table.domain.Item(id)

	var err error
	for retry := 0; retry < maxRetries; retry++ {
		_, err := item.PutAttrs(&attrs)
		if err == nil {
			return nil
		}

		// Wait a bit before trying again with exponential backoff.
		ms := int(rand.Float64() * math.Pow(4, float64(retry)) * 100)
		log.Printf("sdb retry: %d, %d", retry, ms)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
	return err
}

func (table *SDBTable) setProto(id, column string, p proto.Message) error {
	values, err := encodeItem(column, p)
	if err != nil {
		return err
	}
	return table.put(id, values)
}

func (table *SDBTable) delete(id, column string) error {
	item := table.domain.Item(id)
	_, err := item.DeleteAttrNames([]string{column})
	return err
}

func (table *SDBTable) search(query, column string, t reflect.Type) (
	[]proto.Message, error) {
	resp, err := table.domain.Select(query, true, nil)
	if err != nil {
		return nil, err
	}

	result := make([]proto.Message, len(resp.Items))
	i := 0
	for _, v := range(resp.Items) {
		p := reflect.New(t).Interface().(proto.Message)
		found, err := decodeItem(v.Attrs, column, p)
		if err != nil {
			return nil, err
		}
		if !found {
			continue
		}
		result[i] = p
		i++
	}
	return result[:i], nil
}

// Create the database table.
func (table *SDBTable) create() error {
	_, err := table.domain.CreateDomain()
	return err
}


type Iterator struct {
	table *SDBTable
	query string
	consistent bool
	column string
	resp *sdb.SelectResp
	i int
	fetchError error
}

func (it *Iterator) fetch() error {
	var nextToken *string = nil
	if it.resp != nil {
		nextToken = &it.resp.NextToken
	}
	resp, err := it.table.domain.Select(it.query, it.consistent, nextToken)
	if err != nil {
		it.fetchError = err
		return err
	}
	it.resp = resp
	it.i = 0
	return nil
}

func (it *Iterator) Done() bool {
	return it.resp.NextToken == "" && it.i == len(it.resp.Items)
}

func (it *Iterator) Get() ([]byte, error) {
	if it.fetchError != nil {
		return nil, it.fetchError
	}
	if it.Done() {
		return nil, errors.New("Already done")
	}

	for _, v := range it.resp.Items[it.i].Attrs {
		if v.Name == it.column {
			data, err := base64.StdEncoding.DecodeString(v.Value)
			if err != nil {
				return nil, err
			}
			return data, nil
		 }
	}

	return nil, errors.New(fmt.Sprintf("Column missing: %s", it.column))
}

func (it *Iterator) Next() {
	it.i++
	if it.Done() {
		return
	}
	if it.i == len(it.resp.Items) {
		it.fetch()
	}
}

func (table *SDBTable) searchIterator(query, column string) (Iterator, error) {
	it := Iterator{table, query, true, column, nil, 0, nil}
	err := it.fetch()
	return it, err
}



func (table *SDBTable) getAll(column string, t reflect.Type) (
	[]proto.Message, error) {
	query := fmt.Sprintf("select `%s` from `%s`", column, table.domain.Name)
	return table.search(query, column, t)
}

func (table *SDBTable) GetAllIterator(column string) (Iterator, error) {
	query := fmt.Sprintf("select `%s` from `%s`", column, table.domain.Name)
	return table.searchIterator(query, column)
}

func (table *SDBTable) lookupByKeyValue(
	column, key, value string, t reflect.Type) ([]proto.Message, error) {
	query := fmt.Sprintf("select `%s` from `%s` where `%s` = '%s'",
		column, table.domain.Name, key, value)
	return table.search(query, column, t)
}