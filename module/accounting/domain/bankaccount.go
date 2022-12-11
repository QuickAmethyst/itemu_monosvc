package domain

import "database/sql"

type BankAccount struct {
	ID         int64
	AccountID  int64 `db:"account_id"`
	TypeID     int64 `db:"type_id"`
	BankNumber sql.NullString
	Inactive   bool
}
