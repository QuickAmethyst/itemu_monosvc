package domain

import (
	"context"
	sql "database/sql"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
)

type Uom struct {
	ID          int64
	Name        string
	Description sql.NullString
	Decimal     sql.NullInt32
}

type UomRepository interface {
	GetList(ctx context.Context, uom Uom) ([]Uom, qb.Paging, error)
}

type UomUsecase interface {
	GetList(ctx context.Context, uom Uom) ([]Uom, qb.Paging, error)
}