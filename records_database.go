package pds

import (
	"context"
	"time"
)

type RecordsDatabase interface {
	GetRecord(context.Context, string, string, string) (*Record, error)
	AddRecord(context.Context, *Record) error
	UpdateRecord(context.Context, *Record) error
	DeleteRecord(context.Context, *Record) error
	Close() error
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
