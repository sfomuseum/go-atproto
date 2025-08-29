package pds

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/bluesky-social/indigo/atproto/identity"
	"github.com/sfomuseum/go-atproto/plc"
)

type Account struct {
	DID          string `json:"did"`
	Handle       string `json:"handle"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

type CreateAccountResponse struct {
	Account   *Account
	KeyPair   *KeyPair
	Operation *Operation
}

func CreateAccount(ctx context.Context, service string, handle string) (*CreateAccountResponse, error) {

	rsp, err := plc.NewDID(ctx, service, handle)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new DID, %w", err)
	}

	doc := rsp.DID
	did := doc.DID.String()

	op := rsp.Operation
	cid := op.CID().String()

	slog.Info("OK", "did", did, "cid", cid, "pk", rsp.PrivateKey.Multibase())

	// https://github.com/did-method-plc/go-didplc/blob/main/cmd/plcli/main.go#L286

	cl := plc.DefaultClient()

	err = cl.Submit(ctx, did, op)

	if err != nil {
		return nil, fmt.Errorf("Failed to submit operation, %w", err)
	}

	acct := &Account{
		DID:    did,
		Handle: handle,
	}

	acct_kp := &KeyPair{
		DID:                 did,
		Label:               "atproto",
		PrivateKeyMultibase: rsp.PrivateKey.Multibase(),
	}

	acct_op := &Operation{
		DID:       did,
		CID:       cid,
		Operation: op,
	}

	acct_rsp := &CreateAccountResponse{
		Account:   acct,
		KeyPair:   acct_kp,
		Operation: acct_op,
	}

	return acct_rsp, nil
}

func GetAccount(ctx context.Context, db AccountsDatabase, did string) (*Account, error) {
	return db.GetAccount(ctx, did)
}

func GetAccountWithHandle(ctx context.Context, db AccountsDatabase, handle string) (*Account, error) {
	return db.GetAccountWithHandle(ctx, handle)
}

func AddAccount(ctx context.Context, db AccountsDatabase, account *Account) error {

	now := time.Now()
	ts := now.Unix()

	account.Created = ts
	account.LastModified = ts

	return db.AddAccount(ctx, account)
}

func UpdateAccount(ctx context.Context, db AccountsDatabase, account *Account) error {

	now := time.Now()
	ts := now.Unix()

	account.LastModified = ts
	return db.AddAccount(ctx, account)
}

func DeleteAccount(ctx context.Context, db AccountsDatabase, account *Account) error {
	return db.DeleteAccount(ctx, account)
}
