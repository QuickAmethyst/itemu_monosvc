package usecase

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
)

type Reader interface {
	GetAccountClassList(ctx context.Context, stmt sql.AccountClassStatement, p qb.Paging) (result []domain.AccountClass, paging qb.Paging, err error)
	GetAccountClass(ctx context.Context, stmt sql.AccountClassStatement) (accountClass domain.AccountClass, err error)
	GetAccountClassByID(ctx context.Context, id int64) (accountClass domain.AccountClass, err error)
}

type reader struct {
	AccountingSQL sql.SQL
}

func (r *reader) GetAccountClassList(ctx context.Context, stmt sql.AccountClassStatement, p qb.Paging) (result []domain.AccountClass, paging qb.Paging, err error) {
	return r.AccountingSQL.GetAccountClassList(ctx, stmt, p)
}

func (r *reader) GetAccountClass(ctx context.Context, stmt sql.AccountClassStatement) (accountClass domain.AccountClass, err error) {
	return r.AccountingSQL.GetAccountClass(ctx, stmt)
}

func (r *reader) GetAccountClassByID(ctx context.Context, id int64) (accountClass domain.AccountClass, err error) {
	return r.AccountingSQL.GetAccountClassByID(ctx, id)
}

func NewReader(opt *Options) Reader {
	return &reader{opt.AccountingSQL}
}