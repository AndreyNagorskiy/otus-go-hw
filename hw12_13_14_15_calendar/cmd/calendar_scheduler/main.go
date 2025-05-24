package main

import (
	"context"
	"flag"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/amqp"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/jackc/pgx/v5/pgxpool"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg := MustLoad(configFile)
	l := logger.NewLogger(cfg.LogLevel, "calendar_scheduler")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.MakeDBConnectionString())
	if err != nil {
		l.Error("Unable to connect to database", slog.String("error", err.Error()))
		return
	}
	defer dbPool.Close()
	storage := sqlstorage.New(dbPool)

	amqpCfg := amqp.Config{
		URL:          cfg.MakeAMQPConnectionString(),
		Exchange:     amqp.NotificationExchange,
		ExchangeType: "fanout",
	}

	mq, err := amqp.New(amqpCfg, l)
	if err != nil {
		l.Error("Unable to connect to amqp server", slog.String("error", err.Error()))
	}
	defer mq.Close()

	publisher, err := mq.NewPublisher(amqp.ExchangeOptions{
		Durable: true,
	})
	if err != nil {
		l.Error("Failed to create publisher", slog.String("error", err.Error()))
	}
	defer publisher.Close()

	sch := scheduler.NewScheduler(cfg.ReminderInterval, cfg.CleanupInterval, l, storage, publisher)
	go sch.Start(ctx)

	l.Info("Scheduler started")

	<-ctx.Done()

	l.Info("Shutting down gracefully...")
}
