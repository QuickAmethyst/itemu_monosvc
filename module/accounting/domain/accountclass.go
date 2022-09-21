package domain

type AccountClass struct {
	ID       int64
	Name     string
	TypeID   int64 `db:"type_id"`
	Inactive bool
}
