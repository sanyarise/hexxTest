package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/sanyarise/hezzl/internal/model"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	Writer *kafka.Writer
}

func NewKafkaWriter(topic string, host string, port string) *KafkaWriter {
	brokerAddress := fmt.Sprintf("%s:%s", host, port)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})
	return &KafkaWriter{Writer: writer}
}

func (kw *KafkaWriter) LogsKafkaProduce(ctx context.Context, level string, msg string) error {
	message := &model.Log{
		Time:    time.Now().Format(time.RFC3339),
		Level:   level,
		Message: msg,
	}
	m, err := json.Marshal(message)
	if err != nil {
		log.Printf("err on json marshal : %v", err)
	}
	err = kw.Writer.WriteMessages(ctx, kafka.Message{
		Value: m,
	})

	if err != nil {
		log.Printf("could not write message: %v", err)
		return fmt.Errorf("could not write message: %w", err)
	}
	return nil
}
