package sql

import "github.com/gofrs/uuid"

type UserStatement struct {
	ID    uuid.UUID
	Email string
}

type UserProfileStatement struct {
	UserID uuid.UUID
}
