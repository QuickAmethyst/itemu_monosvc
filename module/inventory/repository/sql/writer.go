package sql

import (
	"context"
	"fmt"
	"github.com/QuickAmethyst/monosvc/module/inventory/domain"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
)

type Writer interface {
	StoreUom(ctx context.Context, uom *domain.Uom) (err error)
	UpdateUomByID(ctx context.Context, id int64, uom *domain.Uom) (err error)
}

type writer struct {
	logger logger.Logger
	db     sql.DB
}

func (w *writer) UpdateUomByID(ctx context.Context, id int64, uom *domain.Uom) (err error) {
	_, err = w.updateUom(ctx, uom, UomStatement{ID: id})
	return
}

func (w *writer) updateUom(ctx context.Context, uom *domain.Uom, where UomStatement) (result sql.Result, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateUomFailed, "Failed on build where clause")
		return
	}

	updateQuery := fmt.Sprintf("UPDATE uoms SET name = ?, description = ?, decimal = ? %s", whereClause)

	args := []interface{}{uom.Name, uom.Description, uom.Decimal}
	result, err = w.db.ExecContext(ctx, w.db.Rebind(updateQuery), append(args, whereClauseArgs...)...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateUomFailed, "Update uom failed")
		return
	}

	return
}

func (w *writer) StoreUom(ctx context.Context, uom *domain.Uom) (err error) {
	err = w.db.QueryRowContext(ctx, `
		INSERT INTO uoms (name, description, decimal)
		VALUES ($1, $2, $3)
		RETURNING id
	`, uom.Name, uom.Description, uom.Decimal).Scan(&uom.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreUomFailed, "Insert uom failed")
		return
	}

	return
}

func NewWriter(opt *Options) Writer {
	return &writer{opt.Logger, opt.MasterDB}
}
