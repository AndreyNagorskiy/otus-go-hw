package main

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string   `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq" env-prefix:"RABBITMQ_"`
}

type RabbitMQ struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"PORT" env-default:"5672"`
	Username string `yaml:"username" env:"USERNAME" env-default:"guest"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"guest"`
	Vhost    string `yaml:"vhost" env:"VHOST" env-default:"/"`
}

func MustLoad(cfgFilePath string) Config {
	var cfg Config

	err := cleanenv.ReadConfig(cfgFilePath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return cfg
}

func (c *Config) MakeAMQPConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.RabbitMQ.Username,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
		c.RabbitMQ.Vhost)
}
