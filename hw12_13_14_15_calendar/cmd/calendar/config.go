package main

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
const (
	MemoryStorageType = "memory"
	SQLStorageType    = "sql"
)

type Config struct {
	LogLevel    string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	StorageType string     `yaml:"storage_type" env:"STORAGE_TYPE" env-default:"memory"`
	DB          Database   `yaml:"db"`
	HTTPServer  HTTPServer `yaml:"http_server" env-prefix:"HTTP_"`
	GRPCServer  GrpcServer `yaml:"grpc_server" env-prefix:"GRPC_"`
}

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	Username string `yaml:"username" env:"DB_USERNAME" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:""`
}

type HTTPServer struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"8080"`
}

type GrpcServer struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"8082"`
}

func MustLoad(cfgFilePath string) Config {
	var cfg Config

	err := cleanenv.ReadConfig(cfgFilePath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	validateStorageType(cfg.StorageType)

	return cfg
}

func validateStorageType(storageType string) {
	switch storageType {
	case MemoryStorageType, SQLStorageType:
		return
	default:
		log.Fatalf("unknown storage type: %s", storageType)
	}
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

func (c *Config) MakeHTTPAddr() string {
	return fmt.Sprintf("%s:%d", c.HTTPServer.Host, c.HTTPServer.Port)
}

func (c *Config) MakeGRPCAddr() string {
	return fmt.Sprintf("%s:%d", c.GRPCServer.Host, c.GRPCServer.Port)
}
