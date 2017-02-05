package main

import (
	"github.com/DanielHeckrath/smartcentrix/proto"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *sensorAPI) RegisterSensor(ctx context.Context, req *smartcentrix.RegisterSensorRequest) (*smartcentrix.RegisterSensorResponse, error) {
	// validate input parameter
	if req.SensorId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "SensorID cannot be empty")
	}

	// make sure calling and target user are the same
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot register sensor for different user")

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
	sensor, err = s.sensorRepo.CreateSensor(ctx, sensor)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Some internal operation had an error: Unable to save new sensor: %s", err)
	}

	// create new response
	return &smartcentrix.RegisterSensorResponse{
		Sensor: sensor.proto(),
	}, nil
}

func (s *sensorAPI) UpdateSensor(ctx context.Context, req *smartcentrix.UpdateSensorRequest) (*smartcentrix.UpdateSensorResponse, error) {
	// validate input parameter
	if req.SensorId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "SensorID cannot be empty")
	}

	// make sure calling and target user are the same
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot update sensor for different user")

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
		sensor.RoomID = &req.RoomId.Value
	}

	if req.Status != nil {
		sensor.Status = req.Status.Value
	}

	if req.InUse != nil {
		sensor.InUse = req.InUse.Value
	}

	sensor, err = s.sensorRepo.UpdateSensor(ctx, sensor)

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
	if req.SensorId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "SensorID cannot be empty")
	}

	// make sure calling and target user are the same
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot toggle sensor for different user")

	if err != nil {
		return nil, err
	}

	sensor, err := s.loadSensor(ctx, req.SensorId)

	if err != nil {
		return nil, err
	}

	sensor.Status = req.Status

	sensor, err = s.sensorRepo.UpdateSensor(ctx, sensor)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to save sensor: %s", err)
	}

	res := &smartcentrix.ToggleSensorResponse{
		Sensor: sensor.proto(),
	}

	return res, nil
}

func (s *sensorAPI) ListSensor(ctx context.Context, req *smartcentrix.ListSensorRequest) (*smartcentrix.ListSensorResponse, error) {
	// make sure calling and target user are the same
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot view sensor for different user")

	if err != nil {
		return nil, err
	}

	var sensors []*Sensor
	var rooms []*Room
	var g *errgroup.Group
	g, ctx = errgroup.WithContext(ctx)

	g.Go(func() error {
		sens, err := s.sensorRepo.ListSensors(ctx, req.UserId)

		if err != nil {
			return grpc.Errorf(codes.Internal, "Unable to load sensors: %s", err)
		}

		sensors = sens
		return nil
	})

	g.Go(func() error {
		r, err := s.roomRepo.ListRooms(ctx, req.UserId)

		if err != nil {
			return grpc.Errorf(codes.Internal, "Unable to load rooms: %s", err)
		}

		rooms = r
		return nil
	})

	if err = g.Wait(); err != nil {
		return nil, err
	}

	res := &smartcentrix.ListSensorResponse{
		Sensors: convertSensors(sensors),
		Rooms:   convertRooms(rooms),
	}

	return res, nil
}

func (s *sensorAPI) ShowSensor(ctx context.Context, req *smartcentrix.ShowSensorRequest) (*smartcentrix.ShowSensorResponse, error) {
	// validate input parameter
	if req.SensorId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "SensorID cannot be empty")
	}

	// make sure calling and target user are the same
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot view sensor for different user")

	if err != nil {
		return nil, err
	}

	sensor, err := s.loadSensor(ctx, req.SensorId)

	if err != nil {
		return nil, err
	}

	res := &smartcentrix.ShowSensorResponse{
		Sensor:       sensor.proto(),
		Measurements: convertMeasurements(sensor.Measurements),
	}

	return res, nil
}

func (s *sensorAPI) DeleteSensor(ctx context.Context, req *smartcentrix.DeleteSensorRequest) (*smartcentrix.DeleteSensorResponse, error) {
	// validate input parameter
	if req.SensorId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "SensorID cannot be empty")
	}

	// make sure calling and target user are the same
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot delete sensor for different user")

	if err != nil {
		return nil, err
	}

	sensor, err := s.loadSensor(ctx, req.SensorId)

	if err != nil {
		return nil, err
	}

	err = s.sensorRepo.DeleteSensor(ctx, sensor)

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to delete sensor: %s", err)
	}

	res := &smartcentrix.DeleteSensorResponse{
		Sensor: sensor.proto(),
	}

	return res, nil
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
