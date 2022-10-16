package domain

import (
	"github.com/google/uuid"
	"time"
)

type Journal struct {
	ID        uuid.UUID
	AccountID int64
	Memo      string
	Amount    float64
	CreatedAt time.Time
}
