package pds

import (
	"context"
	"time"

	at_crypto "github.com/bluesky-social/indigo/atproto/crypto"
	"github.com/sfomuseum/go-atproto/crypto"
)

type Key struct {
	DID                 string `json:"did"`
	Label               string `json:"label"`
	PrivateKeyMultibase string `json:"private_key_multibase"`
	Created             int64  `json:"created"`
	LastModified        int64  `json:"last_modified"`
}

func (k *Key) PrivateKeyK256() (*at_crypto.PrivateKeyK256, error) {
	return crypto.PrivateKeyK256FromMultibase(k.PrivateKeyMultibase)
}

func AddKey(ctx context.Context, db KeysDatabase, kp *Key) error {

	now := time.Now()
	ts := now.Unix()

	kp.Created = ts
	kp.LastModified = ts

	return db.AddKey(ctx, kp)
}

func DeleteKey(ctx context.Context, db KeysDatabase, k *Key) error {
	return db.DeleteKey(ctx, k)
}

func DeleteKeysForDID(ctx context.Context, db KeysDatabase, did string) error {
	return db.DeleteKeysForDID(ctx, did)
}
