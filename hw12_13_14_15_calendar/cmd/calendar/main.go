package main

import (
	"context"
	"flag"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/handlers/grpc"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/jackc/pgx/v5/pgxpool"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg := MustLoad(configFile)
	l := logger.NewLogger(cfg.LogLevel, "calendar")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage app.Storage

	switch cfg.StorageType {
	case MemoryStorageType:
		storage = memorystorage.NewStorage()
	case SQLStorageType:
		dbConnectionString := cfg.MakeDBConnectionString()
		err := sqlstorage.Migrate(dbConnectionString, false)
		if err != nil {
			l.Error("Unable to migrate database", slog.String("error", err.Error()))
			return
		}
		dbPool, err := pgxpool.New(ctx, dbConnectionString)
		if err != nil {
			l.Error("Unable to connect to database", slog.String("error", err.Error()))
			return
		}
		defer dbPool.Close()
		storage = sqlstorage.New(dbPool)
	default:
		l.Error("Unsupported storage type", slog.String("storage_type", cfg.StorageType))
	}

	calendar := app.New(l, storage)
	httpServer := internalhttp.NewServer(l, calendar, cfg.MakeHTTPAddr())

	go func() {
		if err := httpServer.Start(); err != nil {
			l.Error("HttpServer failed", slog.String("error", err.Error()))
			cancel()
		}
	}()

	grpcHandler := grpchandler.NewEventHandler(calendar)
	grpcServer := internalgrpc.NewServer(l, grpcHandler)

	go func() {
		if err := grpcServer.Run(cfg.MakeGRPCAddr()); err != nil {
			l.Error("Grpc server failed", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		l.Error("Shutdown failed", slog.String("error", err.Error()))
		panic("server shutdown failed")
	}

	grpcServer.Stop()

	l.Info("Application stopped")
}
