package pds

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"os"

	"github.com/sfomuseum/go-atproto/did"
)

type User struct {
	DID          string   `json:"did"`
	PublicKey    string   `json:"public_key"`
	PrivateKey   string   `json:"private_key"`
	Handle       string   `json:"handle"`
	Aliases      []*Alias `json:"aliases"`
	Created      int64    `json:"created"`
	LastModified int64    `json:"lastmodified"`
}

func CreateUser(ctx context.Context, host string, name string) (*User, error) {

	d, prv_key, err := did.NewDID(name, host)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new DID, %w", err)
	}

	d.Marshal(os.Stdout)

	pub_key, err := d.PublicKey("#atproto")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive public key from DID, %w", err)
	}

	pub_b64 := base64.StdEncoding.EncodeToString(pub_key)
	prv_b64 := base64.StdEncoding.EncodeToString(prv_key)

	aka := d.AlsoKnownAs[0]
	handle := strings.TrimLeft(aka, "at://")

	u := &User{
		DID:        d.Id,
		Handle:     handle,
		PublicKey:  pub_b64,
		PrivateKey: prv_b64,
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
