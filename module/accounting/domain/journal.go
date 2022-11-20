package domain

import (
	"github.com/google/uuid"
	"time"
)

type Journal struct {
	ID        uuid.UUID
	Amount    float64
	TransDate time.Time
	CreatedAt time.Time
}
