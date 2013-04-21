// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package db

import "time"
import "fmt"
import "carpcomm/pb"
import "code.google.com/p/goprotobuf/proto"
import "reflect"

type CommentDB struct {
	table *SDBTable
}

const kCommentColumn = "pb.Comment"
const kCommentKeyUserId = "user_id"
const kCommentKeyObjectId = "object_id"
const kCommentKeyTimestamp = "timestamp"

func NewCommentDB(table *SDBTable) *CommentDB {
	return &CommentDB{table}
}

func (db *CommentDB) Store(comment *pb.Comment) error {
	values, err := encodeItem(kCommentColumn, comment)
	if err != nil {
		return err
	}

	if comment.UserId != nil {
		values[kCommentKeyUserId] = *comment.UserId
	}
	if comment.ObjectId != nil {
		values[kCommentKeyObjectId] = *comment.ObjectId
	}
	if comment.Timestamp != nil {
		values[kCommentKeyTimestamp] = fmt.Sprintf(
			"%016x", *comment.Timestamp)
	}

	return db.table.put(*comment.Id, values)
}

func (db *CommentDB) Lookup(id string) (*pb.Comment, error) {
	s := &pb.Comment{}
	found, err := db.table.getProto(id, kCommentColumn, s)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return s, nil
}

// Results are sorted by timestamp (newest first).
func (db *CommentDB) SearchByObjectId(object_id string) (
	[]*pb.Comment, error) {
	query := fmt.Sprintf(
		"select * from `%s` where `%s` = '%s' and `%s` is not null order by `%s` desc",
		db.table.domain.Name,
		kCommentKeyObjectId,
		object_id,
		kCommentKeyTimestamp,
		kCommentKeyTimestamp)
	result, err := db.table.search(
		query, kCommentColumn, reflect.TypeOf(pb.Comment{}))
	if err != nil {
		return nil, err
	}
	conv := make([]*pb.Comment, len(result))
	for i, v := range(result) {
		conv[i] = v.(*pb.Comment)
	}
	return conv, nil
}

func (db *CommentDB) GetAll() ([]*pb.Comment, error) {
	result, err := db.table.getAll(kCommentColumn,
		reflect.TypeOf(pb.Comment{}))
	if err != nil {
		return nil, err
	}
	conv := make([]*pb.Comment, len(result))
	for i, v := range(result) {
		conv[i] = v.(*pb.Comment)
	}
	return conv, nil
}

func (db *CommentDB) Create() error {
	return db.table.create()
}


// user_id may be nil for anonymous comments
func NewComment(object_id, text string, user_id *string) (
	*pb.Comment, error) {
	c := &pb.Comment{}

	id, err := CryptoRandId()
	if err != nil {
		return nil, err
	}
	c.Id = proto.String(id)
	c.UserId = user_id
	c.ObjectId = &object_id

	now := time.Now().Unix()
	c.Timestamp = &now

	c.Text = &text

	return c, nil
}
