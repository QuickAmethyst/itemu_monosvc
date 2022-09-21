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
	GetAccountClassTypeList(ctx context.Context) (result []domain.AccountClassType)
	GetAccountClassTypeByID(ctx context.Context, id int64) (accountClassType domain.AccountClassType)
}

type reader struct {
	AccountingSQL sql.SQL
}

func (r *reader) GetAccountClassTypeList(_ context.Context) (result []domain.AccountClassType) {
	result = make([]domain.AccountClassType, len(classTypes))
	for id, classType := range classTypes {
		result[id - 1] = classType
	}
	return
}

func (r *reader) GetAccountClassTypeByID(_ context.Context, id int64) (accountClassType domain.AccountClassType) {
	return classTypes[id]
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