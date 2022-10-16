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
	Paging *PagingInput `json:"paging"`
}

type AccountClassTypeInput struct {
	ID int64 `json:"id"`
}

type AccountClassTypesResult struct {
	Data []*AccountClassType `json:"data"`
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
	Inactive *bool  `json:"inactive"`
}

func (w *WriteAccountGroupInput) Domain() (accountGroup domain.AccountGroup, err error) {
	accountGroup.Name = w.Name
	accountGroup.ClassID = w.ClassID

	if w.ParentID != nil {
		if err = accountGroup.ParentID.Scan(*w.ParentID); err != nil {
			return
		}
	}

	if w.Inactive != nil {
		accountGroup.Inactive = *w.Inactive
	}

	return
}

type AccountGroupInput struct {
	ID             int  `json:"id"`
	ParentIDIsNULL bool `json:"parentIDIsNULL"`
}

type Account struct {
	ID       int64 `json:"id"`
	Name     string `json:"name"`
	GroupID  int64    `json:"groupID"`
	Inactive bool  `json:"inactive"`
}

type WriteAccountInput struct {
	Name     string `json:"name"`
	GroupID  int64    `json:"groupID"`
	Inactive *bool  `json:"inactive"`
}

func (w *WriteAccountInput) Domain()  (account domain.Account) {
	account.Name = w.Name
	account.GroupID = w.GroupID

	if w.Inactive != nil {
		account.Inactive = *w.Inactive
	}

	return
}

type AccountInput struct {
	ID int64 `json:"id"`
}

type Journal struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}

type WriteTransactionsInput struct {
	AccountID int64     `json:"accountID"`
	Amount    float64 `json:"amount"`
}
