package pds

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"iter"
	"net/url"

	"github.com/did-method-plc/go-didplc"
	"github.com/sfomuseum/go-atproto"
)

type SQLOperationsDatabase struct {
	OperationsDatabase
	conn   *sql.DB
	engine string
}

func init() {

	ctx := context.Background()
	err := RegisterOperationsDatabase(ctx, "sql", NewSQLOperationsDatabase)

	if err != nil {
		panic(err)
	}
}

func NewSQLOperationsDatabase(ctx context.Context, uri string) (OperationsDatabase, error) {

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

	db := &SQLOperationsDatabase{
		conn:   conn,
		engine: engine,
	}

	return db, nil
}

func (db *SQLOperationsDatabase) GetOperation(ctx context.Context, cid string) (*Operation, error) {

	q := "SELECT cid, did, operation, created, lastmodified FROM operations where cid = ?"
	return db.getOperation(ctx, q, cid)
}

func (db *SQLOperationsDatabase) GetLastOperationForDID(ctx context.Context, did string) (*Operation, error) {

	q := "SELECT cid, did, operation, created, lastmodified FROM operations where did = ? ORDER BY created DESC LIMIT 1"
	return db.getOperation(ctx, q, did)
}

func (db *SQLOperationsDatabase) getOperation(ctx context.Context, q string, args ...interface{}) (*Operation, error) {

	row := db.conn.QueryRowContext(ctx, q, args...)

	var cid string
	var did string
	var str_op string
	var created int64
	var lastmod int64

	err := row.Scan(&cid, &did, &str_op, &created, &lastmod)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, atproto.ErrNotFound
		}

		return nil, err
	}

	as_op, err := db.asOperation(str_op)

	if err != nil {
		return nil, err
	}

	op := &Operation{
		CID:          cid,
		DID:          did,
		Operation:    as_op,
		Created:      created,
		LastModified: lastmod,
	}

	return op, err
}

func (db *SQLOperationsDatabase) AddOperation(ctx context.Context, op *Operation) error {

	enc_op, err := json.Marshal(op.Operation)

	if err != nil {
		return fmt.Errorf("Failed to encode DID document, %w", err)
	}

	q := "INSERT INTO operations (cid, did, operation, created, lastmodified) VALUES (?, ?, ?, ?, ?)"

	_, err = db.conn.ExecContext(ctx, q, op.CID, op.DID, string(enc_op), op.Created, op.LastModified)

	if err != nil {
		return fmt.Errorf("Failed to add operation, %w", err)
	}

	return nil
}

func (db *SQLOperationsDatabase) DeleteOperation(ctx context.Context, op *Operation) error {

	q := "DELETE FROM operations where cid = ?"

	_, err := db.conn.ExecContext(ctx, q, op.CID)

	if err != nil {
		return fmt.Errorf("Failed to delete operation, %w", err)
	}

	return nil
}

func (db *SQLOperationsDatabase) ListOperations(ctx context.Context, opts *ListOperationsOptions) iter.Seq2[*Operation, error] {

	return func(yield func(*Operation, error) bool) {

		q := "SELECT cid, did, operation, created, lastmodified FROM operations ORDER BY created DESC"

		rows, err := db.conn.QueryContext(ctx, q)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {

			var cid string
			var did string
			var str_op string
			var created int64
			var lastmod int64

			err := rows.Scan(&cid, &did, &str_op, &created, &lastmod)

			if err != nil {

				if !yield(nil, err) {
					return
				}

				continue
			}

			as_op, err := db.asOperation(str_op)

			if err != nil {
				if !yield(nil, fmt.Errorf("Failed to decode DID document, %w", err)) {
					return
				}

				continue
			}

			op := &Operation{
				CID:          cid,
				DID:          did,
				Operation:    as_op,
				Created:      created,
				LastModified: lastmod,
			}

			yield(op, nil)
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

func (db *SQLOperationsDatabase) Close() error {
	return db.conn.Close()
}

func (db *SQLOperationsDatabase) asOperation(str_op string) (didplc.Operation, error) {

	oe := didplc.OpEnum{}

	err := oe.UnmarshalJSON([]byte(str_op))

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal operation, %w", err)
	}

	as_op := oe.AsOperation()

	if as_op == nil {
		return nil, fmt.Errorf("Failed to derive operation")
	}

	return as_op, nil
}
