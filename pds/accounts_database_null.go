package pds

import (
	"context"
	"iter"

	"github.com/sfomuseum/go-atproto"
)

type NullAccountsDatabase struct {
	AccountsDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterAccountsDatabase(ctx, "null", NewNullAccountsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullAccountsDatabase(ctx context.Context, uri string) (AccountsDatabase, error) {

	db := &NullAccountsDatabase{}
	return db, nil
}

func (db *NullAccountsDatabase) GetAccount(ctx context.Context, did string) (*Account, error) {
	return nil, atproto.ErrNotFound
}

func (db *NullAccountsDatabase) GetAccountWithHandle(ctx context.Context, handle string) (*Account, error) {
	return nil, atproto.ErrNotFound
}

func (db *NullAccountsDatabase) AddAccount(ctx context.Context, account *Account) error {
	return nil
}

func (db *NullAccountsDatabase) UpdateAccount(ctx context.Context, account *Account) error {
	return nil

}

func (db *NullAccountsDatabase) ListAccounts(ctx context.Context) iter.Seq2[*Account, error] {

	return func(yield func(*Account, error) bool) {
		return
	}
}

func (db *NullAccountsDatabase) Close() error {
	return nil
}
