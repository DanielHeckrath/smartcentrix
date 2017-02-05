package main

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

type sqlMeasurementRepository struct {
	db *gorm.DB
}

func (r *sqlMeasurementRepository) ListMeasurements(ctx context.Context, sensorID string) ([]*Measurement, error) {
	var measurements []*Measurement
	q := r.db.Where("sensor_id = ?", sensorID).Order("timestamp desc").Find(&measurements)

	if q.Error != nil {
		if q.RecordNotFound() {
			return nil, errRecordNotFound
		}
		return nil, q.Error
	}

	return measurements, nil
}

func (r *sqlMeasurementRepository) SaveMeasurement(ctx context.Context, measurement *Measurement) (*Measurement, error) {
	var q *gorm.DB

	if r.db.NewRecord(measurement) {
		q = r.db.Create(measurement)
	} else {
		q = r.db.Save(measurement)
	}

	if q.Error != nil {
		return nil, q.Error
	}

	return measurement, nil
}
