package main

import (
	"flag"
	"fmt"
	"time"

	"log"
	"net"

	"github.com/sanyarise/hezzl/internal/cash"
	"github.com/sanyarise/hezzl/internal/db"
	"github.com/sanyarise/hezzl/internal/pb"
	"github.com/sanyarise/hezzl/internal/server"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)
	userStore, _ := db.NewUserPostgresStore("postgres://postgres:example@localhost:5432/postgres")
	userCash, _ := cash.NewRedisClient("localhost", "6379", 1*time.Hour)
	userServer := server.NewUserServer(userStore, userCash)
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
