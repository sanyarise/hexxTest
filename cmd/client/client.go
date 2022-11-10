package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/google/uuid"
	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chdebug"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/sanyarise/hezzl/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Config struct {
	Host  string `env:"SRV_HOST" envDefault:"0.0.0.0"`
	Port  string `env:"SRV_PORT" envDefault:"8080"`
	DBDsn string `env:"DB_DSN" envDefault:"clickhouse://0.0.0.0:9000/default?sslmode=disable"`
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

func createUser(ctx context.Context, userClient pb.UserServiceClient) {
	var quant int
	fmt.Println("enter the desired number of users to create")
	fmt.Scan(&quant)
	if quant <= 0 {
		fmt.Println("need positive integer")
		return
	}
	for i := 0; i < quant; i++ {
		user := NewUser()
		user.Id = ""
		req := &pb.CreateUserRequest{
			User: user,
		}
		res, err := userClient.CreateUser(ctx, req)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.AlreadyExists {
				// not a big deal
				log.Print("user already exists")
			} else {
				log.Fatal("cannot create user: ", err)
			}
			return
		}

		log.Printf("created user with id: %s", res.Id)
	}
}

func deleteUser(ctx context.Context, userClient pb.UserServiceClient) {
	fmt.Println("Enter user id to delete user")
	var id string
	fmt.Scan(&id)
	req := &pb.DeleteUserRequest{
		Id: id,
	}
	res, err := userClient.DeleteUser(ctx, req)
	if err != nil {
		log.Printf("error on delete user: %v", err)
		return
	}
	fmt.Println(res.Status)
}

func getAllUsers(ctx context.Context, userClient pb.UserServiceClient) {
	req := &pb.AllUsersRequest{}

	stream, err := userClient.GetAllUsers(ctx, req)
	if err != nil {
		fmt.Printf("cannot get all users: %v", err)
		os.Exit(1)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Printf("cannot receive response: %v", err)
			os.Exit(1)
		}
		user := res.GetUser()
		log.Print("- found: ", user.GetId())
		log.Print(" + name: ", user.GetName())
	}
}

func clickhouseFunc(dsn string) {
	log.Println("enter in clickhouse func")
	db := ch.Connect(
		// clickhouse://<user>:<password>@<host>:<port>/<database>?sslmode=disable
		ch.WithDSN(dsn),
	)

	db.AddQueryHook(chdebug.NewQueryHook(
		chdebug.WithVerbose(true),
		chdebug.FromEnv("CHDEBUG"),
	))

	type Logs struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Message string `json:"message"`
	}

	logs := &Logs{}

	rows, err := db.QueryContext(context.Background(), `SELECT * FROM logs`)
	if err != nil {
		log.Printf("err on db.QueryContext:%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&logs.Time,
			&logs.Level,
			&logs.Message,
		); err != nil {
			fmt.Print(err)
		}
		fmt.Println(logs)
	}
	if err != nil {
		fmt.Println(err)
	}

}

func NewUser() *pb.User {
	name := randomUserName()

	user := &pb.User{
		Id:   randomID(),
		Name: name,
	}

	return user
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomID() string {
	return uuid.New().String()
}

func randomUserName() string {
	return randomStringFromSet("John Doe", "Mary Key", "Lucy Now", "Tony Stark", "Jim Beam", "Sasha Grey", "Sherlock Holmes", "Dorian Grey", "Bill Gates", "Steve Jobs")
}

func randomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func main() {
	config := newConfig()
	go run(config.Host, config.Port, config.DBDsn)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ctx.Done()

	log.Println("closing...")
	cancel()
}

func run(host string, port string, dsn string) error {
	serverAddress := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	log.Printf("dial server %s", serverAddress)

	userClient := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for {
		fmt.Println("Input type of operation (create, delete, all, ch, exit)")
		var op string
		fmt.Scan(&op)
		s := strings.ToLower(op)
		switch {
		case s == "create":
			createUser(ctx, userClient)
		case s == "delete":
			deleteUser(ctx, userClient)
		case s == "all":
			getAllUsers(ctx, userClient)
		case s == "ch":
			clickhouseFunc(dsn)
		case s == "exit":
			os.Exit(0)
		default:
			fmt.Println("Unknown type of operation")
		}
	}
}
