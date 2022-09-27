package domain

import "database/sql"

type AccountGroup struct {
	ID       int64
	ParentID sql.NullInt64 `db:"parent_id"`
	ClassID  int64         `db:"class_id"`
	Name     string
	Inactive bool
}
