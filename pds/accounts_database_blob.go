package pds

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"path/filepath"
	"sync"

	"github.com/aaronland/gocloud/blob/bucket"
	"github.com/sfomuseum/go-atproto"
	"gocloud.dev/blob"
)

type BlobAccountsDatabase struct {
	AccountsDatabase
	bucket *blob.Bucket
}

var blob_accounts_register_mu = new(sync.RWMutex)
var blob_accounts_register_map = map[string]bool{}

func init() {

	ctx := context.Background()
	err := RegisterBlobAccountsSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterBlobAccountsSchemes will explicitly register all the schemes associated with ....
func RegisterBlobAccountsSchemes(ctx context.Context) error {

	blob_accounts_register_mu.Lock()
	defer blob_accounts_register_mu.Unlock()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {

		_, exists := blob_accounts_register_map[scheme]

		if exists {
			continue
		}

		err := RegisterAccountsDatabase(ctx, scheme, NewBlobAccountsDatabase)

		if err != nil {
			return fmt.Errorf("Failed to register blob accounts database for '%s', %w", scheme, err)
		}

		blob_accounts_register_map[scheme] = true
	}

	return nil
}

func NewBlobAccountsDatabase(ctx context.Context, uri string) (AccountsDatabase, error) {

	b, err := bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	db := &BlobAccountsDatabase{
		bucket: b,
	}

	return db, nil
}

func (db *BlobAccountsDatabase) GetAccount(ctx context.Context, did string) (*Account, error) {

	path := db.accountPath(did)

	exists, err := db.bucket.Exists(ctx, path)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, atproto.ErrNotFound
	}

	r, err := db.bucket.NewReader(ctx, path, nil)

	if err != nil {
		return nil, err
	}

	defer r.Close()

	var account *Account

	dec := json.NewDecoder(r)
	err = dec.Decode(&account)

	if err != nil {
		return nil, err
	}

	return account, err
}

func (db *BlobAccountsDatabase) GetAccountWithHandle(ctx context.Context, handle string) (*Account, error) {
	return nil, atproto.ErrNotImplemented
}

func (db *BlobAccountsDatabase) AddAccount(ctx context.Context, account *Account) error {
	return db.writeAccount(ctx, account)
}

func (db *BlobAccountsDatabase) UpdateAccount(ctx context.Context, account *Account) error {
	return db.writeAccount(ctx, account)
}

func (db *BlobAccountsDatabase) DeleteAccount(ctx context.Context, account *Account) error {
	return db.bucket.Delete(ctx, account.DID)
}

func (db *BlobAccountsDatabase) ListAccounts(ctx context.Context) iter.Seq2[*Account, error] {

	return func(yield func(*Account, error) bool) {
		yield(nil, atproto.ErrNotImplemented)
		return
	}
}

func (db *BlobAccountsDatabase) Close() error {
	return db.bucket.Close()
}

func (db *BlobAccountsDatabase) writeAccount(ctx context.Context, account *Account) error {

	path := db.accountPath(account.DID)

	wr, err := db.bucket.NewWriter(ctx, path, nil)

	if err != nil {
		return err
	}

	enc := json.NewEncoder(wr)
	err = enc.Encode(account)

	if err != nil {
		return err
	}

	return wr.Close()
}

func (db *BlobAccountsDatabase) accountPath(did string) string {
	fname := fmt.Sprintf("%s.json", did)
	path := filepath.Join("accounts", fname)
	return path
}
