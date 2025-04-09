package main

import (
	"context"
	"flag"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	l := logger.NewLogger(cfg.LogLevel)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage app.Storage

	switch cfg.StorageType {
	case MemoryStorageType:
		storage = memorystorage.New()
	case SQLStorageType:
		dbConnectionString := cfg.MakeDbConnectionString()
		sqlstorage.Migrate(dbConnectionString, false)
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
	server := internalhttp.NewServer(l, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			l.Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := server.Start(ctx); err != nil {
		l.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
