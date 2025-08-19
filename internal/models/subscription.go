package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     sql.NullTime
}
