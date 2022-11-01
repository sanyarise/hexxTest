package main

import (
	"flag"
	"fmt"

	"log"
	"net"

	"github.com/sanyarise/hezzl/internal/pb"
	"github.com/sanyarise/hezzl/internal/service"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)
	userStore, err := service.NewUserPostgresStore("postgres://postgres:example@localhost:5432/postgres")
	userServer := service.NewUserServer(userStore)
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
