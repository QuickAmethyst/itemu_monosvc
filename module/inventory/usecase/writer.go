package usecase

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/inventory/domain"
	"github.com/QuickAmethyst/monosvc/module/inventory/repository/sql"
)

type Writer interface {
	StoreUom(ctx context.Context, uom *domain.Uom) (err error)
}

type writer struct {
	InventorySQL sql.SQL
}

func (w *writer) StoreUom(ctx context.Context, uom *domain.Uom) (err error) {
	return w.InventorySQL.StoreUom(ctx, uom)
}

func NewWriter(opt *Options) Writer {
	return &writer{opt.InventorySQL}
}