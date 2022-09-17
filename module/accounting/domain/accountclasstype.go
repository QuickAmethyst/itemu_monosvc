package domain

type ClassType int

const (
	AssetClassType ClassType = iota + 1
	LiabilitiesClassType
	EquityClassType
	IncomeClassType
	COGSClassType
	ExpenseClassType
)

type AccountClassType struct {
	ID   int
	Name string
}
