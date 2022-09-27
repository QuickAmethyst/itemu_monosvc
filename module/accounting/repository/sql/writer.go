package sql

import (
	"context"
	goErr "errors"
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

	StoreAccountGroup(ctx context.Context, accountClassGroup *domain.AccountGroup) (err error)
	UpdateAccountGroupByID(ctx context.Context, id int64, accountGroup *domain.AccountGroup) (err error)
	DeleteAccountGroupByID(ctx context.Context, id int64) (err error)
}

type writer struct {
	logger logger.Logger
	db     sql.DB
	reader Reader
}

func (w *writer) StoreAccountGroup(ctx context.Context, accountGroup *domain.AccountGroup) (err error) {
	if accountGroup.ParentID.Valid {
		var parentAccountGroup domain.AccountGroup

		parentAccountGroup, err = w.reader.GetAccountGroupByID(ctx, accountGroup.ParentID.Int64)
		if err != nil {
			err = errors.PropagateWithCode(err, errors.GetCode(err), "Error on get parent account group")
			return
		}

		if parentAccountGroup.ParentID.Valid {
			err = errors.PropagateWithCode(goErr.New("invalid parent id"), EcodeParentIDNotValid, "cannot set parent from another account group child")
			return
		}

		accountGroup.ClassID = parentAccountGroup.ClassID
	}

	err = w.db.QueryRowContext(ctx, `
		INSERT INTO account_groups (parent_id, class_id, name, inactive)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, accountGroup.ParentID, accountGroup.ClassID, accountGroup.Name, accountGroup.Inactive,
	).Scan(&accountGroup.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreAccountGroupFailed, "Insert account group failed")
		return
	}

	return
}

func (w *writer) UpdateAccountGroupByID(ctx context.Context, id int64, accountGroup *domain.AccountGroup) (err error) {
	_, err = w.updateAccountGroup(ctx, accountGroup, AccountGroupStatement{ID: id})
	return
}

func (w *writer) updateAccountGroup(ctx context.Context, accountGroup *domain.AccountGroup, where AccountGroupStatement) (result sql.Result, err error) {
	if accountGroup.ParentID.Valid && accountGroup.ID != 0 {
		if accountGroup.ParentID.Int64 == accountGroup.ID {
			err = errors.PropagateWithCode(goErr.New("invalid parent id"), EcodeParentIDNotValid, "cannot set parent with the same account group")
			return
		}
	}

	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountGroupFailed, "Failed on select account group")
		return
	}

	updateQuery := fmt.Sprintf(`
		UPDATE account_groups
		SET parent_id = ?, class_id = ?, name = ?, inactive = ?
		%s
	`, whereClause)

	args := append(
		[]interface{}{accountGroup.ParentID, accountGroup.ClassID, accountGroup.Name, accountGroup.Inactive},
		whereClauseArgs...,
	)

	result, err = w.db.ExecContext(ctx, w.db.Rebind(updateQuery), args...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountGroupFailed, "Update account group failed")
		return
	}

	return
}

func (w *writer) DeleteAccountGroupByID(ctx context.Context, id int64) (err error) {
	_, err = w.deleteAccountGroup(ctx, AccountGroupStatement{ID: id})
	return
}

func (w *writer) deleteAccountGroup(ctx context.Context, where AccountGroupStatement) (result sql.Result, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeDeleteAccountGroupFailed, "Failed on select account class")
		return
	}

	deleteQuery := fmt.Sprintf(`
		DELETE FROM account_classes
		%s
	`, whereClause)

	result, err = w.db.ExecContext(ctx, w.db.Rebind(deleteQuery), whereClauseArgs...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeDeleteAccountGroupFailed, "Update account class failed")
		return
	}

	return
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
		SET name = ?, type_id = ?
		%s
	`, whereClause)

	args := append([]interface{}{accountClass.Name, accountClass.TypeID}, whereClauseArgs...)
	result, err = w.db.ExecContext(ctx, w.db.Rebind(updateQuery), args...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountClassFailed, "Update account class failed")
		return
	}

	return
}

func (w *writer) StoreAccountClass(ctx context.Context, accountClass *domain.AccountClass) (err error) {
	classType := classTypes[accountClass.TypeID]
	if classType.ID == 0 {
		err = errors.PropagateWithCode(err, EcodeStoreAccountClassFailed, "Type not valid")
		return
	}

	err = w.db.QueryRowContext(ctx, `
		INSERT INTO account_classes (name, type_id, inactive)
		VALUES ($1, $2, $3)
		RETURNING id
	`, accountClass.Name, accountClass.TypeID, accountClass.Inactive).Scan(&accountClass.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreAccountClassFailed, "Insert account class failed")
		return
	}

	return
}

func NewWriter(opt *Options, reader Reader) Writer {
	return &writer{opt.Logger, opt.MasterDB, reader}
}
