package model

import (
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"time"
)

type AccountClass struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	TypeID   int64  `json:"typeID"`
	Inactive bool   `json:"inactive"`
}

type WriteAccountClassInput struct {
	Name     string `json:"name"`
	TypeID   int64  `json:"typeID"`
	Inactive bool   `json:"inactive"`
}

func (w *WriteAccountClassInput) Domain() (accountClass domain.AccountClass) {
	accountClass.Name = w.Name
	accountClass.TypeID = w.TypeID
	accountClass.Inactive = w.Inactive
	return
}

type AccountClassInput struct {
	ID int `json:"id"`
}

type AccountClassType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type AccountClassesInput struct {
	Paging PagingInput `json:"paging"`
}

type AccountClassTypeInput struct {
	ID int64 `json:"id"`
}

type AccountClassTypesResult struct {
	Data []AccountClassType `json:"data"`
}

type AccountGroup struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ClassID  int64  `json:"classID"`
	ParentID int64  `json:"parentID"`
	Inactive bool   `json:"inactive"`
}

type WriteAccountGroupInput struct {
	Name     string `json:"name"`
	ClassID  int64  `json:"classID"`
	ParentID *int64 `json:"parentID"`
	Inactive bool   `json:"inactive"`
}

func (w *WriteAccountGroupInput) Domain() (accountGroup domain.AccountGroup, err error) {
	accountGroup.Name = w.Name
	accountGroup.ClassID = w.ClassID
	accountGroup.Inactive = w.Inactive

	if w.ParentID != nil && *w.ParentID > 0 {
		err = accountGroup.ParentID.Scan(*w.ParentID)
		if err != nil {
			return
		}
	}

	return
}

type AccountGroupInput struct {
	ID             int  `json:"id"`
	ParentIDIsNULL bool `json:"parentIDIsNULL"`
}

type Account struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	GroupID  int64  `json:"groupID"`
	Inactive bool   `json:"inactive"`
}

type WriteAccountInput struct {
	Name     string `json:"name"`
	GroupID  int64  `json:"groupID"`
	Inactive bool   `json:"inactive"`
}

func (w *WriteAccountInput) Domain() (account domain.Account) {
	account.Name = w.Name
	account.GroupID = w.GroupID
	account.Inactive = w.Inactive
	return
}

type AccountInput struct {
	ID        int64 `json:"id"`
	ClassType int64 `json:"classType"`
}

type Journal struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	TransDate time.Time `json:"transDate"`
	CreatedAt time.Time `json:"createdAt"`
}

type WriteTransactionRow struct {
	AccountID int64   `json:"accountID"`
	Amount    float64 `json:"amount"`
}

type WriteTransactionInput struct {
	TransDate time.Time             `json:"transDate"`
	Memo      string                `json:"memo"`
	Data      []WriteTransactionRow `json:"data"`
}

type GeneralLedgerPreference struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"accountID"`
}

type WriteGeneralLedgerPreferenceInput struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"accountID"`
}

func (w *WriteGeneralLedgerPreferenceInput) Domain() (preference domain.GeneralLedgerPreference, err error) {
	preference.ID = w.ID
	if err = preference.AccountID.Scan(w.AccountID); err != nil {
		return
	}

	return
}

type GeneralLedgerPreferenceInput struct {
	ID int64 `json:"id"`
}

type FiscalYear struct {
	ID        int64     `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Closed    bool      `json:"closed"`
}

type WriteFiscalYearInput struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Closed    bool      `json:"closed"`
}

func (w *WriteFiscalYearInput) Domain() (fiscalYear domain.FiscalYear) {
	fiscalYear.StartDate = w.StartDate
	fiscalYear.EndDate = w.EndDate
	fiscalYear.Closed = w.Closed
	return
}

type FiscalYearsInput struct {
	Paging PagingInput `json:"paging"`
}

type FiscalYearsResult struct {
	Data   []FiscalYear `json:"data"`
	Paging Paging       `json:"paging"`
}

type BankAccount struct {
	ID         int64  `json:"id"`
	AccountID  int64  `json:"accountID"`
	TypeID     int64  `json:"typeID"`
	BankNumber string `json:"bankNumber"`
	Inactive   bool   `json:"inactive"`
}

type WriteBankAccountInput struct {
	AccountID  int64  `json:"accountID"`
	TypeID     int64  `json:"typeID"`
	BankNumber string `json:"bankNumber"`
	Inactive   bool   `json:"inactive"`
}

func (w *WriteBankAccountInput) Domain() (bankAccount domain.BankAccount, err error) {
	bankAccount.AccountID = w.AccountID
	bankAccount.TypeID = w.TypeID
	bankAccount.Inactive = w.Inactive

	if w.BankNumber != "" {
		if err = bankAccount.BankNumber.Scan(w.BankNumber); err != nil {
			return
		}
	}

	return
}

type BankAccountInput struct {
	ID int64 `json:"id"`
}

type BankAccountType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type BankAccountTypesResult struct {
	Data []BankAccountType `json:"data"`
}

type BankAccountsResult struct {
	Data   []BankAccount `json:"data"`
	Paging Paging        `json:"paging"`
}

type BankAccountsInputScope struct {
	ID int64 `json:"id"`
}

type BankAccountsInput struct {
	Scope  BankAccountsInputScope `json:"scope"`
	Paging PagingInput            `json:"paging"`
}

type BankTransaction struct {
	ID            int64     `json:"id"`
	JournalID     string    `json:"journalID"`
	BankAccountID int64     `json:"bankAccountID"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"createdAt"`
}

type WriteBankTransactionInput struct {
	BankAccountID int64                 `json:"bankAccountID"`
	TransDate     time.Time             `json:"transDate"`
	Memo          string                `json:"memo"`
	Data          []WriteTransactionRow `json:"data"`
}

type AccountClassTransactionIDInput struct {
	Paging *PagingInput `json:"paging"`
}

type AccountClassTransactionResult struct {
	Data        []*TransactionRow `json:"data"`
	TotalAmount float64           `json:"totalAmount"`
}

type TransactionRow struct {
	AccountID int64   `json:"accountID"`
	Amount    float64 `json:"amount"`
}
