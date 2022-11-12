package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"log"
	"net"

	"github.com/sanyarise/hezzl/config"
	"github.com/sanyarise/hezzl/internal/cash"
	"github.com/sanyarise/hezzl/internal/db"
	"github.com/sanyarise/hezzl/internal/pb"
	"github.com/sanyarise/hezzl/internal/queue"
	"github.com/sanyarise/hezzl/internal/server"
	"github.com/sanyarise/hezzl/internal/usecases/cashrepo"
	"github.com/sanyarise/hezzl/internal/usecases/qrepo"
	"github.com/sanyarise/hezzl/internal/usecases/userrepo"

	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config := config.NewConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	pgStore, _ := db.NewUserPostgresStore(config.PGUser, config.PGPass, config.PGHost, config.PGPort)
	store := userrepo.NewUserStorage(pgStore)
	userCash, _ := cash.NewRedisClient(config.RedisHost, config.RedisPort, time.Duration(config.CashTTL))
	cash := cashrepo.NewCashStore(userCash)
	userQueue := queue.NewKafkaWriter(config.KafkaTopic, config.KafkaHost, config.KafkaPort)
	queue := qrepo.NewUserQueue(userQueue)
	userServer := server.NewUserServer(store, cash, queue)
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)
	log.Printf("start server on port :%s", config.Port)
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start listener: ", err)
	}
	go serveGRPC(grpcServer, listener)

	<-ctx.Done()

	log.Println("closing...")

	grpcServer.GracefulStop()
	log.Println("grpc server stopped success")
	cancel()
	pgStore.Close()
	log.Println("database stopped success")
	userCash.Close()
	log.Println("cash client stopped success")
	userQueue.Close()
	log.Println("queue client stopped success")

	return nil
}

func serveGRPC(server *grpc.Server, listener net.Listener) {
	err := server.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
