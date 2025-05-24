package amqp

import amqp "github.com/rabbitmq/amqp091-go"

type Config struct {
	URL          string
	Exchange     string
	ExchangeType string // direct, fanout, topic, headers
	Queue        string
	RoutingKey   string
}

type ExchangeOptions struct {
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type QueueOptions struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type PublishOptions struct {
	Mandatory  bool
	Immediate  bool
	Persistent bool
}

type ConsumeOptions struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}
