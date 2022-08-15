package usecase

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/inventory/domain"
	"github.com/QuickAmethyst/monosvc/module/inventory/repository/sql"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
)

type Reader interface {
	GetUomList(ctx context.Context, stmt sql.UomStatement, p qb.Paging) (result []domain.Uom, paging qb.Paging, err error)
}

type reader struct {
	InventorySQL sql.SQL
}

func (r *reader) GetUomList(ctx context.Context, stmt sql.UomStatement, p qb.Paging) (result []domain.Uom, paging qb.Paging, err error) {
	return r.InventorySQL.GetUomList(ctx, stmt, p)
}

func NewReader(opt *Options) Reader {
	return &reader{opt.InventorySQL}
}
