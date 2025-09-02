package delete

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

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
	logger = logger.With("did", opts.DID)

	accounts_db, err := pds.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		logger.Error("Failed to initialize accounts database", "error", err)
		return err
	}

	defer accounts_db.Close()

	keys_db, err := pds.NewKeysDatabase(ctx, opts.KeysDatabaseURI)

	if err != nil {
		logger.Error("Failed to initialize keys database", "error", err)
		return err
	}

	defer keys_db.Close()

	operations_db, err := pds.NewOperationsDatabase(ctx, opts.OperationsDatabaseURI)

	if err != nil {
		logger.Error("Failed to initialize operations database", "error", err)
		return err
	}

	defer operations_db.Close()

	acct, err := accounts_db.GetAccount(ctx, opts.DID)

	if err != nil {
		logger.Error("Failed to retrieve account", "error", err)
		return fmt.Errorf("Failed to retrieve account, %w", err)
	}

	logger = logger.With("handle", acct.Handle)

	k, err := keys_db.GetKey(ctx, acct.DID, "atproto")

	if err != nil {
		logger.Error("Failed to retrieve atproto key", "error", err)
		return fmt.Errorf("Failed to retrieve key, %w", err)
	}

	pr_key, err := k.PrivateKeyK256()

	if err != nil {
		logger.Error("Failed to derive K256 atproto key", "error", err)
		return fmt.Errorf("Failed to derive private key from multibase, %w", err)
	}

	last_op, err := operations_db.GetLastOperationForDID(ctx, acct.DID)

	if err != nil {
		logger.Error("Failed to retrieve last operation", "error", err)
		return fmt.Errorf("Failed to retrieve last operation for DID, %w", err)
	}

	logger = logger.With("last operation", last_op.CID)

	plc_cl := plc.DefaultClient()

	tombstone_op, err := plc.TombstoneDID(ctx, plc_cl, acct.DID, last_op.CID, pr_key)

	if err != nil {
		logger.Error("Failed to tombstone DID", "error", err)
		return fmt.Errorf("Failed to tombstone DID, %w", err)
	}

	tombstone_cid := tombstone_op.CID().String()

	logger = logger.With("tombstone operation", tombstone_cid)
	logger.Info("DID successfully tombstoned")

	op := &pds.Operation{
		CID:       tombstone_cid,
		DID:       acct.DID,
		Operation: tombstone_op,
	}

	err = pds.AddOperation(ctx, operations_db, op)

	if err != nil {
		logger.Error("Failed to record tombstone operation", "error", err)
		return fmt.Errorf("Failed to add operation for tombstone_op, %w", err)
	}

	err = pds.DeleteAccount(ctx, accounts_db, acct)

	if err != nil {
		logger.Error("Failed to delete account database record", "error", err)
		return err
	}

	list_opts := &pds.ListKeysOptions{
		DID: opts.DID,
	}

	for k, err := range keys_db.ListKeys(ctx, list_opts) {

		if err != nil {
			logger.Error("List keys iterator returned an error", "error", err)
			return err
		}

		err = pds.DeleteKey(ctx, keys_db, k)

		if err != nil {
			logger.Error("Failed to remove key database record", "label", k.Label, "error", err)
			return err
		}
	}

	logger.Info("Account successfully deleted.")
	return nil
}
