package usecase

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
)

type Reader interface {
	GetAllAccountClasses(ctx context.Context, stmt sql.AccountClassStatement) (result []domain.AccountClass, err error)
	GetAccountClass(ctx context.Context, stmt sql.AccountClassStatement) (accountClass domain.AccountClass, err error)
	GetAccountClassByID(ctx context.Context, id int64) (accountClass domain.AccountClass, err error)

	GetAllAccountTypes(ctx context.Context) (result []domain.AccountClassType)
	GetAccountClassTypeByID(ctx context.Context, id int64) (accountClassType domain.AccountClassType)

	GetAllAccountGroups(ctx context.Context, stmt sql.AccountGroupStatement) (result []domain.AccountGroup, err error)
	GetAllTopLevelAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (result []domain.AccountGroup, err error)
	GetAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (accountGroup domain.AccountGroup, err error)
	GetAccountGroupByID(ctx context.Context, id int64) (accountGroup domain.AccountGroup, err error)

	GetAllAccounts(ctx context.Context, stmt sql.AccountStatement) (result []domain.Account, err error)
	GetAccount(ctx context.Context, stmt sql.AccountStatement) (account domain.Account, err error)
	GetAccountByID(ctx context.Context, id int64) (account domain.Account, err error)

	GetAllGeneralLedgerPreferences(ctx context.Context, stmt sql.GeneralLedgerPreferenceStatement) (preferences []domain.GeneralLedgerPreference, err error)

	GetFiscalYearList(ctx context.Context, stmt sql.FiscalYearStatement, p qb.Paging) (result []domain.FiscalYear, paging qb.Paging, err error)

	GetAllBankAccountTypes(ctx context.Context) (bankAccountTypes []domain.BankAccountType)
	GetBankAccountList(ctx context.Context, stmt sql.BankAccountStatement, p qb.Paging) (result []domain.BankAccount, paging qb.Paging, err error)
	GetBankAccount(ctx context.Context, stmt sql.BankAccountStatement) (bankAccount domain.BankAccount, err error)
}

type reader struct {
	AccountingSQL sql.SQL
}

func (r *reader) GetAllBankAccountTypes(ctx context.Context) (bankAccountTypes []domain.BankAccountType) {
	return r.AccountingSQL.GetAllBankAccountTypes(ctx)
}

func (r *reader) GetBankAccountList(ctx context.Context, stmt sql.BankAccountStatement, p qb.Paging) (result []domain.BankAccount, paging qb.Paging, err error) {
	return r.AccountingSQL.GetBankAccountList(ctx, stmt, p)
}

func (r *reader) GetBankAccount(ctx context.Context, stmt sql.BankAccountStatement) (bankAccount domain.BankAccount, err error) {
	return r.AccountingSQL.GetBankAccount(ctx, stmt)
}

func (r *reader) GetFiscalYearList(ctx context.Context, stmt sql.FiscalYearStatement, p qb.Paging) (result []domain.FiscalYear, paging qb.Paging, err error) {
	return r.AccountingSQL.GetFiscalYearList(ctx, stmt, p)
}

func (r *reader) GetAllGeneralLedgerPreferences(ctx context.Context, stmt sql.GeneralLedgerPreferenceStatement) (preferences []domain.GeneralLedgerPreference, err error) {
	return r.AccountingSQL.GetAllGeneralLedgerPreferences(ctx, stmt)
}

func (r *reader) GetAllAccounts(ctx context.Context, stmt sql.AccountStatement) (result []domain.Account, err error) {
	return r.AccountingSQL.GetAllAccounts(ctx, stmt)
}

func (r *reader) GetAccount(ctx context.Context, stmt sql.AccountStatement) (account domain.Account, err error) {
	return r.AccountingSQL.GetAccount(ctx, stmt)
}

func (r *reader) GetAccountByID(ctx context.Context, id int64) (account domain.Account, err error) {
	return r.AccountingSQL.GetAccountByID(ctx, id)
}

func (r *reader) GetAllTopLevelAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (result []domain.AccountGroup, err error) {
	stmt.ParentIDIsNULL = true
	return r.GetAllAccountGroups(ctx, stmt)
}

func (r *reader) GetAllAccountGroups(ctx context.Context, stmt sql.AccountGroupStatement) (result []domain.AccountGroup, err error) {
	return r.AccountingSQL.GetAllAccountGroups(ctx, stmt)
}

func (r *reader) GetAccountGroup(ctx context.Context, stmt sql.AccountGroupStatement) (accountGroup domain.AccountGroup, err error) {
	return r.AccountingSQL.GetAccountGroup(ctx, stmt)
}

func (r *reader) GetAccountGroupByID(ctx context.Context, id int64) (accountGroup domain.AccountGroup, err error) {
	return r.AccountingSQL.GetAccountGroupByID(ctx, id)
}

func (r *reader) GetAllAccountTypes(ctx context.Context) (result []domain.AccountClassType) {
	return r.AccountingSQL.GetAllAccountTypes(ctx)
}

func (r *reader) GetAccountClassTypeByID(ctx context.Context, id int64) (accountClassType domain.AccountClassType) {
	return r.AccountingSQL.GetAccountClassTypeByID(ctx, id)
}

func (r *reader) GetAllAccountClasses(ctx context.Context, stmt sql.AccountClassStatement) (result []domain.AccountClass, err error) {
	return r.AccountingSQL.GetAllAccountClasses(ctx, stmt)
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