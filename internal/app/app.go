package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"pvz-service/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

type App struct {
	restServer *http.Server
	grpcServer *grpc.Server
	config     config.Config
	validator  *validator.Validate
	Pool       *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg config.Config) (*App, error) {
	a := &App{
		config: cfg,
	}

	err := a.setup(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Stop(ctx context.Context) error {
	slog.Info("Gracefully shutdown...")
	var errs []error

	if a.grpcServer != nil {
		stopped := make(chan struct{})

		go func() {
			defer close(stopped)
			a.grpcServer.GracefulStop()
		}()

		select {
		case <-stopped:
		case <-ctx.Done():
			a.grpcServer.Stop()
			errs = append(errs, errors.New("grpc server shutdown timeout"))
		}
	}

	if err := a.restServer.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("rest server shutdown: %w", err))
	}

	if a.Pool != nil {
		a.Pool.Close()
	}

	return errors.Join(errs...)
}

func (a *App) RestRun() error {
	listener, err := net.Listen("tcp", a.config.Server.Rest.Address)
	if err != nil {
		return err
	}
	slog.Info("rest starting on:" + a.config.Server.Rest.Address)
	return a.restServer.Serve(listener)
}

func (a *App) GrpcRun() error {
	listener, err := net.Listen("tcp", a.config.Server.GRPC.Address)
	if err != nil {
		return err
	}
	slog.Info("grpc starting on:" + a.config.Server.GRPC.Address)
	return a.grpcServer.Serve(listener)
}
