package pds

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/aaronland/gocloud/blob/bucket"
	"gocloud.dev/blob"
)

type BlobRecordsDatabase struct {
	RecordsDatabase
	bucket *blob.Bucket
}

func NewBlobRecordsDatabase(ctx context.Context, uri string) (RecordsDatabase, error) {

	b, err := bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	db := &BlobRecordsDatabase{
		bucket: b,
	}

	return db, nil
}

func (db *BlobRecordsDatabase) GetRecord(ctx context.Context, repo string, collection string, rkey string) (*Record, error) {

	path := db.recordPath(repo, collection, rkey)
	r, err := db.bucket.NewReader(ctx, path, nil)

	if err != nil {
		return nil, err
	}

	defer r.Close()

	var record *Record

	dec := json.NewDecoder(r)
	err = dec.Decode(&record)

	if err != nil {
		return nil, err
	}

	return record, err
}

func (db *BlobRecordsDatabase) AddRecord(ctx context.Context, record *Record) error {
	return db.writeRecord(ctx, record)
}

func (db *BlobRecordsDatabase) UpdateRecord(ctx context.Context, record *Record) error {
	return db.writeRecord(ctx, record)
}

func (db *BlobRecordsDatabase) DeleteRecord(ctx context.Context, record *Record) error {
	path := db.recordPath(record.DID, record.Collection, record.RKey)
	return db.bucket.Delete(ctx, path)
}

func (db *BlobRecordsDatabase) Close() error {
	return db.bucket.Close()
}

func (db *BlobRecordsDatabase) writeRecord(ctx context.Context, record *Record) error {

	path := db.recordPath(record.DID, record.Collection, record.RKey)

	wr, err := db.bucket.NewWriter(ctx, path, nil)

	if err != nil {
		return err
	}

	enc := json.NewEncoder(wr)
	err = enc.Encode(record)

	if err != nil {
		return err
	}

	return wr.Close()
}

func (db *BlobRecordsDatabase) recordPath(repo string, collection string, rkey string) string {
	path := filepath.Join(repo, collection)
	path = filepath.Join(path, rkey)
	return path
}
