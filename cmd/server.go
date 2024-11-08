package cmd

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nkbhasker/go-pgx-transaction-example/internal/api"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/db"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/model"
	"github.com/nkbhasker/go-pgx-transaction-example/internal/service"
	"github.com/urfave/cli/v2"
)

var serverCmd = &cli.Command{
	Name:  "server",
	Usage: "Manage the server",
	Subcommands: []*cli.Command{
		startServerCmd,
	},
}

var startServerCmd = &cli.Command{
	Name:  "start",
	Usage: "Start the server",
	Action: func(c *cli.Context) error {
		return startServer(c.Context)
	},
}

func startServer(ctx context.Context) error {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		return errors.New("POSTGRES_URL is required")
	}
	db, err := db.NewDB(postgresURL)
	if err != nil {
		return err
	}
	defer db.Close()
	repo := model.NewRepository(model.NewBaseRepository(db))
	service := service.NewService(service.NewBaseService(repo))
	api := api.NewAPI(service)

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	server := &http.Server{
		Addr:    ":8090",
		Handler: api.Handler(),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		log.Println("Server stopped")
		errCh <- nil
	}()

	go func() {
		<-sigCh
		log.Println("Shutting down server")
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		go func() {
			<-ctx.Done()
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				errCh <- errors.New("graceful shutdown timed out... forcing exit")
			}
			errCh <- ctx.Err()
		}()
		errCh <- server.Shutdown(ctx)
	}()
	log.Println("Server started")

	return <-errCh
}
