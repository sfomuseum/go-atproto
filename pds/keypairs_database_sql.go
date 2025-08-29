package pds

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"net/url"

	"github.com/sfomuseum/go-atproto"
)

type SQLKeyPairsDatabase struct {
	KeyPairsDatabase
	conn   *sql.DB
	engine string
}

func init() {

	ctx := context.Background()
	err := RegisterKeyPairsDatabase(ctx, "sql", NewSQLKeyPairsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLKeyPairsDatabase(ctx context.Context, uri string) (KeyPairsDatabase, error) {

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

	db := &SQLKeyPairsDatabase{
		conn:   conn,
		engine: engine,
	}

	return db, nil
}

func (db *SQLKeyPairsDatabase) GetKeyPair(ctx context.Context, did string, label string) (*KeyPair, error) {

	q := "SELECT did, label, public, private, created, lastmodified FROM keypairs where did = ? AND label = ?"
	return db.getKeyPair(ctx, q, did, label)
}

func (db *SQLKeyPairsDatabase) getKeyPair(ctx context.Context, q string, args ...interface{}) (*KeyPair, error) {

	row := db.conn.QueryRowContext(ctx, q, args...)

	var did string
	var label string
	var public string
	var private string
	var created int64
	var lastmod int64

	err := row.Scan(&did, &label, &public, &private, &created, &lastmod)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, atproto.ErrNotFound
		}

		return nil, err
	}

	kp := &KeyPair{
		DID:                 did,
		Label:               label,
		PublicKeyMultibase:  public,
		PrivateKeyMultibase: private,
		Created:             created,
		LastModified:        lastmod,
	}

	return kp, err
}

func (db *SQLKeyPairsDatabase) AddKeyPair(ctx context.Context, kp *KeyPair) error {

	q := "INSERT INTO keypairs (did, label, public, private, created, lastmodified) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := db.conn.ExecContext(ctx, q, kp.DID, kp.Label, kp.PublicKeyMultibase, kp.PrivateKeyMultibase, kp.Created, kp.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add keypair, %w", err)
	}

	return nil
}

func (db *SQLKeyPairsDatabase) DeleteKeyPair(ctx context.Context, kp *KeyPair) error {

	q := "DELETE FROM keypairs where did = ? AND label = ?"

	_, err := db.conn.ExecContext(ctx, q, kp.DID, kp.Label)

	if err != nil {
		return fmt.Errorf("Failed to delete keypair, %w", err)
	}

	return nil
}

func (db *SQLKeyPairsDatabase) ListKeyPairs(ctx context.Context, opts *ListKeyPairsOptions) iter.Seq2[*KeyPair, error] {

	return func(yield func(*KeyPair, error) bool) {

		q := "SELECT did, label, public, private, created, lastmodified FROM keypairs ORDER BY created DESC"

		rows, err := db.conn.QueryContext(ctx, q)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {

			var did string
			var label string
			var public string
			var private string
			var created int64
			var lastmod int64

			err := rows.Scan(&did, &label, &public, &private, &created, &lastmod)

			if err != nil {

				if !yield(nil, err) {
					return
				}

				continue
			}

			kp := &KeyPair{
				DID:                 did,
				Label:               label,
				PublicKeyMultibase:  public,
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

func (db *SQLKeyPairsDatabase) Close() error {
	return db.conn.Close()
}
