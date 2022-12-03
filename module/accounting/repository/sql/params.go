package sql

import "time"

type TransactionRow struct {
	AccountID int64
	Amount    float64
}

type Transaction struct {
	Date time.Time
	Memo string
	Data []TransactionRow
}
