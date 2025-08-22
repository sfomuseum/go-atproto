package pds

import (
	"context"
	"time"
)

type UsersDatabase interface {
	GetUser(context.Context, string) (*User, error)
	AddUser(context.Context, *User) error
	UpdateUser(context.Context, *User) error
	DeleteUser(context.Context, *User) error
	Close() error
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
