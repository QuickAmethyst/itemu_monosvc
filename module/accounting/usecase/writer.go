package usecase

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"
)

type Writer interface {
	StoreAccountClass(ctx context.Context, accountClasses *domain.AccountClass) (err error)
	UpdateAccountClassByID(ctx context.Context, id int64, accountClass *domain.AccountClass) (err error)
	DeleteAccountClassByID(ctx context.Context, id int64) (err error)

	StoreAccountGroup(ctx context.Context, accountGroup *domain.AccountGroup) (err error)
	UpdateAccountGroupByID(ctx context.Context, id int64, accountGroup *domain.AccountGroup) (err error)
	DeleteAccountGroupByID(ctx context.Context, id int64) (err error)
}

type writer struct {
	AccountingSQL sql.SQL
}

func (w *writer) StoreAccountGroup(ctx context.Context, accountClassGroup *domain.AccountGroup) (err error) {
	return w.AccountingSQL.StoreAccountGroup(ctx, accountClassGroup)
}

func (w *writer) UpdateAccountGroupByID(ctx context.Context, id int64, accountGroup *domain.AccountGroup) (err error) {
	return w.AccountingSQL.UpdateAccountGroupByID(ctx, id, accountGroup)
}

func (w *writer) DeleteAccountGroupByID(ctx context.Context, id int64) (err error) {
	return w.AccountingSQL.DeleteAccountGroupByID(ctx, id)
}

func (w *writer) DeleteAccountClassByID(ctx context.Context, id int64) (err error) {
	return w.AccountingSQL.DeleteAccountClassByID(ctx, id)
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
