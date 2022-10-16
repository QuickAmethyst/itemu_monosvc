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

	StoreAccount(ctx context.Context, account *domain.Account) (err error)
	UpdateAccountByID(ctx context.Context, id int64, account *domain.Account) (err error)
	DeleteAccountByID(ctx context.Context, id int64) (err error)

	StoreTransactions(ctx context.Context, transactions []sql.Transaction) (journal domain.Journal, err error)
}

type writer struct {
	AccountingSQL sql.SQL
}

func (w *writer) StoreTransactions(ctx context.Context, transactions []sql.Transaction) (journal domain.Journal, err error) {
	return w.AccountingSQL.StoreTransactions(ctx, transactions)
}

func (w *writer) StoreAccount(ctx context.Context, account *domain.Account) (err error) {
	return w.AccountingSQL.StoreAccount(ctx, account)
}

func (w *writer) UpdateAccountByID(ctx context.Context, id int64, account *domain.Account) (err error) {
	return w.AccountingSQL.UpdateAccountByID(ctx, id, account)
}

func (w *writer) DeleteAccountByID(ctx context.Context, id int64) (err error) {
	return w.AccountingSQL.DeleteAccountByID(ctx, id)
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
