package domain

import "time"

type UserProfile struct {
	ID        int64
	FullName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
