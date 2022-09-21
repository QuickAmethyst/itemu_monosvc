package model

import "github.com/QuickAmethyst/monosvc/module/accounting/domain"

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

func (w *WriteAccountClassInput) Domain() (accountClass domain.AccountClass, err error) {
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

type AccountClassesResult struct {
	Data   []*AccountClass `json:"data"`
	Paging *Paging         `json:"paging"`
}

type AccountClassTypeInput struct {
	ID int64 `json:"id"`
}

type AccountClassTypesResult struct {
	Data []*AccountClassType `json:"data"`
}
