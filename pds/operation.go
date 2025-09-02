package pds

import (
	"context"
	"time"

	"github.com/did-method-plc/go-didplc"
)

type Operation struct {
	CID          string           `json:"cid"`
	DID          string           `json:"did"`
	Operation    didplc.Operation `json:"operation"`
	Created      int64            `json:"created"`
	LastModified int64            `json:"lastmodified"`
}

func AddOperation(ctx context.Context, db OperationsDatabase, op *Operation) error {

	now := time.Now()
	ts := now.Unix()

	op.Created = ts
	op.LastModified = ts

	return db.AddOperation(ctx, op)
}
