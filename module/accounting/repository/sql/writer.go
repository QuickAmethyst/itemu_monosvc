package sql

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
)

type Writer interface {
	StoreAccountClass(ctx context.Context, accountClass *domain.AccountClass) (err error)
	UpdateAccountClassByID(ctx context.Context, id string, accountClass *domain.AccountClass) (err error)
}

type writer struct {
	logger logger.Logger
	db     sql.DB
}

func (w *writer) UpdateAccountClassByID(ctx context.Context, id string, accountClass *domain.AccountClass) (err error) {
	panic("implement me")
}

func (w *writer) StoreAccountClass(ctx context.Context, accountClass *domain.AccountClass) (err error) {
	stmt, err := w.db.PrepareContext(ctx, `
		INSERT INTO account_classes (name, type)
		VALUES ($1, $2)
		RETURNING ID
	`)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodePreparedStatementFailed, "Prepare statement failed")
		return
	}

	if err = stmt.GetContext(ctx, &accountClass.ID, accountClass.Name, accountClass.Type); err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreAccountClassesFailed, "Exec statement failed")
		return
	}

	return
}

func NewWriter(opt *Options) Writer {
	return &writer{opt.Logger, opt.MasterDB}
}
