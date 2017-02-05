package main

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

type sqlSensorRepository struct {
	db *gorm.DB
}

func (r *sqlSensorRepository) GetSensor(ctx context.Context, sensorID string) (*Sensor, error) {
	sensor := &Sensor{}
	q := r.db.Where("id = ?", sensorID).First(sensor)

	if q.Error != nil {
		if q.RecordNotFound() {
			return nil, errRecordNotFound
		}
		return nil, q.Error
	}

	return sensor, nil
}

func (r *sqlSensorRepository) SaveSensor(ctx context.Context, sensor *Sensor) (*Sensor, error) {
	q := r.db.Save(sensor)

	if q.Error != nil {
		return nil, q.Error
	}

	return sensor, nil
}

func (r *sqlSensorRepository) DeleteSensor(ctx context.Context, sensor *Sensor) error {
	q := r.db.Delete(sensor)

	if q.Error != nil {
		return q.Error
	}

	return nil
}
