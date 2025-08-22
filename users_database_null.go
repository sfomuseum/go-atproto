package pds

import (
	"context"
	"iter"
)

type NullUsersDatabase struct {
	UsersDatabase
}

func init() {

	ctx := context.Background()
	err := RegisterUsersDatabase(ctx, "null", NewNullUsersDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullUsersDatabase(ctx context.Context, uri string) (UsersDatabase, error) {

	db := &NullUsersDatabase{}
	return db, nil
}

func (db *NullUsersDatabase) GetUser(ctx context.Context, did string) (*User, error) {
	return nil, ErrNotFound
}

func (db *NullUsersDatabase) AddUser(ctx context.Context, user *User) error {
	return nil
}

func (db *NullUsersDatabase) UpdateUser(ctx context.Context, user *User) error {
	return nil

}

func (db *NullUsersDatabase) DeleteUser(ctx context.Context, user *User) error {
	return nil
}

func (db *NullUsersDatabase) ListUsers(ctx context.Context) iter.Seq2[*User, error] {

	return func(yield func(*User, error) bool) {
		return
	}
}

func (db *NullUsersDatabase) Close() error {
	return nil
}
