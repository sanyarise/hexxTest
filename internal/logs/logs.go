package logs

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/sanyarise/hezzl/internal/model"
	"github.com/segmentio/kafka-go"
)

func LogMessageCreate(ctx context.Context, level string, msg string) {
	message := &model.Log{
		Time:    time.Now().Format(time.RFC3339),
		Level:   level,
		Message: msg,
	}
	LogsKafkaProducer(ctx, message)
}

func LogsKafkaProducer(ctx context.Context, msg *model.Log) {
	const (
		topic          = "topic1"
		broker1Address = "localhost:9092"
	)
	wr := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
	})
	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("err on json marshal : %v", err)
	}
	err = wr.WriteMessages(ctx, kafka.Message{
		Value: m,
	})

	if err != nil {
		log.Printf("could not write message: %v", err)
	}

}
