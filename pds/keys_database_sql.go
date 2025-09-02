package pds

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"net/url"

	"github.com/sfomuseum/go-atproto"
)

type SQLKeysDatabase struct {
	KeysDatabase
	conn   *sql.DB
	engine string
}

func init() {

	ctx := context.Background()
	err := RegisterKeysDatabase(ctx, "sql", NewSQLKeysDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLKeysDatabase(ctx context.Context, uri string) (KeysDatabase, error) {

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

	db := &SQLKeysDatabase{
		conn:   conn,
		engine: engine,
	}

	return db, nil
}

func (db *SQLKeysDatabase) GetKey(ctx context.Context, did string, label string) (*Key, error) {

	q := "SELECT did, label, private, created, lastmodified FROM keys where did = ? AND label = ?"
	return db.getKey(ctx, q, did, label)
}

func (db *SQLKeysDatabase) getKey(ctx context.Context, q string, args ...interface{}) (*Key, error) {

	row := db.conn.QueryRowContext(ctx, q, args...)

	var did string
	var label string
	var private string
	var created int64
	var lastmod int64

	err := row.Scan(&did, &label, &private, &created, &lastmod)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, atproto.ErrNotFound
		}

		return nil, err
	}

	kp := &Key{
		DID:                 did,
		Label:               label,
		PrivateKeyMultibase: private,
		Created:             created,
		LastModified:        lastmod,
	}

	return kp, err
}

func (db *SQLKeysDatabase) AddKey(ctx context.Context, kp *Key) error {

	q := "INSERT INTO keys (did, label, private, created, lastmodified) VALUES (?, ?, ?, ?, ?)"

	_, err := db.conn.ExecContext(ctx, q, kp.DID, kp.Label, kp.PrivateKeyMultibase, kp.Created, kp.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add key, %w", err)
	}

	return nil
}

func (db *SQLKeysDatabase) DeleteKey(ctx context.Context, kp *Key) error {

	q := "DELETE FROM keys where did = ? AND label = ?"

	_, err := db.conn.ExecContext(ctx, q, kp.DID, kp.Label)

	if err != nil {
		return fmt.Errorf("Failed to delete key, %w", err)
	}

	return nil
}

func (db *SQLKeysDatabase) ListKeys(ctx context.Context, opts *ListKeysOptions) iter.Seq2[*Key, error] {

	return func(yield func(*Key, error) bool) {

		q := "SELECT did, label, private, created, lastmodified FROM keys ORDER BY created DESC"

		rows, err := db.conn.QueryContext(ctx, q)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {

			var did string
			var label string
			var private string
			var created int64
			var lastmod int64

			err := rows.Scan(&did, &label, &private, &created, &lastmod)

			if err != nil {

				if !yield(nil, err) {
					return
				}

				continue
			}

			kp := &Key{
				DID:                 did,
				Label:               label,
				PrivateKeyMultibase: private,
				Created:             created,
				LastModified:        lastmod,
			}

			yield(kp, nil)
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

func (db *SQLKeysDatabase) Close() error {
	return db.conn.Close()
}
