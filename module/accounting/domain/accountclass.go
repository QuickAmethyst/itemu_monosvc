package domain

type AccountClass struct {
	ID       int64
	Name     string
	TypeID   ClassType
	Inactive bool
}
