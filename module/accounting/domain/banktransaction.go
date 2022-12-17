package domain

import (
	"github.com/google/uuid"
	"time"
)

type BankTransaction struct {
	ID            int64
	JournalID     uuid.UUID
	BankAccountID int64
	UserID        int64
	Amount        float64
	Balance       float64
	CreatedAt     time.Time
}
