package main

import (
	"github.com/DanielHeckrath/smartcentrix/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *sensorAPI) RegisterSensor(ctx context.Context, req *smartcentrix.RegisterSensorRequest) (*smartcentrix.RegisterSensorResponse, error) {
	// validate input parameter
	err := s.validateUserSensor(ctx, req.UserId, req.SensorId, "Cannot register sensor for different user")

	if err != nil {
		return nil, err
	}

	// create new sensor struct
	sensor := &Sensor{}
	sensor.ID = req.SensorId
	sensor.UserID = req.UserId

	if req.Name != nil {
		sensor.Name = req.Name.Value
	}

	if req.RoomId != nil {
		sensor.RoomID = &req.RoomId.Value
	}

	// save new sensor in repo
	sensor, err = s.sensorRepo.SaveSensor(ctx, sensor)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Some internal operation had an error: Unable to save new sensor: %s", err)
	}

	// create new response
	return &smartcentrix.RegisterSensorResponse{}, nil
}

func (s *sensorAPI) UpdateSensor(ctx context.Context, req *smartcentrix.UpdateSensorRequest) (*smartcentrix.UpdateSensorResponse, error) {
	// validate input parameter
	err := s.validateUserSensor(ctx, req.UserId, req.SensorId, "Cannot update sensor for different user")

	if err != nil {
		return nil, err
	}

	sensor, err := s.loadSensor(ctx, req.SensorId)

	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		sensor.Name = req.Name.Value
	}

	if req.RoomId != nil {
		sensor.Name = req.RoomId.Value
	}

	if req.Status != nil {
		sensor.Status = req.Status.Value
	}

	sensor, err = s.sensorRepo.SaveSensor(ctx, sensor)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to save sensor: %s", err)
	}

	res := &smartcentrix.UpdateSensorResponse{
		Sensor: sensor.proto(),
	}

	return res, nil
}

func (s *sensorAPI) ToggleSensor(ctx context.Context, req *smartcentrix.ToggleSensorRequest) (*smartcentrix.ToggleSensorResponse, error) {
	// validate input parameter
	err := s.validateUserSensor(ctx, req.UserId, req.SensorId, "Cannot toggle sensor for different user")

	if err != nil {
		return nil, err
	}

	sensor, err := s.loadSensor(ctx, req.SensorId)

	if err != nil {
		return nil, err
	}

	sensor.Status = req.Status

	sensor, err = s.sensorRepo.SaveSensor(ctx, sensor)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to save sensor: %s", err)
	}

	res := &smartcentrix.ToggleSensorResponse{
		Sensor: sensor.proto(),
	}

	return res, nil
}

func (s *sensorAPI) ListSensor(context.Context, *smartcentrix.ListSensorRequest) (*smartcentrix.ListSensorResponse, error) {
	return nil, errorNotImplemented
}

func (s *sensorAPI) ShowSensor(context.Context, *smartcentrix.ShowSensorRequest) (*smartcentrix.ShowSensorResponse, error) {
	return nil, errorNotImplemented
}

func (s *sensorAPI) DeleteSensor(context.Context, *smartcentrix.DeleteSensorRequest) (*smartcentrix.DeleteSensorResponse, error) {
	return nil, errorNotImplemented
}

func (s *sensorAPI) validateUserSensor(ctx context.Context, userID, sensorID, msg string) error {
	// validate input parameter
	if userID == "" {
		return grpc.Errorf(codes.InvalidArgument, "UserID cannot be empty")
	}

	if sensorID == "" {
		return grpc.Errorf(codes.InvalidArgument, "SensorID cannot be empty")
	}

	// authorize user
	user, err := s.authorizeUser(ctx)

	if err != nil {
		return err
	}

	// check if request sender is target user
	if user.ID != userID {
		return grpc.Errorf(codes.PermissionDenied, "You are not allowed to access this resource: %s", msg)
	}

	return nil
}

func (s *sensorAPI) loadSensor(ctx context.Context, sensorID string) (*Sensor, error) {
	sensor, err := s.sensorRepo.GetSensor(ctx, sensorID)

	if err != nil {
		if err == errRecordNotFound {
			return nil, grpc.Errorf(codes.NotFound, "There is no sensor with this id registered")
		}
		return nil, grpc.Errorf(codes.Internal, "Unable to load sensor: %s", err)
	}

	return sensor, nil
}
