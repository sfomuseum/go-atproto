package create

import (
	"context"
	"flag"
	"log/slog"

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

	accounts_db, err := pds.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return err
	}

	defer accounts_db.Close()

	keypairs_db, err := pds.NewKeyPairsDatabase(ctx, opts.KeyPairsDatabaseURI)

	if err != nil {
		return err
	}

	defer keypairs_db.Close()

	operations_db, err := pds.NewOperationsDatabase(ctx, opts.OperationsDatabaseURI)

	if err != nil {
		return err
	}

	defer operations_db.Close()

	logger := slog.Default()

	rsp, err := pds.CreateAccount(ctx, opts.Service, opts.Handle)

	if err != nil {
		return err
	}

	logger = logger.With("did", rsp.Account.DID)
	logger = logger.With("cid", rsp.Operation.CID)
	logger = logger.With("keypair", rsp.KeyPair.Label)

	err = pds.AddAccount(ctx, accounts_db, rsp.Account)

	if err != nil {
		logger.Error("Failed to add account to database", "error", err)
		return err
	}

	err = pds.AddKeyPair(ctx, keypairs_db, rsp.KeyPair)

	if err != nil {
		logger.Error("Failed to add keypair to database", "error", err)
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
