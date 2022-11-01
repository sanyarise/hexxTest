package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/sanyarise/hezzl/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer is the server that provides user services
type UserServer struct {
	Store UserStore
	pb.UnimplementedUserServiceServer
}

// NewUserServer returns a new UserServer
func NewUserServer(store UserStore) *UserServer {
	return &UserServer{
		Store: store,
	}
}

func (server *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := req.GetUser()
	log.Printf("receive a create user request with id: %s", user.Id)

	if len(user.Id) > 0 {
		// check if it's a valid UUID
		_, err := uuid.Parse(user.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "user ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new user ID: %v", err)
		}
		user.Id = id.String()
	}
	if ctx.Err() == context.Canceled {
		log.Print("request is canceled")
		return nil, status.Error(codes.Canceled, "deadline is canceled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	// save user to store
	err := server.Store.SaveUser(ctx, user)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save user to the store: %v", err)
	}
	log.Printf("saved user with id: %s", user.Id)

	res := &pb.CreateUserResponse{
		Id: user.Id,
	}
	return res, nil
}

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
		log.Print("request is canceled")
		return nil, status.Error(codes.Canceled, "deadline is canceled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}
	// delete user from store
	err := server.Store.DeleteUser(ctx, id)
	if err != nil {
		code := codes.Internal
		return nil, status.Errorf(code, "cannot delete user from the store: %v", err)
	}
	log.Printf("deleted user with id: %s success", id)

	res := &pb.DeleteUserResponse{
		Status: "Success",
	}
	return res, nil
}

// GetAllUsers is a server-streaming RPC to get all users
func (server *UserServer) GetAllUsers(
	req *pb.AllUsersRequest,
	stream pb.UserService_GetAllUsersServer,
) error {
	res, err := server.Store.GetAllUsers(stream.Context())

	for resUser := range res {
		response := &pb.AllUsersResponse{User: resUser}
		err := stream.Send(response)
		if err != nil {
			return err
		}
		
	}

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}
