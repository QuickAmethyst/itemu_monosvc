package domain

import "database/sql"

type BankAccount struct {
	ID         int64
	AccountID  int64 `db:"account_id"`
	Type       int64
	BankNumber sql.NullString
	Inactive   bool
}
