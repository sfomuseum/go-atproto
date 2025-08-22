package pds

import (
	"context"
	"iter"
)

type NullRecordsDatabase struct {
	RecordsDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterRecordsDatabase(ctx, "null", NewNullRecordsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullRecordsDatabase(ctx context.Context, uri string) (RecordsDatabase, error) {

	db := &NullRecordsDatabase{}
	return db, nil
}

func (db *NullRecordsDatabase) GetRecord(ctx context.Context, repo string, collection string, rkey string) (*Record, error) {

	return nil, ErrNotFound
}

func (db *NullRecordsDatabase) AddRecord(ctx context.Context, record *Record) error {
	return nil
}

func (db *NullRecordsDatabase) UpdateRecord(ctx context.Context, record *Record) error {
	return nil

}

func (db *NullRecordsDatabase) DeleteRecord(ctx context.Context, record *Record) error {
	return nil
}

func (db *NullRecordsDatabase) ListRecords(ctx context.Context, opts *ListRecordsOptions) iter.Seq2[*Record, error] {

	return func(yield func(*Record, error) bool) {
		return
	}
}

func (db *NullRecordsDatabase) Close() error {
	return nil
}
