package pds

import (
	"context"
	"time"
)

type KeyPair struct {
	DID                 string `json:"did"`
	Label               string `json:"label"`
	PublicKeyMultibase  string `json:"public_key_multibase"`
	PrivateKeyMultibase string `json:"private_key_multibase"`
	Created             int64  `json:"created"`
	LastModified        int64  `json:"last_modified"`
}

func AddKeyPair(ctx context.Context, db KeyPairsDatabase, kp *KeyPair) error {

	now := time.Now()
	ts := now.Unix()

	kp.Created = ts
	kp.LastModified = ts

	return db.AddKeyPair(ctx, kp)
}
