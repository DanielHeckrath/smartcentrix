package main

import (
	"fmt"
	"os"
	"time"

	smartcentrix "github.com/DanielHeckrath/smartcentrix/proto"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/juju/errors"
)

var (
	envDatabaseHost     = "MYSQL_HOST"
	envDatabaseName     = "MYSQL_DATABASE"
	envDatabaseUser     = "MYSQL_USER"
	envDatabasePassword = "MYSQL_PASSWORD"
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

	opts := fmt.Sprintf("%[1]s:%[2]s@tcp(%[3]s)/%[4]s?charset=utf8&parseTime=True&loc=Local", username, password, host, database)
	db, err := gorm.Open("mysql", opts)

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

// BaseModel adds basic columns for our models
type BaseModel struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// User represents a single user in our database
type User struct {
	BaseModel

	Name     string
	Password string

	Sensors []Sensor `gorm:"ForeignKey:UserID"`
	Rooms   []Room   `gorm:"ForeignKey:UserID"`
	Devices []Device `gorm:"ForeignKey:UserID"`
}

// Device represents a users mobile device
type Device struct {
	BaseModel

	Type      int64
	PushToken string

	UserID string `gorm:"index"`
}

// Room is a mechanic to group several sensors
type Room struct {
	BaseModel

	Name    string
	Sensors []Sensor `gorm:"ForeignKey:RoomID"`

	UserID string `gorm:"index"`
}

func (r Room) proto() *smartcentrix.Room {
	room := smartcentrix.Room{
		Id:   r.ID,
		Name: r.Name,
	}

	return &room
}

func convertRooms(r []*Room) []*smartcentrix.Room {
	rooms := make([]*smartcentrix.Room, 0, len(r))

	for _, room := range r {
		rooms = append(rooms, room.proto())
	}

	return rooms
}

// Sensor is a single smartcentrix sensor
type Sensor struct {
	BaseModel

	Name            string
	LastMeasurement int64
	Status          bool
	InUse           bool
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
		sensor.RoomId = *s.RoomID
	}

	return &sensor
}

func convertSensors(s []*Sensor) []*smartcentrix.Sensor {
	sensors := make([]*smartcentrix.Sensor, 0, len(s))

	for _, sensor := range s {
		sensors = append(sensors, sensor.proto())
	}

	return sensors
}

// Measurement is a single entry in a sensors measurements time series
type Measurement struct {
	BaseModel

	Timestamp int64
	Value     int64

	SensorID string `gorm:"index"`
}

func (m Measurement) proto() *smartcentrix.Measurement {
	measurement := smartcentrix.Measurement{
		Id:        m.ID,
		Timestamp: m.Timestamp,
		Value:     m.Value,
	}

	return &measurement
}

func convertMeasurements(m []Measurement) []*smartcentrix.Measurement {
	measurements := make([]*smartcentrix.Measurement, 0, len(m))

	for _, measurement := range m {
		measurements = append(measurements, measurement.proto())
	}

	return measurements
}
