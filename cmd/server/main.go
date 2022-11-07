package main

import (
	"fmt"
	"time"

	"log"
	"net"

	"github.com/sanyarise/hezzl/config"
	"github.com/sanyarise/hezzl/internal/cash"
	"github.com/sanyarise/hezzl/internal/db"
	"github.com/sanyarise/hezzl/internal/logs"
	"github.com/sanyarise/hezzl/internal/pb"
	"github.com/sanyarise/hezzl/internal/server"

	"google.golang.org/grpc"
)

func main() {
	config := config.NewConfig()
	userStore, _ := db.NewUserPostgresStore(config.PGUser, config.PGPass, config.PGHost, config.PGPort)
	userCash, _ := cash.NewRedisClient(config.RedisHost, config.RedisPort, time.Duration(config.CashTTL))
	userQueue := logs.NewKafkaWriter(config.KafkaTopic, config.KafkaHost, config.KafkaPort)
	userServer := server.NewUserServer(userStore, userCash, userQueue)
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)
	log.Printf("start server on port :%s", config.Port)
	address := fmt.Sprintf("0.0.0.0:%s", config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
