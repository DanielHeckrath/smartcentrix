package main

import (
	"github.com/DanielHeckrath/smartcentrix/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *sensorAPI) RegisterRoom(ctx context.Context, req *smartcentrix.RegisterRoomRequest) (*smartcentrix.RegisterRoomResponse, error) {
	// validate input parameter
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot register room for different user")

	if err != nil {
		return nil, err
	}

	// create new room struct
	room := &Room{}
	room.UserID = req.UserId
	room.Name = req.Name

	// save new room in repo
	room, err = s.roomRepo.SaveRoom(ctx, room)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Some internal operation had an error: Unable to save new room: %s", err)
	}

	// create new response
	return &smartcentrix.RegisterRoomResponse{
		Room: room.proto(),
	}, nil
}

func (s *sensorAPI) UpdateRoom(ctx context.Context, req *smartcentrix.UpdateRoomRequest) (*smartcentrix.UpdateRoomResponse, error) {
	// validate input parameter
	if req.RoomId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "RoomID cannot be empty")
	}

	// make sure calling and target user are the same
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot update room for different user")

	if err != nil {
		return nil, err
	}

	room, err := s.loadRoom(ctx, req.RoomId)

	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		room.Name = req.Name.Value
	}

	room, err = s.roomRepo.SaveRoom(ctx, room)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to save room: %s", err)
	}

	res := &smartcentrix.UpdateRoomResponse{
		Room: room.proto(),
	}

	return res, nil
}

func (s *sensorAPI) DeleteRoom(ctx context.Context, req *smartcentrix.DeleteRoomRequest) (*smartcentrix.DeleteRoomResponse, error) {
	// validate input parameter
	if req.RoomId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "RoomID cannot be empty")
	}

	// make sure calling and target user are the same
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot toggle room for different user")

	if err != nil {
		return nil, err
	}

	room, err := s.loadRoom(ctx, req.RoomId)

	if err != nil {
		return nil, err
	}

	err = s.roomRepo.DeleteRoom(ctx, room)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to delete room: %s", err)
	}

	res := &smartcentrix.DeleteRoomResponse{
		Room: room.proto(),
	}

	return res, nil
}

func (s *sensorAPI) loadRoom(ctx context.Context, roomID string) (*Room, error) {
	room, err := s.roomRepo.GetRoom(ctx, roomID)

	if err != nil {
		if err == errRecordNotFound {
			return nil, grpc.Errorf(codes.NotFound, "There is no room with this id registered")
		}
		return nil, grpc.Errorf(codes.Internal, "Unable to load room: %s", err)
	}

	return room, nil
}
