package domain

import "time"

type FiscalYear struct {
	ID        int64
	StartDate time.Time
	EndDate   time.Time
	Closed    bool
}
