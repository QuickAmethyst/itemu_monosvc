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
	DeleteAccountClassByID(ctx context.Context, id int64) (err error)
}

type writer struct {
	logger logger.Logger
	db     sql.DB
}

func (w *writer) DeleteAccountClassByID(ctx context.Context, id int64) (err error) {
	_, err = w.deleteAccountClass(ctx, AccountClassStatement{ID: id})
	return
}

func (w *writer) deleteAccountClass(ctx context.Context, where AccountClassStatement) (result sql.Result, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountClassFailed, "Failed on select account class")
		return
	}

	deleteQuery := fmt.Sprintf(`
		DELETE FROM account_classes
		%s
	`, whereClause)

	result, err = w.db.ExecContext(ctx, w.db.Rebind(deleteQuery), whereClauseArgs...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeDeleteAccountClassFailed, "Update account class failed")
		return
	}

	return
}

func (w *writer) UpdateAccountClassByID(ctx context.Context, id int64, accountClass *domain.AccountClass) (err error) {
	_, err = w.updateAccountClass(ctx, accountClass, AccountClassStatement{ID: id})
	return
}

func (w *writer) updateAccountClass(ctx context.Context, accountClass *domain.AccountClass, where AccountClassStatement) (result sql.Result, err error) {
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

	args := []interface{}{accountClass.Name, accountClass.Type}
	result, err = w.db.ExecContext(ctx, w.db.Rebind(updateQuery), append(args, whereClauseArgs...)...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountClassFailed, "Update account class failed")
		return
	}

	return
}

func (w *writer) StoreAccountClass(ctx context.Context, accountClass *domain.AccountClass) (err error) {
	err = w.db.QueryRowContext(ctx, `
		INSERT INTO account_classes (name, type, inactive)
		VALUES ($1, $2, $3)
		RETURNING id
	`, accountClass.Name, accountClass.Type, accountClass.Inactive).Scan(&accountClass.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreAccountClassFailed, "Insert account class failed")
		return
	}

	return
}

func NewWriter(opt *Options) Writer {
	return &writer{opt.Logger, opt.MasterDB}
}
