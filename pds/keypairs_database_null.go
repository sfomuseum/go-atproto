package pds

import (
	"context"
	"iter"

	"github.com/sfomuseum/go-atproto"
)

type NullKeyPairsDatabase struct {
	KeyPairsDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterKeyPairsDatabase(ctx, "null", NewNullKeyPairsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullKeyPairsDatabase(ctx context.Context, uri string) (KeyPairsDatabase, error) {

	db := &NullKeyPairsDatabase{}
	return db, nil
}

func (db *NullKeyPairsDatabase) GetKeyPair(ctx context.Context, did string, label string) (*KeyPair, error) {
	return nil, atproto.ErrNotFound
}

func (db *NullKeyPairsDatabase) AddKeyPair(ctx context.Context, kp *KeyPair) error {
	return nil
}

func (db *NullKeyPairsDatabase) DeleteKeyPair(ctx context.Context, kp *KeyPair) error {
	return nil
}

func (db *NullKeyPairsDatabase) ListKeyPairs(ctx context.Context, opts *ListKeyPairsOptions) iter.Seq2[*KeyPair, error] {
	return func(yield func(*KeyPair, error) bool) {}
}

func (db *NullKeyPairsDatabase) Close() error {
	return nil
}
