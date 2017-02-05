package main

import (
	"github.com/juju/errors"
	"golang.org/x/net/context"
)

var (
	errRecordNotFound = errors.New("record not found")
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, userID string) (*User, error)
	GetUserWithName(ctx context.Context, name string) (*User, error)
	GetUserWithCredentials(ctx context.Context, name, password string) (*User, error)
	SaveUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, user *User) error
}

type RoomRepository interface {
	GetRoom(ctx context.Context, roomID string) (*Room, error)
	SaveRoom(ctx context.Context, room *Room) (*Room, error)
	DeleteRoom(ctx context.Context, room *Room) error
}

type SensorRepository interface {
	GetSensor(ctx context.Context, sensorID string) (*Sensor, error)
	SaveSensor(ctx context.Context, sensor *Sensor) (*Sensor, error)
	DeleteSensor(ctx context.Context, sensor *Sensor) error
}
