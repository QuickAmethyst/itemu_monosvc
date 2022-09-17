package domain

type AccountClass struct {
	ID       int64
	Name     string
	Type     ClassType
	Inactive bool
}
