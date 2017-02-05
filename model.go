package main

import (
	"fmt"
	"os"
	"time"

	smartcentrix "github.com/DanielHeckrath/smartcentrix/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/juju/errors"
)

var (
	envDatabaseHost     = "POSTGRES_HOST"
	envDatabaseName     = "POSTGRES_DATABASE"
	envDatabaseUser     = "POSTGRES_USER"
	envDatabasePassword = "POSTGRES_PASSWORD"
)

var (
	errHostEmpty     = dbError("host", envDatabaseHost)
	errNameEmpty     = dbError("database", envDatabaseName)
	errUserEmpty     = dbError("user", envDatabaseUser)
	errPasswordEmpty = dbError("password", envDatabasePassword)
)

func dbError(key string, envVar string) error {
	return errors.Errorf("Database %s cannot be empty. Please set %s via %s environment variable", key, key, envVar)
}

func newDatabase() (*gorm.DB, error) {
	// read databse host address from environment variables
	host := os.Getenv(envDatabaseHost)

	if host == "" {
		return nil, errHostEmpty
	}

	// read databse user address from environment variables
	username := os.Getenv(envDatabaseUser)

	if username == "" {
		return nil, errUserEmpty
	}

	// read databse password address from environment variables
	password := os.Getenv(envDatabasePassword)

	if password == "" {
		return nil, errPasswordEmpty
	}

	// read databse name address from environment variables
	database := os.Getenv(envDatabaseName)

	if database == "" {
		return nil, errNameEmpty
	}

	opts := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, username, database, password)
	db, err := gorm.Open("postgres", opts)

	if err != nil {
		return nil, err
	}

	var (
		user        = &User{}
		device      = &Device{}
		room        = &Room{}
		sensor      = &Sensor{}
		measurement = &Measurement{}
	)

	db.AutoMigrate(
		user,
		device,
		room,
		sensor,
		measurement,
	)

	// add foreign keys - this is not done during auto migration
	db.Model(device).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(sensor).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(room).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	// sensors can exist without an associated room, so set to null on delete
	db.Model(sensor).AddForeignKey("room_id", "rooms(id)", "SET NULL", "CASCADE")
	db.Model(measurement).AddForeignKey("sensor_id", "sensors(id)", "CASCADE", "CASCADE")

	return db, nil
}

type BaseModel struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type User struct {
	BaseModel

	Name     string
	Password string

	Sensors []Sensor `gorm:"ForeignKey:UserID"`
	Rooms   []Room   `gorm:"ForeignKey:UserID"`
	Devices []Device `gorm:"ForeignKey:UserID"`
}

type Device struct {
	BaseModel

	Type      int64
	PushToken string

	UserID string `gorm:"index"`
}

type Room struct {
	BaseModel

	Name    string
	Sensors []Sensor `gorm:"ForeignKey:RoomID"`

	UserID string `gorm:"index"`
}

type Sensor struct {
	BaseModel

	Name            string
	LastMeasurement int64
	Status          bool
	Measurements    []Measurement `gorm:"ForeignKey:SensorID"`

	UserID string `gorm:"index"`
	// RoomID must be a pointer for optional foreign key constraints
	RoomID *string
}

func (s Sensor) proto() *smartcentrix.Sensor {
	sensor := smartcentrix.Sensor{
		Id:              s.ID,
		Name:            s.Name,
		LastMeasurement: s.LastMeasurement,
		Status:          s.Status,
	}

	if s.RoomID != nil {
		sensor.RoomId = &wrappers.StringValue{
			Value: *s.RoomID,
		}
	}

	return &sensor
}

type Measurement struct {
	BaseModel

	Timestamp int64
	Value     int64

	SensorID string `gorm:"index"`
}
