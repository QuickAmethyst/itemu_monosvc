package domain

import (
	"github.com/google/uuid"
)

type GeneralLedger struct {
	ID        uuid.UUID
	JournalID uuid.UUID
	AccountID int64
	Amount    float64
	CreatedBy uuid.UUID
}
