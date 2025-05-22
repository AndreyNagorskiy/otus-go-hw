package main

import (
	"context"
	"flag"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg := MustLoad(configFile)
	l := logger.NewLogger(cfg.LogLevel)

	//TODO init RabbitMQ structure

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.MakeDBConnectionString())
	if err != nil {
		l.Error("Unable to connect to database", slog.String("error", err.Error()))
		return
	}
	defer dbPool.Close()
	storage := sqlstorage.New(dbPool)

	sch := scheduler.NewScheduler(cfg.ReminderInterval, cfg.CleanupInterval, l, storage)
	go sch.Start(ctx)

	l.Info("Scheduler started")

	<-ctx.Done()

	l.Info("Shutting down gracefully...")
}
