package domain

const (
	CashAccountType int64 = iota + 1
	ChequingAccountType
	SavingAccountType
	CreditAccountType
)

type BankAccountType struct {
	ID   int64
	Name string
}
