package pds

import (
	"context"
	"time"

	"github.com/sfomuseum/go-atproto"
)

type User struct {
	DID          string   `json:"did"`
	PublicKey    string   `json:"public_key"`
	Handle       *Handle  `json:"handle"`
	Aliases      []*Alias `json:"aliases"`
	Created      int64    `json:"created"`
	LastModified int64    `json:"lastmodified"`
}

func CreateUser(ctx context.Context) (*User, error) {
	return nil, atproto.ErrNotImplemented
}

func GetUser(ctx context.Context, db UsersDatabase, did string) (*User, error) {
	return db.GetUser(ctx, did)
}

func AddUser(ctx context.Context, db UsersDatabase, user *User) error {

	now := time.Now()
	ts := now.Unix()

	user.Created = ts
	user.LastModified = ts

	return db.AddUser(ctx, user)
}

func UpdateUser(ctx context.Context, db UsersDatabase, user *User) error {

	now := time.Now()
	ts := now.Unix()

	user.LastModified = ts
	return db.AddUser(ctx, user)
}

func DeleteUser(ctx context.Context, db UsersDatabase, user *User) error {
	return db.DeleteUser(ctx, user)
}
