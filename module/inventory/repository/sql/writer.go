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
	return w.updateUom(ctx, uom, UomStatement{ID: id})
}

func (w *writer) updateUom(ctx context.Context, uom *domain.Uom, where UomStatement) (err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateUomFailed, "Failed on select uom")
		return
	}

	updateQuery := fmt.Sprintf("UPDATE uoms SET name = ?, description = ?, decimal = ? %s", whereClause)
	stmt, err := w.db.PrepareContext(ctx, w.db.Rebind(updateQuery))
	if err != nil {
		err = errors.PropagateWithCode(err, EcodePreparedStatementFailed, "Prepare statement failed")
		return
	}

	args := []interface{}{uom.Name, uom.Description, uom.Decimal}
	_, err = stmt.ExecContext(ctx, append(args, whereClauseArgs...)...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateUomFailed, "Exec statement failed")
		return
	}

	return
}

func (w *writer) StoreUom(ctx context.Context, uom *domain.Uom) (err error) {
	stmt, err := w.db.PrepareContext(ctx, "INSERT INTO uoms (name, description, decimal) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		err = errors.PropagateWithCode(err, EcodePreparedStatementFailed, "Prepare statement failed")
		return
	}

	if err = stmt.GetContext(ctx, &uom.ID, uom.Name, uom.Description, uom.Decimal); err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreUomFailed, "Exec statement failed")
		return
	}

	return
}

func NewWriter(opt *Options) Writer {
	return &writer{opt.Logger, opt.MasterDB}
}
