package pds

import (
	"context"
)

type UsersDatabase interface {
	GetUser(context.Context, string) (*User, error)
	AddUser(context.Context, *User) error
	UpdateUser(context.Context, *User) error
	DeleteUser(context.Context, *User) error
}
