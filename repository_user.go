package main

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

type sqlUserRepository struct {
	db *gorm.DB
}

func (r *sqlUserRepository) RegisterUser(ctx context.Context, user *User) (*User, error) {
	q := r.db.Create(user)

	if q.Error != nil {
		return nil, q.Error
	}

	return user, nil
}

func (r *sqlUserRepository) GetUser(ctx context.Context, userID string) (*User, error) {
	user := &User{}
	q := r.db.Where("id = ?", userID).First(user)

	if q.Error != nil {
		if q.RecordNotFound() {
			return nil, errRecordNotFound
		}
		return nil, q.Error
	}

	return user, nil
}

func (r *sqlUserRepository) GetUserWithName(ctx context.Context, name string) (*User, error) {
	user := &User{}
	q := r.db.Where("name = ?", name).First(user)

	if q.Error != nil {
		if q.RecordNotFound() {
			return nil, errRecordNotFound
		}
		return nil, q.Error
	}

	return user, nil
}

func (r *sqlUserRepository) GetUserWithCredentials(ctx context.Context, name, password string) (*User, error) {
	user := &User{}
	q := r.db.Where("name = ? AND password = ?", name, password).First(user)

	if q.Error != nil {
		if q.RecordNotFound() {
			return nil, errRecordNotFound
		}
		return nil, q.Error
	}

	return user, nil
}

func (r *sqlUserRepository) SaveUser(ctx context.Context, user *User) (*User, error) {
	q := r.db.Save(user)

	if q.Error != nil {
		return nil, q.Error
	}

	return user, nil
}

func (r *sqlUserRepository) DeleteUser(ctx context.Context, user *User) error {
	q := r.db.Delete(user)

	if q.Error != nil {
		return q.Error
	}

	return nil
}
