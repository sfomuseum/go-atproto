package create

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-atproto"
	"github.com/sfomuseum/go-atproto/pds"
	"github.com/sfomuseum/go-atproto/plc"
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

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	logger := slog.Default()
	logger = logger.With("handle", opts.Handle)
	logger = logger.With("service", opts.Service)

	accounts_db, err := pds.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return err
	}

	defer accounts_db.Close()

	keys_db, err := pds.NewKeysDatabase(ctx, opts.KeysDatabaseURI)

	if err != nil {
		return err
	}

	defer keys_db.Close()

	operations_db, err := pds.NewOperationsDatabase(ctx, opts.OperationsDatabaseURI)

	if err != nil {
		return err
	}

	defer operations_db.Close()

	acct, err := pds.GetAccountWithHandle(ctx, accounts_db, opts.Handle)

	if acct != nil {
		logger.Error("Handle already exists", "with did", acct.DID)
		return fmt.Errorf("Handle already taken")
	}

	if err != nil && err != atproto.ErrNotFound {
		logger.Error("Failed to determine if handle exists", "error", err)
		return err
	}

	plc_cl := plc.DefaultClient()

	rsp, err := pds.CreateAccount(ctx, plc_cl, opts.Service, opts.Handle)

	if err != nil {
		logger.Error("Failed to create account", "error", err)
		return err
	}

	logger = logger.With("did", rsp.Account.DID)
	logger = logger.With("cid", rsp.Operation.CID)
	logger = logger.With("key", rsp.Key.Label)

	err = pds.AddAccount(ctx, accounts_db, rsp.Account)

	if err != nil {
		logger.Error("Failed to add account to database", "error", err)
		return err
	}

	err = pds.AddKey(ctx, keys_db, rsp.Key)

	if err != nil {
		logger.Error("Failed to add key to database", "error", err)
		return err
	}

	err = pds.AddOperation(ctx, operations_db, rsp.Operation)

	if err != nil {
		logger.Error("Failed to add operation to database", "error", err)
		return err
	}

	logger.Info("New account created")
	return nil
}
