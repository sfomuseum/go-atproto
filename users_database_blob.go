package pds

import (
	"context"
	"encoding/json"

	"github.com/aaronland/gocloud/blob/bucket"
	"gocloud.dev/blob"
)

type BlobUsersDatabase struct {
	UsersDatabase
	bucket *blob.Bucket
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

	r, err := db.bucket.NewReader(ctx, did, nil)

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

func (db *BlobUsersDatabase) Close() error {
	return db.bucket.Close()
}

func (db *BlobUsersDatabase) writeUser(ctx context.Context, user *User) error {

	wr, err := db.bucket.NewWriter(ctx, user.DID, nil)

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
