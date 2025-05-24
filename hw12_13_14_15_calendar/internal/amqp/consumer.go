package amqp

import (
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch     *amqp.Channel
	cfg    Config
	mq     *Client
	mu     sync.Mutex
	closed bool
}

func (c *Client) NewConsumer(queueOpts QueueOptions, exchangeOpts ExchangeOptions) (*Consumer, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, err := c.newChannel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		c.cfg.Exchange,
		c.cfg.ExchangeType,
		exchangeOpts.Durable,
		exchangeOpts.AutoDelete,
		exchangeOpts.Internal,
		exchangeOpts.NoWait,
		exchangeOpts.Args,
	); err != nil {
		ch.Close()

		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = ch.QueueDeclare(
		c.cfg.Queue,
		queueOpts.Durable,
		queueOpts.AutoDelete,
		queueOpts.Exclusive,
		queueOpts.NoWait,
		queueOpts.Args,
	)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := ch.QueueBind(
		c.cfg.Queue,
		c.cfg.RoutingKey,
		c.cfg.Exchange,
		queueOpts.NoWait,
		queueOpts.Args,
	); err != nil {
		ch.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &Consumer{
		ch:     ch,
		cfg:    c.cfg,
		mq:     c,
		closed: false,
	}, nil
}

func (c *Consumer) Consume(handler func(d amqp.Delivery), opts ConsumeOptions) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return ErrClosed
	}

	msgs, err := c.ch.Consume(
		c.cfg.Queue,
		opts.Consumer,
		opts.AutoAck,
		opts.Exclusive,
		opts.NoLocal,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	return c.ch.Close()
}
