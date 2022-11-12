package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/segmentio/kafka-go"
	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chdebug"
)

type Config struct {
	KafkaTopic string `env:"KAFKA_TOPIC" envDefault:"topic1"`
	KafkaHost  string `env:"KAFKA_HOST" envDefault:"kafka"`
	KafkaPort  string `env:"KAFKA_PORT" envDefault:"9192"`
	KafkaGroup string `env:"KAFKA_GROUP" envDefault:"my-group"`
	DBDsn      string `env:"DB_DSN" envDefault:"clickhouse://clickhouse:9000/default?sslmode=disable"`
}

type Log struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

var (
	config Config
	once   sync.Once
)

func main() {
	go run()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ctx.Done()

	log.Println("closing...")
	cancel()
}

func run() error {
	config := newConfig()

	broker1Address := fmt.Sprintf("%s:%s", config.KafkaHost, config.KafkaPort)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address},
		Topic:   config.KafkaTopic,
		GroupID: config.KafkaGroup,
	})

	out := make(chan []byte, 100)

	go kafkaReciever(reader, out)

	var jsonStr []byte
	for {
		jsonStr = <-out
		l := new(Log)
		err := json.Unmarshal(jsonStr, l)
		if err != nil {
			log.Printf("error on json.Unmarshal: %v", err)
		}
		clickHouseLogWriter(config.DBDsn, l)
	}
}

// newConfig returns new configuration
func newConfig() *Config {
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

// kafkaReciever reads log messages from kafka
func kafkaReciever(r *kafka.Reader, out chan []byte) {
	defer r.Close()
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("could not read message " + err.Error())
		}
		fmt.Println("received: ", string(msg.Value))
		out <- msg.Value
	}
}

// ckickHouseLogWriter writes log messages to clickhouse
func clickHouseLogWriter(dsn string, l *Log) {
	fmt.Println("enter in clickHouseWriter")

	db := ch.Connect(
		// clickhouse://<user>:<password>@<host>:<port>/<database>?sslmode=disable
		ch.WithDSN(dsn),
	)

	db.AddQueryHook(chdebug.NewQueryHook(
		chdebug.WithVerbose(true),
		chdebug.FromEnv("CHDEBUG"),
	))
	defer db.Close()

	span := &Log{
		Time:    l.Time,
		Level:   l.Level,
		Message: l.Message,
	}

	res, err := db.NewInsert().Model(span).Exec(context.Background())
	if err != nil {
		log.Printf("error on insert log in clickhouse: %v", err)
	} else {
		log.Printf("insert into clickhouse success: %v", res)
	}
}
