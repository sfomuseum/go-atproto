package create

import (
	"context"
	"flag"
	_ "log/slog"

	"github.com/sfomuseum/go-pds"
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

	users_db, err := pds.NewUsersDatabase(ctx, opts.UsersDatabaseURI)

	if err != nil {
		return err
	}

	defer users_db.Close()

	u, err := pds.CreateUser(ctx)

	if err != nil {
		return err
	}

	err = pds.AddUser(ctx, users_db, u)

	if err != nil {
		return err
	}

	return nil
}
