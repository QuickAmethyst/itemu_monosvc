package sql

import (
	"context"
	"fmt"
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
)

type Writer interface {
	StoreAccountClass(ctx context.Context, accountClass *domain.AccountClass) (err error)
	UpdateAccountClassByID(ctx context.Context, id int64, accountClass *domain.AccountClass) (err error)
}

type writer struct {
	logger logger.Logger
	db     sql.DB
}

func (w *writer) UpdateAccountClassByID(ctx context.Context, id int64, accountClass *domain.AccountClass) (err error) {
	return w.updateAccountClass(ctx, accountClass, AccountClassStatement{ID: id})
}

func (w *writer) updateAccountClass(ctx context.Context, accountClass *domain.AccountClass, where AccountClassStatement) (err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountClassFailed, "Failed on select account class")
		return
	}

	updateQuery := fmt.Sprintf(`
		UPDATE account_classes
		SET name = ?, type = ?
		%s
	`, whereClause)
	stmt, err := w.db.PrepareContext(ctx, w.db.Rebind(updateQuery))
	if err != nil {
		err = errors.PropagateWithCode(err, EcodePreparedStatementFailed, "Prepare statement failed")
		return
	}

	args := []interface{}{accountClass.Name, accountClass.Type}
	_, err = stmt.ExecContext(ctx, append(args, whereClauseArgs...)...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountClassFailed, "Exec statement failed")
		return
	}

	return
}

func (w *writer) StoreAccountClass(ctx context.Context, accountClass *domain.AccountClass) (err error) {
	stmt, err := w.db.PrepareContext(ctx, `
		INSERT INTO account_classes (name, type, inactive)
		VALUES ($1, $2, $3)
		RETURNING ID
	`)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodePreparedStatementFailed, "Prepare statement failed")
		return
	}

	if err = stmt.GetContext(ctx, &accountClass.ID, accountClass.Name, accountClass.Type, accountClass.Inactive); err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreAccountClassFailed, "Exec statement failed")
		return
	}

	return
}

func NewWriter(opt *Options) Writer {
	return &writer{opt.Logger, opt.MasterDB}
}
