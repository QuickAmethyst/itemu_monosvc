package usecase

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"
)

type Reader interface {
	GetAllAccountClass(ctx context.Context, stmt sql.AccountClassStatement) (result []domain.AccountClass, err error)
	GetAccountClass(ctx context.Context, stmt sql.AccountClassStatement) (accountClass domain.AccountClass, err error)
	GetAccountClassByID(ctx context.Context, id int64) (accountClass domain.AccountClass, err error)

	GetAccountClassTypeList(ctx context.Context) (result []domain.AccountClassType)
	GetAccountClassTypeByID(ctx context.Context, id int64) (accountClassType domain.AccountClassType)

	GetAllAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (result []domain.AccountGroup, err error)
	GetAllTopLevelAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (result []domain.AccountGroup, err error)
	GetAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (accountGroup domain.AccountGroup, err error)
	GetAccountGroupByID(ctx context.Context, id int64) (accountGroup domain.AccountGroup, err error)
}

type reader struct {
	AccountingSQL sql.SQL
}

func (r *reader) GetAllTopLevelAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (result []domain.AccountGroup, err error) {
	stmt.ParentIDIsNULL = true
	return r.GetAllAccountGroup(ctx, stmt)
}

func (r *reader) GetAllAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (result []domain.AccountGroup, err error) {
	return r.AccountingSQL.GetAllAccountGroup(ctx, stmt)
}

func (r *reader) GetAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (accountGroup domain.AccountGroup, err error) {
	return r.AccountingSQL.GetAccountGroup(ctx, stmt)
}

func (r *reader) GetAccountGroupByID(ctx context.Context, id int64) (accountGroup domain.AccountGroup, err error) {
	return r.AccountingSQL.GetAccountGroupByID(ctx, id)
}

func (r *reader) GetAccountClassTypeList(ctx context.Context) (result []domain.AccountClassType) {
	return r.AccountingSQL.GetAccountClassTypeList(ctx)
}

func (r *reader) GetAccountClassTypeByID(ctx context.Context, id int64) (accountClassType domain.AccountClassType) {
	return r.AccountingSQL.GetAccountClassTypeByID(ctx, id)
}

func (r *reader) GetAllAccountClass(ctx context.Context, stmt sql.AccountClassStatement) (result []domain.AccountClass, err error) {
	return r.AccountingSQL.GetAllAccountClass(ctx, stmt)
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