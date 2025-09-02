package pds

import (
	"context"
	"iter"

	"github.com/sfomuseum/go-atproto"
)

type NullKeysDatabase struct {
	KeysDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterKeysDatabase(ctx, "null", NewNullKeysDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullKeysDatabase(ctx context.Context, uri string) (KeysDatabase, error) {

	db := &NullKeysDatabase{}
	return db, nil
}

func (db *NullKeysDatabase) GetKey(ctx context.Context, did string, label string) (*Key, error) {
	return nil, atproto.ErrNotFound
}

func (db *NullKeysDatabase) AddKey(ctx context.Context, kp *Key) error {
	return nil
}

func (db *NullKeysDatabase) DeleteKey(ctx context.Context, kp *Key) error {
	return nil
}

func (db *NullKeysDatabase) DeleteKeysForDID(ctx context.Context, did string) error {
	return nil
}

func (db *NullKeysDatabase) ListKeys(ctx context.Context, opts *ListKeysOptions) iter.Seq2[*Key, error] {
	return func(yield func(*Key, error) bool) {}
}

func (db *NullKeysDatabase) Close() error {
	return nil
}
