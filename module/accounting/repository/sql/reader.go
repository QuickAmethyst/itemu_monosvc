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
	GetAccountClassList(ctx context.Context, stmt AccountClassStatement, p qb.Paging) (result []domain.AccountClass, paging qb.Paging, err error)
	GetAccountClass(ctx context.Context, stmt AccountClassStatement) (accountClass domain.AccountClass, err error)
	GetAccountClassByID(ctx context.Context, id int64) (accountClass domain.AccountClass, err error)
}

type reader struct {
	db sql.DB
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
		SELECT id, name, description, decimal
		FROM account_classes
		%s
	`, whereClause)

	if err = r.db.GetContext(ctx, &accountClass, r.db.Rebind(selectQuery), whereClauseArgs); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountClassFailed, "Failed on get account class")
		return
	}

	return
}

func (r *reader) GetAccountClassList(ctx context.Context, stmt AccountClassStatement, p qb.Paging) (result []domain.AccountClass, paging qb.Paging, err error) {
	result = make([]domain.AccountClass, 0)
	paging = p
	paging.Normalize()

	fromClause := "FROM account_classes"
	limitClause, limitClauseArgs := p.BuildQuery()
	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountClassListFailed, "Failed on build where clause")
		return
	}

	selectQuery := fmt.Sprintf("SELECT id, name, type_id, inactive %s %s %s", fromClause, whereClause, limitClause)
	countQuery := fmt.Sprintf("SELECT COUNT(*) %s %s", fromClause, whereClause)

	if err = r.db.SelectContext(ctx, &result, r.db.Rebind(selectQuery), append(whereClauseArgs, limitClauseArgs...)...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountClassListFailed, "Failed on select account class")
		return
	}

	if err = r.db.GetContext(ctx, &paging.Total, r.db.Rebind(countQuery), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountClassListCountFailed, "Failed on select count account class")
		return
	}

	return
}

func NewReader(opt *Options) Reader {
	return &reader{opt.SlaveDB}
}