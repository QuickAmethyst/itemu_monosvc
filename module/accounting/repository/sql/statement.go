package sql

import (
	"time"
)

type AccountClassStatement struct {
	ID int64
}

type AccountGroupStatement struct {
	ID             int64
	ParentID       int64
	ParentIDIsNULL bool
}

type AccountStatement struct {
	ID        int64
	ClassType int64 `qb:"-"`
}

type GeneralLedgerPreferenceStatement struct {
	ID int64
}

type FiscalYearStatement struct {
	ID           int64
	ClosedNotEQ  bool
	StartDateGTE time.Time
	StartDateLTE time.Time
	EndDateGTE   time.Time
	EndDateLTE   time.Time
}

type BankAccountStatement struct {
	ID int64
}
