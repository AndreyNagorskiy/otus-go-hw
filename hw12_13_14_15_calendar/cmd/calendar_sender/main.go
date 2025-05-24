package main

import (
	"context"
	"flag"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/amqp"
	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rabbitmq/amqp091-go"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg := MustLoad(configFile)
	l := logger.NewLogger(cfg.LogLevel, "calendar_sender")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	amqpCfg := amqp.Config{
		URL:          cfg.MakeAMQPConnectionString(),
		Exchange:     amqp.NotificationExchange,
		ExchangeType: "fanout",
		Queue:        "logs-processor",
		RoutingKey:   "",
	}

	mq, err := amqp.New(amqpCfg, l)
	if err != nil {
		l.Error("Unable to connect to amqp server", slog.String("error", err.Error()))
	}
	defer func() {
		if err := mq.Close(); err != nil {
			l.Error("Failed to close AMQP connection", slog.String("error", err.Error()))
		}
	}()

	consumer, err := mq.NewConsumer(
		amqp.QueueOptions{Durable: true},
		amqp.ExchangeOptions{Durable: true},
	)
	if err != nil {
		l.Error("Failed to create consumer", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			l.Error("Failed to close consumer", slog.String("error", err.Error()))
		}
	}()

	done := make(chan struct{})
	consumeErr := make(chan error, 1)

	go func() {
		err := consumer.Consume(func(d amqp091.Delivery) {
			select {
			case <-ctx.Done():
				// При получении сигнала завершения возвращаем сообщение в очередь
				if err := d.Nack(false, true); err != nil {
					l.Error("Failed to Nack message", slog.String("error", err.Error()))
				}
				return
			default:
				l.Info("Processing message", slog.String("message", string(d.Body)))
				if err := d.Ack(false); err != nil {
					l.Error("Failed to Ack message", slog.String("error", err.Error()))
				}
			}
		}, amqp.ConsumeOptions{
			Consumer: "notifications-handler",
			AutoAck:  false,
		})
		if err != nil {
			consumeErr <- err
		}
		close(done)
	}()

	l.Info("Sender started successfully")

	// Ожидание событий завершения
	select {
	case <-ctx.Done():
		l.Info("Shutdown signal received, initiating graceful shutdown...")

		// Даем фиксированное время на завершение
		select {
		case <-done:
			l.Info("Consumer finished processing messages")
		case err := <-consumeErr:
			l.Info("Consumer stopped with error", slog.String("error", err.Error()))
		case <-time.After(5 * time.Second):
			l.Info("Shutdown timeout reached, forcing exit")
		}

	case err := <-consumeErr:
		l.Error("Consumer failed", slog.String("error", err.Error()))
		cancel()

	case <-done:
		l.Info("Consumer finished normally")
	}

	l.Info("Service shutdown completed")
}
