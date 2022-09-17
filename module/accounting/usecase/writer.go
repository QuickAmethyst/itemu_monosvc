package usecase

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"
)

type Writer interface {
	StoreAccountClass(ctx context.Context, accountClasses *domain.AccountClass) (err error)
	UpdateAccountClassByID(ctx context.Context, id int64, accountClass *domain.AccountClass) (err error)
}

type writer struct {
	AccountingSQL sql.SQL
}

func (w *writer) UpdateAccountClassByID(ctx context.Context, id int64, accountClass *domain.AccountClass) (err error) {
	return w.AccountingSQL.UpdateAccountClassByID(ctx, id, accountClass)
}

func (w *writer) StoreAccountClass(ctx context.Context, accountClasses *domain.AccountClass) (err error) {
	return w.AccountingSQL.StoreAccountClass(ctx, accountClasses)
}

func NewWriter(opt *Options) Writer {
	return &writer{opt.AccountingSQL}
}
