package domain

import "database/sql"

type GeneralLedgerPreference struct {
	ID        int64
	AccountID sql.NullInt64 `db:"account_id"`
}
