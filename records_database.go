package pds

import (
	"context"
)

type RecordsDatabase interface {
	GetRecord(context.Context, string, string, string) (*Record, error)
	AddRecord(context.Context, *Record) error
	UpdateRecord(context.Context, *Record) error
	DeleteRecord(context.Context, *Record) error
}
