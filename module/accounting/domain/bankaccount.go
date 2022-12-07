package domain

import "database/sql"

const (
	CashAccountType int64 = iota + 1
	ChequingAccountType
	SavingAccountType
	CreditAccountType
)

type BankAccount struct {
	ID         int64
	AccountID  int64 `db:"account_id"`
	Type       int64
	BankNumber sql.NullString
	Inactive   bool
}
