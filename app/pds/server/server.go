package server

import (
	"context"
	"flag"
	"log/slog"
	"net/http"

	aa_server "github.com/aaronland/go-http/v3/server"
	"github.com/sfomuseum/go-atproto/http/xrpc/com/atproto/repo"
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

	mux := http.NewServeMux()

	getrecord_opts := &repo.GetRecordHandlerOptions{}

	getrecord_handler, err := repo.GetRecordHandler(getrecord_opts)

	if err != nil {
		return err
	}

	mux.Handle(repo.GetRecordHandlerURI, getrecord_handler)

	s, err := aa_server.NewServer(ctx, opts.ServerURI)

	if err != nil {
		return err
	}

	slog.Info("Listening for requests", "address", s.Address())
	return s.ListenAndServe(ctx, mux)
}
