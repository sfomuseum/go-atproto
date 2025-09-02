package pds

import (
	"context"
	"database/sql"
	_ "encoding/json"
	"fmt"
	"iter"
	"net/url"
	_ "strings"

	_ "github.com/bluesky-social/indigo/atproto/identity"
	"github.com/sfomuseum/go-atproto"
)

type SQLAccountsDatabase struct {
	AccountsDatabase
	conn *sql.DB
}

func init() {

	ctx := context.Background()
	err := RegisterAccountsDatabase(ctx, "sql", NewSQLAccountsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLAccountsDatabase(ctx context.Context, uri string) (AccountsDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	engine := u.Host
	dsn := q.Get("dsn")

	if engine == "" {
		return nil, fmt.Errorf("Missing database engine")
	}

	if dsn == "" {
		return nil, fmt.Errorf("Missing DSN string")
	}

	conn, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Unable to create database (%s) because %v", engine, err)
	}

	db := &SQLAccountsDatabase{
		conn: conn,
	}

	return db, nil
}

func (db *SQLAccountsDatabase) GetAccount(ctx context.Context, did string) (*Account, error) {

	q := "SELECT did, handle, created, lastmodified FROM accounts where did = ?"
	return db.getAccount(ctx, q, did)
}

func (db *SQLAccountsDatabase) GetAccountWithHandle(ctx context.Context, handle string) (*Account, error) {

	q := "SELECT did, handle, created, lastmodified FROM accounts where handle = ?"
	return db.getAccount(ctx, q, handle)
}

func (db *SQLAccountsDatabase) getAccount(ctx context.Context, q string, args ...any) (*Account, error) {

	row := db.conn.QueryRowContext(ctx, q, args...)

	var did string
	var handle string
	var created int64
	var lastmod int64

	err := row.Scan(&did, &handle, &created, &lastmod)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, atproto.ErrNotFound
		}

		return nil, err
	}

	/*
		doc_r := strings.NewReader(str_doc)
		var doc *identity.DIDDocument

		dec := json.NewDecoder(doc_r)
		err = dec.Decode(&doc)

		if err != nil {
			return nil, fmt.Errorf("Failed to decode DID document, %w", err)
		}
	*/

	u := &Account{
		DID:          did,
		Handle:       handle,
		Created:      created,
		LastModified: lastmod,
	}

	return u, err
}

func (db *SQLAccountsDatabase) AddAccount(ctx context.Context, account *Account) error {

	/*
		enc_doc, err := json.Marshal(account.DIDDocument)

		if err != nil {
			return fmt.Errorf("Failed to encode DID document, %w", err)
		}
	*/

	q := "INSERT INTO accounts (did, handle, created, lastmodified) VALUES (?, ?, ?, ?)"

	_, err := db.conn.ExecContext(ctx, q, account.DID, account.Handle, account.Created, account.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add account, %w", err)
	}

	return nil
}

func (db *SQLAccountsDatabase) UpdateAccount(ctx context.Context, account *Account) error {

	/*
		enc_doc, err := json.Marshal(account.DIDDocument)

		if err != nil {
			return fmt.Errorf("Failed to encode DID document, %w", err)
		}
	*/

	q := "UPDATE accounts SET handle = ?, lastmodified = ? WHERE did = ?"

	_, err := db.conn.ExecContext(ctx, q, account.Handle, account.LastModified, account.DID)

	if err != nil {
		return fmt.Errorf("Failed to update account, %w", err)
	}

	return nil

}

func (db *SQLAccountsDatabase) ListAccounts(ctx context.Context) iter.Seq2[*Account, error] {

	return func(yield func(*Account, error) bool) {

		q := "SELECT did, handle, created, lastmodified FROM accounts ORDER BY created DESC"

		rows, err := db.conn.QueryContext(ctx, q)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {

			var did string
			var handle string
			var created int64
			var lastmod int64

			err := rows.Scan(&did, &handle, &created, &lastmod)

			if err != nil {

				if !yield(nil, err) {
					return
				}

				continue
			}

			/*
				doc_r := strings.NewReader(str_doc)
				var doc *identity.DIDDocument

				dec := json.NewDecoder(doc_r)
				err = dec.Decode(&doc)

				if err != nil {

					if !yield(nil, fmt.Errorf("Failed to decode DID document, %w", err)) {
						return
					}

					continue
				}
			*/

			u := &Account{
				DID:          did,
				Handle:       handle,
				Created:      created,
				LastModified: lastmod,
			}

			yield(u, nil)
		}

		err = rows.Close()

		if err != nil {
			yield(nil, err)
			return
		}

		err = rows.Err()

		if err != nil {
			yield(nil, err)
			return
		}
	}
}

func (db *SQLAccountsDatabase) Close() error {
	return db.conn.Close()
}
