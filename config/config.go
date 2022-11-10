package config

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/caarlos0/env/v6"
)

// Config describe configuration
type Config struct {
	Port       string `env:"SRV_PORT" envDefault:"8080"`
	Host       string `env:"SRV_HOST" envDefault:"0.0.0.0"`
	PGUser     string `env:"PG_USER" envDefault:"postgres"`
	PGPass     string `env:"PG_PASSWORD" envDefault:"example"`
	PGPort     string `env:"PG_PORT" envDefault:"5432"`
	PGHost     string `env:"PG_HOST" envDefault:"postgres"`
	CashTTL    int    `env:"SRV_CASH_TTL" envDefault:"1"`
	RedisPort  string `env:"REDIS_PORT" envDefault:"6379"`
	RedisHost  string `env:"REDIS_HOST" envDefault:"redis"`
	KafkaPort  string `env:"KAFKA_PORT" envDefault:"9192"`
	KafkaHost  string `env:"KAFKA_HOST" envDefault:"kafka"`
	KafkaTopic string `env:"KAFKA_TOPIC" envDefault:"topic1"`
}

var (
	config Config
	once   sync.Once
)

// NewConfig create service config
func NewConfig() *Config {
	once.Do(func() {
		if err := env.Parse(&config); err != nil {
			log.Fatalf("Can't load configuration: %s", err)
		}
		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Load config successful %v", string(configBytes))
	})
	return &config
}
