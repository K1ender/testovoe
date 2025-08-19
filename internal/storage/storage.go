package storage

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage struct {
	Subscription SubscriptionStorage
}

func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{Subscription: NewPostgresSubscriptionStorage(db)}
}
