package sql

import (
	"github.com/google/uuid"
	"time"
)

type TransactionRow struct {
	AccountID int64
	Amount    float64
}

type Transaction struct {
	Date      time.Time
	Memo      string
	Data      []TransactionRow
	journalID uuid.UUID
}

type BankTransaction struct {
	BankAccountID       int64
	bankTransactionType BankTransactionType
	Transaction
}
