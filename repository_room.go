package main

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type sqlRoomRepository struct {
	db *gorm.DB
}

func (r *sqlRoomRepository) GetRoom(ctx context.Context, roomID string) (*Room, error) {
	room := &Room{}
	q := r.db.Where("id = ?", roomID).First(room)

	if q.Error != nil {
		if q.RecordNotFound() {
			return nil, errRecordNotFound
		}
		return nil, q.Error
	}

	return room, nil
}

func (r *sqlRoomRepository) ListRooms(ctx context.Context, userID string) ([]*Room, error) {
	var rooms []*Room
	q := r.db.Where("user_id = ?", userID).Find(&rooms)

	if q.Error != nil {
		if q.RecordNotFound() {
			return nil, errRecordNotFound
		}
		return nil, q.Error
	}

	return rooms, nil
}

func (r *sqlRoomRepository) SaveRoom(ctx context.Context, room *Room) (*Room, error) {
	var q *gorm.DB

	if r.db.NewRecord(room) {
		room.ID = uuid.NewV4().String()
		q = r.db.Create(room)
	} else {
		q = r.db.Save(room)
	}

	if q.Error != nil {
		return nil, q.Error
	}

	return room, nil
}

func (r *sqlRoomRepository) DeleteRoom(ctx context.Context, room *Room) error {
	q := r.db.Delete(room)

	if q.Error != nil {
		return q.Error
	}

	return nil
}
