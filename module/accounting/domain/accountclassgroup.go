package domain

type AccountGroup struct {
	ID       int64
	ParentID int64
	ClassID  int64
	Name     string
	Inactive bool
}
