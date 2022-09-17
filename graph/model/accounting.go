package model

import "github.com/QuickAmethyst/monosvc/module/accounting/domain"

type AccountClass struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type uint   `json:"type"`
}

type WriteAccountClassesInput struct {
	Name string `json:"name"`
	Type uint   `json:"type"`
}

func (w *WriteAccountClassesInput) Domain() (accountClass domain.AccountClass, err error) {
	accountClass.Name = w.Name
	accountClass.Type = domain.ClassType(w.Type)
	return
}
