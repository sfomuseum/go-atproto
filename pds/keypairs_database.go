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

type KeyPairsDatabase interface {
	GetKeyPair(context.Context, string, string) (*KeyPair, error)
	AddKeyPair(context.Context, *KeyPair) error
	ListKeyPairs(context.Context) iter.Seq2[*KeyPair, error]
	Close() error
}

var keypairs_database_roster roster.Roster

// KeyPairsDatabaseInitializationFunc is a function defined by individual keypairs_database package and used to create
// an instance of that keypairs_database
type KeyPairsDatabaseInitializationFunc func(ctx context.Context, uri string) (KeyPairsDatabase, error)

// RegisterKeyPairsDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `KeyPairsDatabase` instances by the `NewKeyPairsDatabase` method.
func RegisterKeyPairsDatabase(ctx context.Context, scheme string, init_func KeyPairsDatabaseInitializationFunc) error {

	err := ensureKeyPairsDatabaseRoster()

	if err != nil {
		return err
	}

	return keypairs_database_roster.Register(ctx, scheme, init_func)
}

func ensureKeyPairsDatabaseRoster() error {

	if keypairs_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		keypairs_database_roster = r
	}

	return nil
}

// NewKeyPairsDatabase returns a new `KeyPairsDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `KeyPairsDatabaseInitializationFunc`
// function used to instantiate the new `KeyPairsDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterKeyPairsDatabase` method.
func NewKeyPairsDatabase(ctx context.Context, uri string) (KeyPairsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := keypairs_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(KeyPairsDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func KeyPairsDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureKeyPairsDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range keypairs_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
