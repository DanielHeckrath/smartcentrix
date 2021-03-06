package main

import (
	"github.com/juju/errors"
	"golang.org/x/net/context"
)

var (
	errRecordNotFound = errors.New("record not found")
)

// UserRepository is a repository for user data
type UserRepository interface {
	// RegisterUser creates a new user if it does not exists
	RegisterUser(ctx context.Context, user *User) (*User, error)
	// GetUser loads a user with it's id
	GetUser(ctx context.Context, userID string) (*User, error)
	// GetUserWithName loads a user with it's name
	GetUserWithName(ctx context.Context, name string) (*User, error)
	// GetUserWithCredentials loads a user with a credentials combination
	GetUserWithCredentials(ctx context.Context, name, password string) (*User, error)
	// SaveUser updates an existing user
	SaveUser(ctx context.Context, user *User) (*User, error)
	// DeleteUser deletes an existing user
	DeleteUser(ctx context.Context, user *User) error
}

// RoomRepository is a repository for room data
type RoomRepository interface {
	// GetRoom loads a room with it's id
	GetRoom(ctx context.Context, roomID string) (*Room, error)
	// SaveRoom updates or creates a room
	SaveRoom(ctx context.Context, room *Room) (*Room, error)
	// ListRooms loads all rooms for a user
	ListRooms(ctx context.Context, userID string) ([]*Room, error)
	// DeleteRoom deletes an existing room
	DeleteRoom(ctx context.Context, room *Room) error
}

// SensorRepository is a repository for sensor data
type SensorRepository interface {
	// GetSensor loads a sensor with it's id
	GetSensor(ctx context.Context, sensorID string) (*Sensor, error)
	// CreateSensor creates a new sensor
	CreateSensor(ctx context.Context, sensor *Sensor) (*Sensor, error)
	// UpdateSensor updates an existing sensor
	UpdateSensor(ctx context.Context, sensor *Sensor) (*Sensor, error)
	// ListSensors loads all sensors for a user
	ListSensors(ctx context.Context, userID string) ([]*Sensor, error)
	// DeleteSensor deletes an existing sensor
	DeleteSensor(ctx context.Context, sensor *Sensor) error
}

// MeasurementRepository is a repository for sensor measurements
type MeasurementRepository interface {
	// SaveMeasurement updates or creates a measurement
	SaveMeasurement(ctx context.Context, measurement *Measurement) (*Measurement, error)
	// ListMeasurements loads all measurements for a sensor
	ListMeasurements(ctx context.Context, sensorID string) ([]*Measurement, error)
}
