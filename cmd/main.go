package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zchelalo/expense-control-back/pkg/bootstrap"
	"github.com/zchelalo/expense-control-back/pkg/observability"
	"go.uber.org/zap"
)

func main() {
	cfg, err := bootstrap.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	log := bootstrap.GetLogger()
	defer bootstrap.SyncLogger()

	shutdownTracing, err := observability.InitTracing(context.Background(), cfg.ServiceName, cfg.Environment, cfg.OtelExporterOtlpEndpoint)
	if err != nil {
		log.Fatal("cannot init tracing", zap.Error(err))
	}
	defer shutdownTracing(context.Background())

	app, err := bootstrap.InitApp(log, cfg)
	if err != nil {
		log.Fatal("cannot initialize application", zap.Error(err))
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("server starting")
		if err := app.Server.Start(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigs:
		log.Info("signal received, shutting down", zap.String("signal", sig.String()))
	case err := <-errCh:
		log.Error("server error", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = app.Cleanup(ctx)
	log.Info("shutdown complete")
}