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

type SQLUsersDatabase struct {
	UsersDatabase
	conn *sql.DB
}

func init() {

	ctx := context.Background()
	err := RegisterUsersDatabase(ctx, "sql", NewSQLUsersDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLUsersDatabase(ctx context.Context, uri string) (UsersDatabase, error) {

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

	db := &SQLUsersDatabase{
		conn: conn,
	}

	return db, nil
}

func (db *SQLUsersDatabase) GetUser(ctx context.Context, did string) (*User, error) {

	q := "SELECT handle, created, lastmodified FROM users where did = ?"

	row := db.conn.QueryRowContext(ctx, q, did)

	var str_doc string
	var handle string
	var created int64
	var lastmod int64

	err := row.Scan(&str_doc, &handle, &created, &lastmod)

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

	u := &User{
		DID:          did,
		Handle:       handle,
		Created:      created,
		LastModified: lastmod,
	}

	return u, err
}

func (db *SQLUsersDatabase) AddUser(ctx context.Context, user *User) error {

	/*
		enc_doc, err := json.Marshal(user.DIDDocument)

		if err != nil {
			return fmt.Errorf("Failed to encode DID document, %w", err)
		}
	*/

	q := "INSERT INTO users (did, handle, created, lastmodified) VALUES (?, ?, ?, ?)"

	_, err := db.conn.ExecContext(ctx, q, user.DID, user.Handle, user.Created, user.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add user, %w", err)
	}

	return nil
}

func (db *SQLUsersDatabase) UpdateUser(ctx context.Context, user *User) error {

	/*
		enc_doc, err := json.Marshal(user.DIDDocument)

		if err != nil {
			return fmt.Errorf("Failed to encode DID document, %w", err)
		}
	*/

	q := "UPDATE users SET handle = ?, lastmodified = ? WHERE did = ?"

	_, err := db.conn.ExecContext(ctx, q, user.Handle, user.LastModified, user.DID)

	if err != nil {
		return fmt.Errorf("Failed to update user, %w", err)
	}

	return nil

}

func (db *SQLUsersDatabase) DeleteUser(ctx context.Context, user *User) error {

	q := "DELETE FROM users where did = ?"

	_, err := db.conn.ExecContext(ctx, q, user.DID)

	if err != nil {
		return fmt.Errorf("Failed to delete user, %w", err)
	}

	return nil
}

func (db *SQLUsersDatabase) ListUsers(ctx context.Context) iter.Seq2[*User, error] {

	return func(yield func(*User, error) bool) {

		q := "SELECT did, handle, created, lastmodified FROM users ORDER BY created DESC"

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

			err := rows.Scan(&did, handle, &created, &lastmod)

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

			u := &User{
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

func (db *SQLUsersDatabase) Close() error {
	return db.conn.Close()
}
