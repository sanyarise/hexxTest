package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/sanyarise/hezzl/internal/model"
	"github.com/sanyarise/hezzl/internal/usecases/qrepo"
	"github.com/segmentio/kafka-go"
)

var _ qrepo.Queue = &KafkaWriter{}

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

func (kw *KafkaWriter) Close() {
	kw.Writer.Close()
}

func (kw *KafkaWriter) Enqueue(ctx context.Context, level string, msg string) error {
	var once sync.Once
	once.Do(func() {
		KafkaCreateTopic()
	})
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

func KafkaCreateTopic() {
	topic := "topic1"

	conn, err := kafka.Dial("tcp", "kafka:9192")
	if err != nil {
		log.Printf("error on kafka.Dial: %v\n", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Printf("error on conn.Controller: %v\n", err)
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Printf("error on controllerConn kafka.Dial: %v\n", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Printf("error on controllerConn.CreateTopics: %v\n", err)
	}
}
