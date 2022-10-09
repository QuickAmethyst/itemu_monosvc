package domain

type Account struct {
	ID       int64
	Name     string
	GroupID  int64 `db:"group_id"`
	Inactive bool
}
