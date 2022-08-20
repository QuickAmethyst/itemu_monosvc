package model

import (
	"github.com/QuickAmethyst/monosvc/module/inventory/domain"
)

type Uom struct {
	ID          int64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Decimal     int32  `json:"decimal"`
}

type UomsInput struct {
	Paging *PagingInput `json:"paging"`
}

type UomsResult struct {
	Data   []*Uom  `json:"data"`
	Paging *Paging `json:"paging"`
}

type WriteUomInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Decimal     int32  `json:"decimal"`
}

func (w *WriteUomInput) Domain() (uom domain.Uom, err error) {
	uom.Name = w.Name

	if err = uom.Description.Scan(w.Description); err != nil {
		return
	}

	if err = uom.Decimal.Scan(w.Decimal); err != nil {
		return
	}

	return
}
