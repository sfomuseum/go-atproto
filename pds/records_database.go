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

type ListRecordsOptions struct {
	Repo       string
	Collection string
}

type RecordsDatabase interface {
	GetRecord(context.Context, string, string, string) (*Record, error)
	AddRecord(context.Context, *Record) error
	UpdateRecord(context.Context, *Record) error
	DeleteRecord(context.Context, *Record) error
	ListRecords(context.Context, *ListRecordsOptions) iter.Seq2[*Record, error]
	Close() error
}

var record_database_roster roster.Roster

// RecordsDatabaseInitializationFunc is a function defined by individual record_database package and used to create
// an instance of that record_database
type RecordsDatabaseInitializationFunc func(ctx context.Context, uri string) (RecordsDatabase, error)

// RegisterRecordsDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `RecordsDatabase` instances by the `NewRecordsDatabase` method.
func RegisterRecordsDatabase(ctx context.Context, scheme string, init_func RecordsDatabaseInitializationFunc) error {

	err := ensureRecordsDatabaseRoster()

	if err != nil {
		return err
	}

	return record_database_roster.Register(ctx, scheme, init_func)
}

func ensureRecordsDatabaseRoster() error {

	if record_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		record_database_roster = r
	}

	return nil
}

// NewRecordsDatabase returns a new `RecordsDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `RecordsDatabaseInitializationFunc`
// function used to instantiate the new `RecordsDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterRecordsDatabase` method.
func NewRecordsDatabase(ctx context.Context, uri string) (RecordsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := record_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(RecordsDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func RecordsDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureRecordsDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range record_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
