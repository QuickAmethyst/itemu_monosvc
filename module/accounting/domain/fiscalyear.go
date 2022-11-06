package domain

import "time"

type FiscalYear struct {
	ID        int64
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	Closed    bool
}
