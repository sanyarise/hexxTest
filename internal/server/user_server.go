package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sanyarise/hezzl/internal/cash"
	"github.com/sanyarise/hezzl/internal/db"
	"github.com/sanyarise/hezzl/internal/logs"
	"github.com/sanyarise/hezzl/internal/pb"
	"github.com/sanyarise/hezzl/internal/usecases/userrepo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer is the server that provides user services
type UserServer struct {
	Store userrepo.UserStore
	Cash  cash.CashStore
	Queue *logs.KafkaWriter
	pb.UnimplementedUserServiceServer
}

// NewUserServer returns a new UserServer
func NewUserServer(store userrepo.UserStore, cash cash.CashStore, queue *logs.KafkaWriter) *UserServer {
	return &UserServer{
		Store: store,
		Cash:  cash,
		Queue: queue,
	}
}

// CreateUser create new user and save in store
func (server *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := req.GetUser()
	log.Printf("receive a create user request with id: %s", user.Id)

	if len(user.Id) > 0 {
		// check if it's a valid UUID
		_, err := uuid.Parse(user.Id)
		if err != nil {
			msg := status.Errorf(codes.InvalidArgument, "user ID is not a valid UUID: %v", err)
			server.Queue.LogsKafkaProduce(ctx, "ERROR", msg.Error())
			return nil, msg
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			msg := status.Errorf(codes.Internal, "cannot generate a new user ID: %v", err)
			server.Queue.LogsKafkaProduce(ctx, "ERROR", msg.Error())
			return nil, msg
		}
		user.Id = id.String()
	}
	if ctx.Err() == context.Canceled {
		log.Println("request is canceled")
		msg := status.Error(codes.Canceled, "request is canceled")
		server.Queue.LogsKafkaProduce(ctx, "ERROR", msg.Error())
		return nil, msg
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Println("deadline is exceeded")
		msg := status.Error(codes.DeadlineExceeded, "deadline is exceeded")
		server.Queue.LogsKafkaProduce(ctx, "ERROR", msg.Error())
		return nil, msg
	}
	// save user to store
	err := server.Store.SaveUser(ctx, user)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, db.ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		msg := status.Errorf(code, "cannot save user to the store: %v", err)
		server.Queue.LogsKafkaProduce(ctx, "ERROR", msg.Error())
		return nil, msg
	}
	msg := fmt.Sprintf("user with id %v created successfully", user.Id)
	log.Println(msg)
	server.Queue.LogsKafkaProduce(ctx, "INFO", msg)

	res := &pb.CreateUserResponse{
		Id: user.Id,
	}
	return res, nil
}

// DeleteUser delete user by id
func (server *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id := req.Id
	log.Printf("receive a delete user request with id: %s", id)

	if len(id) > 0 {
		// check if it's a valid UUID
		_, err := uuid.Parse(id)

		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "user ID is not a valid UUID: %v", err)
		}
	}

	if ctx.Err() == context.Canceled {
		log.Println("request is canceled")
		return nil, status.Error(codes.Canceled, "context is canceled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Println("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	// delete user from store
	err := server.Store.DeleteUser(ctx, id)
	if err != nil {
		code := codes.Internal
		return nil, status.Errorf(code, "cannot delete user from the store: %v", err)
	}
	log.Printf("deleted user with id: %s success\n", id)

	res := &pb.DeleteUserResponse{
		Status: fmt.Sprintf("Delete user with id %s success\n", id),
	}
	return res, nil
}

// GetAllUsers is a server-streaming RPC to get all users
func (server *UserServer) GetAllUsers(req *pb.AllUsersRequest, stream pb.UserService_GetAllUsersServer) error {
	req.Request = "Get all users"
	if ok := server.Cash.CheckCash(req.Request); !ok {
		log.Println("CheckCash() returns false")
		res, err := server.Store.GetAllUsers(stream.Context())
		if err != nil {
			log.Printf("error on server.Store.GetAllUsers() %v", err)
			return err
		}
		ctx := context.Background()
		err = server.Cash.CreateCash(ctx, res, req.Request)
		if err != nil {
			log.Printf("error on CreateCash: %v", err)
			return err
		}
		log.Println("Create cash success")
	}
	res, err := server.Cash.GetCash(req.Request)
	if err != nil {
		return status.Errorf(codes.Internal, "error on GetCash(): %v", err)
	}
	log.Println("Get cash success")
	for _, v := range res {
		err := stream.Send(&pb.AllUsersResponse{User: v})
		if err != nil {
			return nil
		}
	}
	return nil
}
