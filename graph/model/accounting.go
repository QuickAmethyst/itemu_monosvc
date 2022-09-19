package model

import "github.com/QuickAmethyst/monosvc/module/accounting/domain"

type AccountClass struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	TypeID   uint   `json:"typeID"`
	Inactive bool   `json:"inactive"`
}

type WriteAccountClassesInput struct {
	Name     string `json:"name"`
	TypeID   uint   `json:"typeID"`
	Inactive bool   `json:"inactive"`
}

func (w *WriteAccountClassesInput) Domain() (accountClass domain.AccountClass, err error) {
	accountClass.Name = w.Name
	accountClass.TypeID = domain.ClassType(w.TypeID)
	accountClass.Inactive = w.Inactive
	return
}
