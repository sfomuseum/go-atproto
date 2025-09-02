package pds

import (
	"context"
	"iter"

	"github.com/sfomuseum/go-atproto"
)

type NullOperationsDatabase struct {
	OperationsDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterOperationsDatabase(ctx, "null", NewNullOperationsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullOperationsDatabase(ctx context.Context, uri string) (OperationsDatabase, error) {
	db := &NullOperationsDatabase{}
	return db, nil
}

func (db *NullOperationsDatabase) GetOperation(ctx context.Context, cid string) (*Operation, error) {
	return nil, atproto.ErrNotFound
}

func (db *NullOperationsDatabase) GetLastOperationForDID(ctx context.Context, did string) (*Operation, error) {
	return nil, atproto.ErrNotFound
}

func (db *NullOperationsDatabase) AddOperation(ctx context.Context, op *Operation) error {
	return nil
}

func (db *NullOperationsDatabase) DeleteOperation(ctx context.Context, op *Operation) error {
	return nil
}

func (db *NullOperationsDatabase) ListOperations(ctx context.Context, opts *ListOperationsOptions) iter.Seq2[*Operation, error] {
	return func(yield func(*Operation, error) bool) {}
}

func (db *NullOperationsDatabase) Close() error {
	return nil
}
