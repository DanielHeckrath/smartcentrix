package main

import (
	"github.com/DanielHeckrath/smartcentrix/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc"
)

var errorNotImplemented = grpc.Errorf(codes.Unimplemented, "This method is not yet implemented")

type sensorAPI struct {
}

func (*sensorAPI) RegisterUser(context.Context, *smartcentrix.RegisterUserRequest) (*smartcentrix.RegisterUserResponse, error) {
	return nil, errorNotImplemented
}

func (*sensorAPI) RegisterSensor(context.Context, *smartcentrix.RegisterSensorRequest) (*smartcentrix.RegisterSensorResponse, error) {
	return nil, errorNotImplemented
}

func (*sensorAPI) ListSensor(context.Context, *smartcentrix.ListSensorRequest) (*smartcentrix.ListSensorResponse, error) {
	return nil, errorNotImplemented
}

func (*sensorAPI) ShowSensor(context.Context, *smartcentrix.ShowSensorRequest) (*smartcentrix.ShowSensorResponse, error) {
	return nil, errorNotImplemented
}

func (*sensorAPI) DeleteSensor(context.Context, *smartcentrix.DeleteSensorRequest) (*smartcentrix.DeleteSensorResponse, error) {
	return nil, errorNotImplemented
}

func (*sensorAPI) ListSensorMeasurement(context.Context, *smartcentrix.ListSensorMeasurementRequest) (*smartcentrix.ListSensorMeasurementResponse, error) {
	return nil, errorNotImplemented
}

func (*sensorAPI) UpdateSensorMeasurement(context.Context, *smartcentrix.UpdateSensorMeasurementRequest) (*smartcentrix.UpdateSensorMeasurementResponse, error) {
	return nil, errorNotImplemented
}

