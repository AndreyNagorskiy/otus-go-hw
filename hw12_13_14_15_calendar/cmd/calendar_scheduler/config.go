package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel         string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	DB               Database      `yaml:"db"`
	RabbitMQ         RabbitMQ      `yaml:"rabbitmq" env-prefix:"RABBITMQ_"`
	ReminderInterval time.Duration `yaml:"reminder_interval" env:"REMINDER_INTERVAL" env-default:"1m"`
	CleanupInterval  time.Duration `yaml:"cleanup_interval" env:"CLEANUP_INTERVAL" env-default:"24h"`
}

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	Username string `yaml:"username" env:"DB_USERNAME" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:""`
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

func (c *Config) MakeDBConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.DB.Username,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}

func (c *Config) MakeAMQPConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.RabbitMQ.Username,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
		c.RabbitMQ.Vhost)
}
