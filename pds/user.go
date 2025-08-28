package pds

import (
	"context"
	"fmt"
	_ "net/http"
	"time"

	"github.com/bluesky-social/indigo/atproto/identity"
	_ "github.com/bluesky-social/indigo/plc"
	"github.com/sfomuseum/go-atproto/plc"
	_ "github.com/sfomuseum/go-atproto/plc/api"
)

type User struct {
	Id           string                `json:"id"`
	DID          *identity.DIDDocument `json:"did"`
	PrivateKey   string                `json:"private_key"`
	Handle       string                `json:"handle"`
	Aliases      []*Alias              `json:"aliases"`
	Created      int64                 `json:"created"`
	LastModified int64                 `json:"lastmodified"`
}

func CreateUser(ctx context.Context, service string, handle string) (*User, error) {

	rsp, err := plc.NewDID(ctx, service, handle)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new DID, %w", err)
	}

	doc := rsp.DID
	id := doc.DID.String()

	/*
		err = api.Create(ctx, id, rsp.PlcOperation)

		if err != nil {
			return nil, fmt.Errorf("Failed to create PLC for DID, %w", err)
		}
	*/

	// To do: Private key, wut??

	u := &User{
		Id:     id,
		DID:    doc,
		Handle: handle,
	}

	return u, nil
}

func GetUser(ctx context.Context, db UsersDatabase, did string) (*User, error) {
	return db.GetUser(ctx, did)
}

func GetUserWithHandle(ctx context.Context, db UsersDatabase, handle string) (*User, error) {
	return db.GetUserWithHandle(ctx, handle)
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
