package domain

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

func (u *User) PasswordValid(v string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(v))
}

func (u *User) SetPassword(pass string) (err error) {
	var hashedPass []byte
	hashedPass, err = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.Password = string(hashedPass)
	return
}
