package pds

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type UsersDatabase interface {
	GetUser(context.Context, string) (*User, error)
	AddUser(context.Context, *User) error
	UpdateUser(context.Context, *User) error
	DeleteUser(context.Context, *User) error
	ListUsers(context.Context) iter.Seq2[*User, error]
	Close() error
}

var user_database_roster roster.Roster

// UsersDatabaseInitializationFunc is a function defined by individual user_database package and used to create
// an instance of that user_database
type UsersDatabaseInitializationFunc func(ctx context.Context, uri string) (UsersDatabase, error)

// RegisterUsersDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `UsersDatabase` instances by the `NewUsersDatabase` method.
func RegisterUsersDatabase(ctx context.Context, scheme string, init_func UsersDatabaseInitializationFunc) error {

	err := ensureUsersDatabaseRoster()

	if err != nil {
		return err
	}

	return user_database_roster.Register(ctx, scheme, init_func)
}

func ensureUsersDatabaseRoster() error {

	if user_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		user_database_roster = r
	}

	return nil
}

// NewUsersDatabase returns a new `UsersDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `UsersDatabaseInitializationFunc`
// function used to instantiate the new `UsersDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterUsersDatabase` method.
func NewUsersDatabase(ctx context.Context, uri string) (UsersDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := user_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(UsersDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func UsersDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureUsersDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range user_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
