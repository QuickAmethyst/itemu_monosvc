package sql

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/inventory/domain"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
)

type Writer interface {
	StoreUom(ctx context.Context, uom *domain.Uom) (err error)
}

type writer struct {
	logger logger.Logger
	db     sql.DB
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
