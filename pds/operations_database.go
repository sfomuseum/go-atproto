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

type OperationsDatabase interface {
	GetOperation(context.Context, string) (*Operation, error)
	AddOperation(context.Context, *Operation) error
	ListOperations(context.Context) iter.Seq2[*Operation, error]
	Close() error
}

var operation_database_roster roster.Roster

// OperationsDatabaseInitializationFunc is a function defined by individual operation_database package and used to create
// an instance of that operation_database
type OperationsDatabaseInitializationFunc func(ctx context.Context, uri string) (OperationsDatabase, error)

// RegisterOperationsDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `OperationsDatabase` instances by the `NewOperationsDatabase` method.
func RegisterOperationsDatabase(ctx context.Context, scheme string, init_func OperationsDatabaseInitializationFunc) error {

	err := ensureOperationsDatabaseRoster()

	if err != nil {
		return err
	}

	return operation_database_roster.Register(ctx, scheme, init_func)
}

func ensureOperationsDatabaseRoster() error {

	if operation_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		operation_database_roster = r
	}

	return nil
}

// NewOperationsDatabase returns a new `OperationsDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `OperationsDatabaseInitializationFunc`
// function used to instantiate the new `OperationsDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterOperationsDatabase` method.
func NewOperationsDatabase(ctx context.Context, uri string) (OperationsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := operation_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(OperationsDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func OperationsDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureOperationsDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range operation_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
