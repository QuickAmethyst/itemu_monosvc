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
	GetAllAccountClasses(ctx context.Context, stmt AccountClassStatement) (result []domain.AccountClass, err error)
	GetAccountClass(ctx context.Context, stmt AccountClassStatement) (accountClass domain.AccountClass, err error)
	GetAccountClassByID(ctx context.Context, id int64) (accountClass domain.AccountClass, err error)

	GetAllAccountTypes(ctx context.Context) (result []domain.AccountClassType)
	GetAccountClassTypeByID(ctx context.Context, id int64) (accountClassType domain.AccountClassType)

	GetAllAccountGroups(ctx context.Context, stmt AccountGroupStatement) (result []domain.AccountGroup, err error)
	GetAccountGroup(ctx context.Context, stmt AccountGroupStatement) (accountGroup domain.AccountGroup, err error)
	GetAccountGroupByID(ctx context.Context, id int64) (accountGroup domain.AccountGroup, err error)

	GetAllAccounts(ctx context.Context, stmt AccountStatement) (result []domain.Account, err error)
	GetAccount(ctx context.Context, stmt AccountStatement) (account domain.Account, err error)
	GetAccountByID(ctx context.Context, id int64) (account domain.Account, err error)

	ValidatePreferences(ctx context.Context, preferences []domain.GeneralLedgerPreference) (err error)
	GetAllGeneralLedgerPreferences(ctx context.Context, stmt GeneralLedgerPreferenceStatement) (preferences []domain.GeneralLedgerPreference, err error)

	GetFiscalYearList(ctx context.Context, stmt FiscalYearStatement, p qb.Paging) (result []domain.FiscalYear, paging qb.Paging, err error)
	GetFiscalYear(ctx context.Context, stmt FiscalYearStatement) (fiscalYear domain.FiscalYear, err error)
	GetActiveFiscalYear(ctx context.Context) (fiscalYear domain.FiscalYear, err error)
}

type reader struct {
	db sql.DB
}

func (r *reader) GetFiscalYear(ctx context.Context, statement FiscalYearStatement) (fiscalYear domain.FiscalYear, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(statement)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetFiscalYearFailed, "Failed on get fiscal year")
		return
	}

	query := fmt.Sprintf("SELECT id, start_date, end_date, closed FROM fiscal_years %s", whereClause)
	if err = r.db.GetContext(ctx, &fiscalYear, r.db.Rebind(query), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetFiscalYearFailed, "Failed on get fiscal year failed")
		return
	}

	return
}

func (r *reader) GetActiveFiscalYear(ctx context.Context) (fiscalYear domain.FiscalYear, err error) {
	fiscalYear, err = r.GetFiscalYear(ctx, FiscalYearStatement{ClosedNotEQ: true})
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetActiveFiscalYearFailed, "Failed on get active fiscal year")
		return
	}

	return
}

func (r *reader) GetFiscalYearList(ctx context.Context, stmt FiscalYearStatement, p qb.Paging) (result []domain.FiscalYear, paging qb.Paging, err error) {
	result = make([]domain.FiscalYear, 0)
	paging = p
	paging.Normalize()

	fromClause := "FROM fiscal_years"
	limitClause, limitClauseArgs := paging.BuildQuery()
	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetFiscalYearListFailed, "Failed on select fiscal year")
		return
	}
	fmt.Println("HELO", fromClause, whereClause, limitClause, append(whereClauseArgs, limitClauseArgs...))
	selectQuery := fmt.Sprintf("SELECT id, start_date, end_date, closed %s %s %s", fromClause, whereClause, limitClause)
	countQuery := fmt.Sprintf("SELECT COUNT(*) %s %s", fromClause, whereClause)

	if err = r.db.SelectContext(ctx, &result, r.db.Rebind(selectQuery), append(whereClauseArgs, limitClauseArgs...)...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetFiscalYearListFailed, "Failed on select fiscal year")
		return
	}

	if err = r.db.GetContext(ctx, &paging.Total, r.db.Rebind(countQuery), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetFiscalYearListFailed, "Failed on select count fiscal year")
		return
	}

	return
}

