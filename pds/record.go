package pds

import (
	"context"
	"fmt"
	"time"
)

type Record struct {
	CID          string `json:"cid"`
	DID          string `json:"did"`
	Collection   string `json:"collection"`
	RKey         string `json:"rkey"`
	Value        string `json:"value"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

func (r *Record) BlockURI() string {
	return fmt.Sprintf("repo:%s/%s/%s", r.DID, r.Collection, r.RKey)
}

func GetRecord(ctx context.Context, db RecordsDatabase, repo string, collection string, rkey string) (*Record, error) {
	return db.GetRecord(ctx, repo, collection, rkey)
}

func AddRecord(ctx context.Context, db RecordsDatabase, record *Record) error {

	now := time.Now()
	ts := now.Unix()

	record.Created = ts
	record.LastModified = ts

	return db.AddRecord(ctx, record)
}

func UpdateRecord(ctx context.Context, db RecordsDatabase, record *Record) error {

	now := time.Now()
	ts := now.Unix()

	record.LastModified = ts
	return db.AddRecord(ctx, record)
}

func DeleteRecord(ctx context.Context, db RecordsDatabase, record *Record) error {
	return db.DeleteRecord(ctx, record)
}
