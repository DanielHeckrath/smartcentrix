package main

import (
	"github.com/DanielHeckrath/smartcentrix/proto"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *sensorAPI) RegisterUser(ctx context.Context, req *smartcentrix.RegisterUserRequest) (*smartcentrix.RegisterUserResponse, error) {
	// validate input parameter
	if req.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Username cannot be empty")
	}

	if req.Password == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Password cannot be empty")
	}

	// check if user already exists
	user, err := s.userRepo.GetUserWithName(ctx, req.Name)

	if err != nil && err != errRecordNotFound {
		return nil, grpc.Errorf(codes.Internal, "Unable to check if user exists: %s", err)
	}

	if user != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "A user with this name already exists")
	}

	userID := uuid.NewV4().String()
	token, err := generateToken(userID)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Cannot create new access token for user: %s", err)
	}

	// create new user
	user = &User{}
	user.ID = userID
	user.Name = req.Name
	user.Password = req.Password // TODO do not save password in plaintext!

	// save user in database
	user, err = s.userRepo.SaveUser(ctx, user)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to save new user: %s", err)
	}

	// create registration response
	res := &smartcentrix.RegisterUserResponse{
		User: &smartcentrix.User{
			Id:   user.ID,
			Name: user.Name,
		},
		Token: token,
	}

	return res, nil
}

func (s *sensorAPI) LoginUser(ctx context.Context, req *smartcentrix.RegisterUserRequest) (*smartcentrix.RegisterUserResponse, error) {
	// validate input parameter
	if req.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Username cannot be empty")
	}

	if req.Password == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Password cannot be empty")
	}

	// check if user already exists
	user, err := s.userRepo.GetUserWithCredentials(ctx, req.Name, req.Password)

	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "Unable to login with credentials: Username or password is wrong")
	}

	token, err := generateToken(user.ID)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Cannot create new access token for user: %s", err)
	}

	// create registration response
	res := &smartcentrix.RegisterUserResponse{
		User: &smartcentrix.User{
			Id:   user.ID,
			Name: user.Name,
		},
		Token: token,
	}

	return res, nil
}
