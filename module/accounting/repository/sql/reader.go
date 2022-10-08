package sql

import (
	"context"
	"fmt"
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
)

type Reader interface {
	GetAllAccountClass(ctx context.Context, stmt AccountClassStatement) (result []domain.AccountClass, err error)
	GetAccountClass(ctx context.Context, stmt AccountClassStatement) (accountClass domain.AccountClass, err error)
	GetAccountClassByID(ctx context.Context, id int64) (accountClass domain.AccountClass, err error)

	GetAccountClassTypeList(ctx context.Context) (result []domain.AccountClassType)
	GetAccountClassTypeByID(ctx context.Context, id int64) (accountClassType domain.AccountClassType)

	GetAllAccountGroup(ctx context.Context, stmt AccountGroupStatement) (result []domain.AccountGroup, err error)
	GetAccountGroup(ctx context.Context, stmt AccountGroupStatement) (accountGroup domain.AccountGroup, err error)
	GetAccountGroupByID(ctx context.Context, id int64) (accountGroup domain.AccountGroup, err error)
}

type reader struct {
	db sql.DB
}

func (r *reader) GetAccountClassTypeList(ctx context.Context) (result []domain.AccountClassType) {
	result = make([]domain.AccountClassType, len(classTypes))
	for id, classType := range classTypes {
		result[id-1] = classType
	}

	return
}

func (r *reader) GetAccountClassTypeByID(ctx context.Context, id int64) (accountClassType domain.AccountClassType) {
	return classTypes[id]
}

func (r *reader) GetAllAccountGroup(ctx context.Context, stmt AccountGroupStatement) (result []domain.AccountGroup, err error) {
	result = make([]domain.AccountGroup, 0)
	fromClause := "FROM account_groups"

	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAllAccountGroupFailed, "Failed on build where clause")
		return
	}

	selectQuery := fmt.Sprintf("SELECT id, parent_id, class_id, name, inactive %s %s", fromClause, whereClause)
	if err = r.db.SelectContext(ctx, &result, r.db.Rebind(selectQuery), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAllTopLevelAccountGroupFailed, "Failed on select all account group")
		return
	}

	return
}

func (r *reader) GetAccountGroup(ctx context.Context, stmt AccountGroupStatement) (accountGroup domain.AccountGroup, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountGroupFailed, "Failed on build where clause")
		return
	}

	selectQuery := fmt.Sprintf(`
		SELECT id, parent_id, class_id, name, inactive
		FROM account_groups
		%s
	`, whereClause)

	if err = r.db.GetContext(ctx, &accountGroup, r.db.Rebind(selectQuery), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountGroupFailed, "Failed on get account group")
		return
	}

	return
}

func (r *reader) GetAccountGroupByID(ctx context.Context, id int64) (accountGroup domain.AccountGroup, err error) {
	return r.GetAccountGroup(ctx, AccountGroupStatement{ID: id})
}

func (r *reader) GetAccountClassByID(ctx context.Context, id int64) (accountClass domain.AccountClass, err error) {
	return r.GetAccountClass(ctx, AccountClassStatement{ID: id})
}

func (r *reader) GetAccountClass(ctx context.Context, stmt AccountClassStatement) (accountClass domain.AccountClass, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountClassFailed, "Failed on build where clause")
		return
	}

	selectQuery := fmt.Sprintf(`
		SELECT id, name, type_id, inactive
		FROM account_classes
		%s
	`, whereClause)

	if err = r.db.GetContext(ctx, &accountClass, r.db.Rebind(selectQuery), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountClassFailed, "Failed on get account class")
		return
	}

	return
}

func (r *reader) GetAllAccountClass(ctx context.Context, stmt AccountClassStatement) (result []domain.AccountClass, err error) {
	result = make([]domain.AccountClass, 0)
	fromClause := "FROM account_classes"
	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountClassListFailed, "Failed on build where clause")
		return
	}

	selectQuery := fmt.Sprintf("SELECT id, name, type_id, inactive %s %s", fromClause, whereClause)

	if err = r.db.SelectContext(ctx, &result, r.db.Rebind(selectQuery), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountClassListFailed, "Failed on select account class")
		return
	}

	return
}

func NewReader(opt *Options) Reader {
	return &reader{db: opt.SlaveDB}
}
