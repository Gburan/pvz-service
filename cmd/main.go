package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"pvz-service/internal/app"
	"pvz-service/internal/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	ctx := context.TODO()

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		slog.Error("cannot setup app:", "error", err)
		os.Exit(1)
	}
	go func() {
		if err = a.RestRun(); err != nil {
			slog.Error("stop rest server:", "error", err)
			os.Exit(1)
		}
	}()
	go func() {
		if err = a.GrpcRun(); err != nil {
			slog.Error("stop grpc server:", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("got interruption signal")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = a.Stop(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
}
