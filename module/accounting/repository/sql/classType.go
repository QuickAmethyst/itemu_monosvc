package sql

import "github.com/QuickAmethyst/monosvc/module/accounting/domain"

const (
	AssetClassType int64 = iota + 1
	LiabilitiesClassType
	EquityClassType
	IncomeClassType
	COGSClassType
	ExpenseClassType
)

var classTypes = map[int64]domain.AccountClassType{
	AssetClassType:       {AssetClassType, "Asset"},
	LiabilitiesClassType: {LiabilitiesClassType, "Liability"},
	EquityClassType:      {EquityClassType, "Equity"},
	IncomeClassType:      {IncomeClassType, "Income"},
	COGSClassType:        {COGSClassType, "Cost of Good Solds"},
	ExpenseClassType:     {ExpenseClassType, "Expense"},
}

