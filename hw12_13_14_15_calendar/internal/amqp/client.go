package amqp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrClosed        = errors.New("connection is closed")
	ErrReconnectFail = errors.New("failed to reconnect")
)

const NotificationExchange = "notifications"

type IPublisher interface {
	Publish(routingKey string, body []byte, opts PublishOptions) error
	PublishWithContext(ctx context.Context, routingKey string, body []byte, opts PublishOptions) error
	Close() error
}

type IConsumer interface {
	Consume(handler func(d amqp.Delivery), opts ConsumeOptions) error
	Close() error
}

type Client struct {
	cfg      Config
	l        logger.Logger
	conn     *amqp.Connection
	mu       sync.RWMutex
	done     chan struct{}
	notifier chan struct{}
	closed   bool
}

func New(cfg Config, l logger.Logger) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	c := &Client{
		cfg:      cfg,
		l:        l,
		conn:     conn,
		done:     make(chan struct{}),
		notifier: make(chan struct{}, 1),
	}

	go c.monitorConnection()
	return c, nil
}

func (c *Client) monitorConnection() {
	closeChan := c.conn.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-c.done:
			return
		case err, ok := <-closeChan:
			if !ok {
				return
			}
			if err != nil {
				c.l.Info("Connection closed", slog.String("error", err.Error()))
			}
			c.reconnect()
		}
	}
}

func (c *Client) reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return ErrClosed
	}

	var conn *amqp.Connection
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		time.Sleep(time.Duration(i+1) * time.Second)

		conn, err = amqp.Dial(c.cfg.URL)
		if err == nil {
			c.conn = conn
			c.l.Info("Successfully reconnected to RabbitMQ")
			select {
			case c.notifier <- struct{}{}:
			default:
			}
			return nil
		}
		c.l.Info(fmt.Sprintf("Reconnect attempt %d failed: %v", i+1, err))
	}

	c.closed = true
	return ErrReconnectFail
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	close(c.done)
	c.closed = true

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

func (c *Client) newChannel() (*amqp.Channel, error) {
	if c.closed {
		return nil, ErrClosed
	}
	return c.conn.Channel()
}
