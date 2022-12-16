package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Journal struct {
	ID        uuid.UUID
	Amount    float64
	CreatedAt time.Time
	TransDate time.Time
	Memo      sql.NullString
	DeletedAt time.Time
}
