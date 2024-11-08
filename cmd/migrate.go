package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
)

const (
	driverName      = "postgres"
	migrationsTable = "migrations"
	migrationsDir   = "migrations"
)

var migrateCmd = &cli.Command{
	Name:  "migrate",
	Usage: "Manage the database migrations",
	Subcommands: []*cli.Command{
		migrateUpCmd,
	},
}

var migrateUpCmd = &cli.Command{
	Name:  "up",
	Usage: "Run the database migrations",
	Action: func(c *cli.Context) error {
		return migrateUp(c.Context)
	},
}

func migrateUp(_ context.Context) error {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		return errors.New("POSTGRES_URL is required")
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	db, err := sql.Open(driverName, postgresURL)
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: migrationsTable,
	})
	if err != nil {
		return err
	}
	instance, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/%s", wd, migrationsDir),
		driverName,
		driver,
	)
	if err != nil {
		return err
	}

	return instance.Up()
}