func (r *reader) GetAllGeneralLedgerPreferences(ctx context.Context, stmt GeneralLedgerPreferenceStatement) (preferences []domain.GeneralLedgerPreference, err error) {
	preferences = make([]domain.GeneralLedgerPreference, 0)
	fromClause := "FROM general_ledger_preferences"
	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAllGeneralLedgerPreferencesFailed, "Failed on get all general ledger preferences")
		return
	}

	query := fmt.Sprintf("SELECT id, account_id %s %s", fromClause, whereClause)
	if err = r.db.SelectContext(ctx, &preferences, r.db.Rebind(query), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAllGeneralLedgerPreferencesFailed, "Failed on get all general ledger preferences")
		return
	}

	return
}

func (r *reader) GetAllAccounts(ctx context.Context, stmt AccountStatement) (result []domain.Account, err error) {
	result = make([]domain.Account, 0)
	fromClause := "FROM accounts"
	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err , EcodeGetAllAccountsFailed, "Failed on get all accounts")
		return
	}

	query := fmt.Sprintf("SELECT id, name, group_id, inactive %s %s", fromClause, whereClause)
	if err = r.db.SelectContext(ctx, &result, r.db.Rebind(query), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAllAccountGroupsFailed, "Failed on get all accounts")
		return
	}

	return
}

func (r *reader) GetAccount(ctx context.Context, stmt AccountStatement) (account domain.Account, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountFailed, "Failed on get account failed")
		return
	}

	query := fmt.Sprintf("SELECT id, name, group_id, inactive FROM accounts %s", whereClause)
	if err = r.db.GetContext(ctx, &account, r.db.Rebind(query), whereClauseArgs...); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAccountFailed, "Failed on get account failed")
		return
	}

	return
}

func (r *reader) GetAccountByID(ctx context.Context, id int64) (account domain.Account, err error) {
	return r.GetAccount(ctx, AccountStatement{ID: id})
}

func (r *reader) GetAllAccountTypes(ctx context.Context) (result []domain.AccountClassType) {
	result = make([]domain.AccountClassType, len(classTypes))
	for id, classType := range classTypes {
		result[id-1] = classType
	}

	return
}

func (r *reader) GetAccountClassTypeByID(ctx context.Context, id int64) (accountClassType domain.AccountClassType) {
	return classTypes[id]
}

func (r *reader) GetAllAccountGroups(ctx context.Context, stmt AccountGroupStatement) (result []domain.AccountGroup, err error) {
	result = make([]domain.AccountGroup, 0)
	fromClause := "FROM account_groups"

	whereClause, whereClauseArgs, err := qb.NewWhereClause(stmt)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetAllAccountGroupsFailed, "Failed on build where clause")
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

func (r *reader) GetAllAccountClasses(ctx context.Context, stmt AccountClassStatement) (result []domain.AccountClass, err error) {
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


func (r *reader) ValidatePreferences(ctx context.Context, preferences []domain.GeneralLedgerPreference) (err error) {
	var fieldErrors []errors.FieldError

	query := `
		SELECT account_classes.id, account_classes.name, account_classes.type_id, account_classes.inactive
		FROM account_classes, account_groups, accounts
		WHERE 
			accounts.group_id = account_groups.id AND account_groups.class_id = account_classes.id AND
			accounts.id = ?
	`
	for _, preference := range preferences {
		var accountClass domain.AccountClass
		field := fmt.Sprintf("%d", preference.ID)

		if preference.AccountID.Valid && preference.AccountID.Int64 != 0 {
			err = r.db.GetContext(ctx, &accountClass, query, preference.AccountID)
			if err == sql.ErrNoRows {
				fieldErrors = append(fieldErrors, errors.FieldError{Field: field, Message: "Account not found"})
			} else if err != nil {
				fieldErrors = append(fieldErrors, errors.FieldError{Field: field, Message: err.Error()})
				continue
			}
		}

		if preference.ID == int64(RetainedEarnings) && !IsBalanceSheetAccount(accountClass.TypeID) {
			fieldErrors = append(fieldErrors, errors.FieldError{Field: field, Message: "Account must be one of the balance sheet account"})
			continue
		}
	}

	if len(fieldErrors) == 0 {
		return nil
	}

	return errors.ValidationErrors(fieldErrors)
}

func NewReader(opt *Options) Reader {
	return &reader{db: opt.SlaveDB}
}
