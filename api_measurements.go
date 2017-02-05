package main

import (
	"fmt"
	"log"

	"github.com/DanielHeckrath/smartcentrix-notifications/notification"
	"github.com/DanielHeckrath/smartcentrix/proto"
	"github.com/golang/protobuf/ptypes/wrappers"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const measurementAlertThreshold int64 = 5 * 1000 // 5kWh

func (s *sensorAPI) ListSensorMeasurement(context.Context, *smartcentrix.ListSensorMeasurementRequest) (*smartcentrix.ListSensorMeasurementResponse, error) {
	return nil, errorNotImplemented
}

func (s *sensorAPI) UpdateSensorMeasurement(ctx context.Context, req *smartcentrix.UpdateSensorMeasurementRequest) (*smartcentrix.UpdateSensorMeasurementResponse, error) {
	// validate input parameter
	if req.SensorId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "SensorID cannot be empty")
	}

	if len(req.Measurements) < 1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Measurements cannot be empty")
	}

	// make sure calling and target user are the same
	err := s.validateUserOwnership(ctx, req.UserId, "Cannot update measurements for different users sensor")

	sensor, err := s.loadSensor(ctx, req.SensorId)

	if err != nil {
		return nil, err
	}

	measurements := make([]*Measurement, 0, len(req.Measurements))

	for _, m := range req.Measurements {
		measurement := &Measurement{}
		measurement.ID = m.Id
		measurement.Timestamp = m.Timestamp
		measurement.Value = m.Value
		measurement.SensorID = req.SensorId

		measurement, err = s.measurementRepo.SaveMeasurement(ctx, measurement)

		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "Unable to save measurement: %s", err)
		}

		measurements = append(measurements, measurement)
	}

	// check if we need to send alert notification in background
	go s.calculateSensorAlert(context.Background(), sensor, measurements)

	return &smartcentrix.UpdateSensorMeasurementResponse{}, nil
}

func (s *sensorAPI) calculateSensorAlert(ctx context.Context, sensor *Sensor, measurements []*Measurement) {
	var sendAlert bool
	var m *Measurement

	for _, measurement := range measurements {
		if measurement.Value > measurementAlertThreshold {
			sendAlert = true
			m = measurement
			break
		}
	}

	if !sendAlert {
		return
	}

	n := notification.Notification{
		Title: &wrappers.StringValue{Value: "One of your sensors is over it's threshold"},
		Body:  &wrappers.StringValue{Value: fmt.Sprintf("Sensor %s is over it's threshold with a value of %d", sensor.Name, m.Value)},
	}

	client, err := getNotificationServiceClient()

	if err != nil {
		log.Printf("Unable to create notification service client: %s", err)
	}

	req := &notification.SendNotificationRequest{
		Notification: &n,
		// TODO: load device push tokens
	}

	_, err = client.SendNotification(context.Background(), req)

	if err != nil {
		log.Printf("Unable to send notification: %s", err)
	}
}
