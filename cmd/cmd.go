package cmd

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "app",
	Usage: "Run the app",
	Commands: []*cli.Command{
		serverCmd,
		migrateCmd,
	},
}

func Execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return app.RunContext(ctx, os.Args)
}
