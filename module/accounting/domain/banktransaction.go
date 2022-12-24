package domain

import (
	"github.com/google/uuid"
	"time"
)

type BankTransaction struct {
	ID            int64
	JournalID     uuid.UUID
	BankAccountID int64
	UserID        uuid.UUID
	Amount        float64
	Balance       float64
	Memo          string
	TransDate     time.Time
	CreatedAt     time.Time
}
