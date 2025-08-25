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

type BlobUsersDatabase struct {
	UsersDatabase
	bucket *blob.Bucket
}

var blob_users_register_mu = new(sync.RWMutex)
var blob_users_register_map = map[string]bool{}

func init() {

	ctx := context.Background()
	err := RegisterBlobUsersSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterBlobUsersSchemes will explicitly register all the schemes associated with ....
func RegisterBlobUsersSchemes(ctx context.Context) error {

	blob_users_register_mu.Lock()
	defer blob_users_register_mu.Unlock()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {

		_, exists := blob_users_register_map[scheme]

		if exists {
			continue
		}

		err := RegisterUsersDatabase(ctx, scheme, NewBlobUsersDatabase)

		if err != nil {
			return fmt.Errorf("Failed to register blob users database for '%s', %w", scheme, err)
		}

		blob_users_register_map[scheme] = true
	}

	return nil
}

func NewBlobUsersDatabase(ctx context.Context, uri string) (UsersDatabase, error) {

	b, err := bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	db := &BlobUsersDatabase{
		bucket: b,
	}

	return db, nil
}

func (db *BlobUsersDatabase) GetUser(ctx context.Context, did string) (*User, error) {

	path := db.userPath(did)

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

	var user *User

	dec := json.NewDecoder(r)
	err = dec.Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, err
}

func (db *BlobUsersDatabase) AddUser(ctx context.Context, user *User) error {
	return db.writeUser(ctx, user)
}

func (db *BlobUsersDatabase) UpdateUser(ctx context.Context, user *User) error {
	return db.writeUser(ctx, user)
}

func (db *BlobUsersDatabase) DeleteUser(ctx context.Context, user *User) error {
	return db.bucket.Delete(ctx, user.DID)
}

func (db *BlobUsersDatabase) ListUsers(ctx context.Context) iter.Seq2[*User, error] {

	return func(yield func(*User, error) bool) {
		yield(nil, atproto.ErrNotImplemented)
		return
	}
}

func (db *BlobUsersDatabase) Close() error {
	return db.bucket.Close()
}

func (db *BlobUsersDatabase) writeUser(ctx context.Context, user *User) error {

	path := db.userPath(user.DID)

	wr, err := db.bucket.NewWriter(ctx, path, nil)

	if err != nil {
		return err
	}

	enc := json.NewEncoder(wr)
	err = enc.Encode(user)

	if err != nil {
		return err
	}

	return wr.Close()
}

func (db *BlobUsersDatabase) userPath(did string) string {
	fname := fmt.Sprintf("%s.json", did)
	path := filepath.Join("users", fname)
	return path
}
