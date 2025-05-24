package amqp

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch     *amqp.Channel
	cfg    Config
	opts   ExchangeOptions
	mq     *Client
	mu     sync.Mutex
	closed bool
}

func (c *Client) NewPublisher(exchangeOpts ExchangeOptions) (*Publisher, error) {
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

	return &Publisher{
		ch:     ch,
		cfg:    c.cfg,
		opts:   exchangeOpts,
		mq:     c,
		closed: false,
	}, nil
}

func (p *Publisher) Publish(routingKey string, body []byte, opts PublishOptions) error {
	return p.PublishWithContext(context.Background(), routingKey, body, opts)
}

func (p *Publisher) PublishWithContext(ctx context.Context, routingKey string, body []byte, opts PublishOptions) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return ErrClosed
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		deliveryMode := amqp.Transient
		if opts.Persistent {
			deliveryMode = amqp.Persistent
		}

		err := p.ch.Publish(
			p.cfg.Exchange,
			routingKey,
			opts.Mandatory,
			opts.Immediate,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         body,
				DeliveryMode: deliveryMode,
				Timestamp:    time.Now(),
			},
		)
		if err != nil {
			if err := p.reconnect(); err != nil {
				return fmt.Errorf("failed to publish and reconnect: %w", err)
			}
			// Retry once after reconnect
			return p.PublishWithContext(ctx, routingKey, body, opts)
		}

		return nil
	}
}

func (p *Publisher) reconnect() error {
	p.mq.mu.RLock()
	defer p.mq.mu.RUnlock()

	if p.mq.closed {
		return ErrClosed
	}

	ch, err := p.mq.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		p.cfg.Exchange,
		p.cfg.ExchangeType,
		p.opts.Durable,
		p.opts.AutoDelete,
		p.opts.Internal,
		p.opts.NoWait,
		p.opts.Args,
	); err != nil {
		ch.Close()
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	oldCh := p.ch
	p.ch = ch
	oldCh.Close()

	return nil
}

func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	return p.ch.Close()
}
