package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chdebug"
)

type Log struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func main() {
	out := make(chan []byte, 100)
	go Reciever(out)
	var jsonStr []byte
	for {
		jsonStr = <-out
		l := new(Log)
		err := json.Unmarshal(jsonStr, l)
		if err != nil {
			log.Printf("error on json.Unmarshal: %v", err)
		}
		clickHouseWriter(l)
	}
}

func Reciever(out chan []byte) {
	const (
		topic          = "topic1"
		broker1Address = "localhost:9092"
	)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
		out <- msg.Value
	}
}

func clickHouseWriter(l *Log) {
	fmt.Println("enter in clH")

	db := ch.Connect(
		// clickhouse://<user>:<password>@<host>:<port>/<database>?sslmode=disable
		ch.WithDSN("clickhouse://0.0.0.0:9000/default?sslmode=disable"),
	)

	db.AddQueryHook(chdebug.NewQueryHook(
		chdebug.WithVerbose(true),
		chdebug.FromEnv("CHDEBUG"),
	))

	span := &Log{
		Time:    l.Time,
		Level:   l.Level,
		Message: l.Message,
	}

	res2, err := db.NewInsert().Model(span).Exec(context.Background())
	if err != nil {
		log.Printf("error on insert log in clickhouse: %v", err)
	} else {
		log.Printf("insert into clickhouse success: %v", res2)
	}
}
