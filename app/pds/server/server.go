package server

import (
	"context"
	"flag"
	"log/slog"
	"net/http"

	aa_server "github.com/aaronland/go-http/v3/server"
	"github.com/sfomuseum/go-atproto/http/xrpc/com/atproto/identity"
	"github.com/sfomuseum/go-atproto/http/xrpc/com/atproto/repo"
	"github.com/sfomuseum/go-atproto/pds"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	var users_db pds.UsersDatabase
	var records_db pds.RecordsDatabase

	mux := http.NewServeMux()

	// Resolve handle

	resolve_handle_opts := &identity.ResolveHandleHandlerOptions{
		UsersDatabase: users_db,
	}

	resolve_handle, err := identity.ResolveHandleHandler(resolve_handle_opts)

	if err != nil {
		return err
	}

	mux.Handle(identity.ResolveHandleHandlerURI, resolve_handle)

	// Get record

	get_record_opts := &repo.GetRecordHandlerOptions{
		RecordsDatabase: records_db,
	}

	get_record, err := repo.GetRecordHandler(get_record_opts)

	if err != nil {
		return err
	}

	mux.Handle(repo.GetRecordHandlerURI, get_record)

	s, err := aa_server.NewServer(ctx, opts.ServerURI)

	if err != nil {
		return err
	}

	slog.Info("Listening for requests", "address", s.Address())
	return s.ListenAndServe(ctx, mux)
}
