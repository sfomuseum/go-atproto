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
	logger = logger.With("did", opts.DID)

	rsp, err := pds.DeleteAccount(ctx, opts.DID)

	if err != nil {
		return err
	}

	logger = logger.With("cid", rsp.Operation.CID)

	err = pds.RemoveAccount(ctx, accounts_db, rsp.Account)

	if err != nil {
		logger.Error("Failed to add account to database", "error", err)
		return err
	}

	list_opts := &pds.ListKeyPairsOptions{
		DID: opts.DID,
	}

	for kp, err := range keypairs_db.ListKeyPairs(ctx, list_opts) {

		if err != nil {
			logger.Error("List keypairs iterator returned an error", "error", err)
			return err
		}

		err = pds.DeleteKeyPair(ctx, keypairs_db, kp)

		if err != nil {
			logger.Error("Failed to remove keypair to database", "label", kp.Label, "error", err)
			return err
		}
	}

	err = pds.AddOperation(ctx, operations_db, rsp.Operation)

	if err != nil {
		logger.Error("Failed to add operation to database", "error", err)
		return err
	}

	logger.Info("New account created")
	return nil
}
