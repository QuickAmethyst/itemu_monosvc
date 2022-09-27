package model

import (
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
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
