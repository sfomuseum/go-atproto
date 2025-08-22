package pds

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"path/filepath"
	"sync"

	"github.com/aaronland/gocloud/blob/bucket"
	"gocloud.dev/blob"
)

type BlobRecordsDatabase struct {
	RecordsDatabase
	bucket *blob.Bucket
}

var blob_records_register_mu = new(sync.RWMutex)
var blob_records_register_map = map[string]bool{}

func init() {

	ctx := context.Background()
	err := RegisterBlobRecordsSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterBlobRecordsSchemes will explicitly register all the schemes associated with ....
func RegisterBlobRecordsSchemes(ctx context.Context) error {

	blob_records_register_mu.Lock()
	defer blob_records_register_mu.Unlock()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {

		_, exists := blob_records_register_map[scheme]

		if exists {
			continue
		}

		err := RegisterRecordsDatabase(ctx, scheme, NewBlobRecordsDatabase)

		if err != nil {
			return fmt.Errorf("Failed to register blob records database for '%s', %w", scheme, err)
		}

		blob_records_register_map[scheme] = true
	}

	return nil
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

	exists, err := db.bucket.Exists(ctx, path)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrNotFound
	}

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

func (db *BlobRecordsDatabase) ListRecords(ctx context.Context, opts *ListRecordsOptions) iter.Seq2[*Record, error] {

	return func(yield func(*Record, error) bool) {
		yield(nil, ErrNotImplemented)
		return
	}
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

	fname := fmt.Sprintf("%s.json", rkey)

	path := filepath.Join("records", repo)
	path = filepath.Join(path, collection)
	path = filepath.Join(path, fname)
	return path
}
