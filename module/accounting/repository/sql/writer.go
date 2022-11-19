package sql

import (
	"context"
	goErr "errors"
	"fmt"
	"github.com/QuickAmethyst/monosvc/module/accounting/domain"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
	"github.com/google/uuid"
	"time"
)

type Writer interface {
	StoreAccountClass(ctx context.Context, accountClass *domain.AccountClass) (err error)
	UpdateAccountClassByID(ctx context.Context, id int64, accountClass *domain.AccountClass) (err error)
	DeleteAccountClassByID(ctx context.Context, id int64) (err error)

	StoreAccountGroup(ctx context.Context, accountClassGroup *domain.AccountGroup) (err error)
	UpdateAccountGroupByID(ctx context.Context, id int64, accountGroup *domain.AccountGroup) (err error)
	DeleteAccountGroupByID(ctx context.Context, id int64) (err error)

	StoreAccount(ctx context.Context, account *domain.Account) (err error)
	UpdateAccountByID(ctx context.Context, id int64, account *domain.Account) (err error)
	DeleteAccountByID(ctx context.Context, id int64) (err error)

	StoreTransactions(ctx context.Context, userID uuid.UUID, transactions []Transaction) (journal *domain.Journal, err error)

	UpdateGeneralLedgerPreferenceByID(ctx context.Context, id int64, preference *domain.GeneralLedgerPreference) (err error)
	UpdateGeneralLedgerPreferences(ctx context.Context, preferences []domain.GeneralLedgerPreference) (err error)

	StoreFiscalYear(ctx context.Context, fiscalYear *domain.FiscalYear) (err error)
	UpdateFiscalYearByID(ctx context.Context, id int64, fiscalYear *domain.FiscalYear) (err error)
}

type writer struct {
	logger logger.Logger
	db     sql.DB
	reader Reader
}

func (w *writer) UpdateFiscalYearByID(ctx context.Context, id int64, fiscalYear *domain.FiscalYear) (err error) {
	if _, err = w.db.Updates(ctx, "fiscal_years", fiscalYear, &FiscalYearStatement{ID: id}); err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateFiscalYearFailed, "Update fiscal year failed")
		return
	}

	return
}

func (w *writer) StoreFiscalYear(ctx context.Context, fiscalYear *domain.FiscalYear) (err error) {
	var intersectedFiscalYear domain.FiscalYear

	if fiscalYear.EndDate.Before(fiscalYear.StartDate) {
		err = errors.PropagateWithCode(fmt.Errorf("invalid fiscal year date"), EcodeValidateFiscalYearFailed, "Fiscal year end date must after start date")
		return
	}

	args := []interface{}{fiscalYear.StartDate, fiscalYear.StartDate, fiscalYear.EndDate, fiscalYear.EndDate}
	query := `
		SELECT id, start_date, end_date, closed
		FROM fiscal_years
		WHERE (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?)
	`

	err = w.db.GetContext(ctx, &intersectedFiscalYear, w.db.Rebind(query), args...)
	if err != nil && err != sql.ErrNoRows {
		err = errors.PropagateWithCode(err, EcodeStoreFiscalYearFailed, "Store fiscal year failed")
		return
	}

	if err == nil {
		err = errors.PropagateWithCode(fmt.Errorf("intersect with other fiscal year"), EcodeStoreFiscalYearFailed, "Fiscal year intersect with other fiscal year")
		return
	}

	query = "INSERT INTO fiscal_years (start_date, end_date, closed) VALUES ($1, $2, $3) RETURNING id"
	err = w.db.QueryRowContext(ctx, w.db.Rebind(query), fiscalYear.StartDate, fiscalYear.EndDate, fiscalYear.Closed).Scan(&fiscalYear.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreFiscalYearFailed, "Store fiscal year failed")
		return
	}

	return
}

func (w *writer) UpdateGeneralLedgerPreferences(ctx context.Context, preferences []domain.GeneralLedgerPreference) (err error) {
	if err = w.reader.ValidatePreferences(ctx, preferences); err != nil {
		err = errors.PropagateWithCode(err, EcodeValidatePreferencesFailed, "Update general ledger preferences failed")
		return
	}

	err = w.db.Transaction(ctx, nil, func(tx *sql.Tx) error {
		for _, preference := range preferences {
			err = w.mustUpdateGeneralLedgerPreferenceByID(ctx, preference.ID, &preference)
			if err != nil {
				err = errors.PropagateWithCode(err, EcodeUpdateGeneralLedgerPreferenceFailed, "Update general ledger preferences failed")
				return err
			}
		}

		return nil
	})

	return
}

func (w *writer) UpdateGeneralLedgerPreferenceByID(ctx context.Context, id int64, preference *domain.GeneralLedgerPreference) (err error) {
	err = w.reader.ValidatePreferences(ctx, []domain.GeneralLedgerPreference{
		{ID: id, AccountID: preference.AccountID},
	})

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeValidatePreferencesFailed, "Update general ledger preferences failed")
		return
	}

	if err = w.mustUpdateGeneralLedgerPreferenceByID(ctx, id, preference); err != nil {
		err = errors.PropagateWithCode(goErr.New("update general ledger preference by id failed"), EcodeUpdateGeneralLedgerPreferenceFailed, "creator unknown")
		return
	}

	return
}

func (w *writer) mustUpdateGeneralLedgerPreferenceByID(ctx context.Context, id int64, preference *domain.GeneralLedgerPreference) (err error) {
	if _, err = w.db.Updates(ctx, "general_ledger_preferences", preference, &GeneralLedgerPreferenceStatement{ID: id}); err != nil {
		err = errors.PropagateWithCode(goErr.New("update general ledger preference by id failed"), EcodeUpdateGeneralLedgerPreferenceFailed, "creator unknown")
		return
	}

	return
}

func (w *writer) StoreTransactions(ctx context.Context, userID uuid.UUID, transactions []Transaction) (journal *domain.Journal, err error) {
	var (
		gls           []domain.GeneralLedger
		journalAmount float64
		balanceAmount float64
	)

	if userID == uuid.Nil {
		err = errors.PropagateWithCode(goErr.New("creator unknown"), EcodeStoreTransactionCreatedByRequired, "creator unknown")
		return
	}

	now := time.Now()
	journalID := uuid.New()

	for _, transaction := range transactions {
		if transaction.Amount == 0 {
			continue
		}

		gls = append(gls, domain.GeneralLedger{
			ID:        uuid.New(),
			JournalID: journalID,
			AccountID: transaction.AccountID,
			CreatedBy: userID,
			Amount:    transaction.Amount,
		})

		balanceAmount += transaction.Amount
		if transaction.Amount > 0 {
			journalAmount += transaction.Amount
		}
	}

	// check mapAccountGeneralLedger is empty.
	// empty is occur when all the transaction amount is zero
	if len(gls) == 0 {
		return nil, nil
	}

	if balanceAmount != 0 {
		err = errors.PropagateWithCode(fmt.Errorf("transaction not balance"), EcodeTransactionNotBalance, "Transaction not balance")
		return
	}

	journal = &domain.Journal{
		ID:        journalID,
		Amount:    journalAmount,
		CreatedAt: now,
	}

	err = w.db.Transaction(ctx, nil, func(tx *sql.Tx) error {
		err = w.StoreJournalTx(tx, ctx, journal)

		if err != nil {
			err = errors.PropagateWithCode(err, EcodeStoreTransactionAtJournalFailed, "Store transaction failed")
			return err
		}

		err = w.StoreGeneralLedgersTx(tx, ctx, gls)
		if err != nil {
			err = errors.PropagateWithCode(err, EcodeStoreTransactionAtGeneralLedgerFailed, "Store transaction failed")
			return err
		}

		return nil
	})

	return
}

func (w *writer) StoreJournalTx(tx *sql.Tx, ctx context.Context, journal *domain.Journal) (err error) {
	query := "INSERT INTO journals (id, amount, created_at) VALUES ($1, $2, $3) RETURNING id"
	err = tx.QueryRowContext(ctx, query, journal.ID, journal.Amount, journal.CreatedAt).Scan(&journal.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreJournalFailed, "Store journal failed")
		return
	}

	return
}

func (w *writer) StoreGeneralLedgersTx(tx *sql.Tx, ctx context.Context, gls []domain.GeneralLedger) (err error) {
	var params []interface{}

	query := `Insert INTO general_ledgers (id, journal_id, account_id, amount, created_by) VALUES `

	for _, gl := range gls {
		query += "(?,?,?,?,?),"
		params = append(params, gl.ID, gl.JournalID, gl.AccountID, gl.Amount, gl.CreatedBy)
	}

	query = query[:len(query)-1] // remove trailing ","
	_, err = tx.ExecContext(ctx, w.db.Rebind(query), params...)
	return
}

func (w *writer) StoreAccount(ctx context.Context, account *domain.Account) (err error) {
	err = w.db.QueryRowContext(ctx, `
		INSERT INTO accounts (name, group_id, inactive)
		VALUES ($1, $2, $3)
		RETURNING id
	`, account.Name, account.GroupID, account.Inactive).Scan(&account.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreAccountFailed, "Store account failed")
		return
	}

	return
}

func (w *writer) UpdateAccountByID(ctx context.Context, id int64, account *domain.Account) (err error) {
	dest := map[string]interface{}{
		"name":     account.Name,
		"group_id": account.GroupID,
		"inactive": account.Inactive,
	}

	if _, err = w.db.Updates(ctx, "accounts", dest, &AccountStatement{ID: id}); err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountFailed, "Update account by id failed")
		return
	}

	return
}

func (w *writer) DeleteAccountByID(ctx context.Context, id int64) (err error) {
	if _, err = w.db.Delete(ctx, "accounts", &AccountStatement{ID: id}); err != nil {
		err = errors.PropagateWithCode(err, EcodeDeleteAccountFailed, "Delete account by id failed")
		return
	}

	return
}

func (w *writer) StoreAccountGroup(ctx context.Context, accountGroup *domain.AccountGroup) (err error) {
	if accountGroup.ParentID.Valid && accountGroup.ParentID.Int64 != 0 {
		var parentAccountGroup domain.AccountGroup

		parentAccountGroup, err = w.reader.GetAccountGroupByID(ctx, accountGroup.ParentID.Int64)
		if err != nil {
			err = errors.PropagateWithCode(err, errors.GetCode(err), "Error on get parent account group")
			return
		}

		if parentAccountGroup.ParentID.Valid {
			err = errors.PropagateWithCode(goErr.New("invalid parent id"), EcodeParentIDNotValid, "cannot set parent from another account group child")
			return
		}

		accountGroup.ClassID = parentAccountGroup.ClassID
	}

	err = w.db.QueryRowContext(ctx, `
		INSERT INTO account_groups (parent_id, class_id, name, inactive)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, accountGroup.ParentID, accountGroup.ClassID, accountGroup.Name, accountGroup.Inactive,
	).Scan(&accountGroup.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreAccountGroupFailed, "Insert account group failed")
		return
	}

	return
}

func (w *writer) UpdateAccountGroupByID(ctx context.Context, id int64, accountGroup *domain.AccountGroup) (err error) {
	_, err = w.updateAccountGroup(ctx, accountGroup, AccountGroupStatement{ID: id})
	return
}

func (w *writer) updateAccountGroup(ctx context.Context, accountGroup *domain.AccountGroup, where AccountGroupStatement) (result sql.Result, err error) {
	if accountGroup.ParentID.Valid && accountGroup.ID != 0 {
		if accountGroup.ParentID.Int64 == accountGroup.ID {
			err = errors.PropagateWithCode(goErr.New("invalid parent id"), EcodeParentIDNotValid, "cannot set parent with the same account group")
			return
		}
	}

	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountGroupFailed, "Failed on select account group")
		return
	}

	updateQuery := fmt.Sprintf(`
		UPDATE account_groups
		SET parent_id = ?, class_id = ?, name = ?, inactive = ?
		%s
	`, whereClause)

	args := append(
		[]interface{}{accountGroup.ParentID, accountGroup.ClassID, accountGroup.Name, accountGroup.Inactive},
		whereClauseArgs...,
	)

	result, err = w.db.ExecContext(ctx, w.db.Rebind(updateQuery), args...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountGroupFailed, "Update account group failed")
		return
	}

	return
}

func (w *writer) DeleteAccountGroupByID(ctx context.Context, id int64) (err error) {
	_, err = w.deleteAccountGroup(ctx, AccountGroupStatement{ID: id})
	return
}

func (w *writer) deleteAccountGroup(ctx context.Context, where AccountGroupStatement) (result sql.Result, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeDeleteAccountGroupFailed, "Failed on select account class")
		return
	}

	deleteQuery := fmt.Sprintf(`
		DELETE FROM account_classes
		%s
	`, whereClause)

	result, err = w.db.ExecContext(ctx, w.db.Rebind(deleteQuery), whereClauseArgs...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeDeleteAccountGroupFailed, "Update account class failed")
		return
	}

	return
}

func (w *writer) DeleteAccountClassByID(ctx context.Context, id int64) (err error) {
	_, err = w.deleteAccountClass(ctx, AccountClassStatement{ID: id})
	return
}

func (w *writer) deleteAccountClass(ctx context.Context, where AccountClassStatement) (result sql.Result, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountClassFailed, "Failed on delete account class")
		return
	}

	deleteQuery := fmt.Sprintf(`
		DELETE FROM account_classes
		%s
	`, whereClause)

	result, err = w.db.ExecContext(ctx, w.db.Rebind(deleteQuery), whereClauseArgs...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeDeleteAccountClassFailed, "Failed on delete account class")
		return
	}

	return
}

func (w *writer) UpdateAccountClassByID(ctx context.Context, id int64, accountClass *domain.AccountClass) (err error) {
	_, err = w.updateAccountClass(ctx, accountClass, AccountClassStatement{ID: id})
	return
}

func (w *writer) updateAccountClass(ctx context.Context, accountClass *domain.AccountClass, where AccountClassStatement) (result sql.Result, err error) {
	whereClause, whereClauseArgs, err := qb.NewWhereClause(where)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountClassFailed, "Failed on select account class")
		return
	}

	updateQuery := fmt.Sprintf(`
		UPDATE account_classes
		SET name = ?, type_id = ?
		%s
	`, whereClause)

	args := append([]interface{}{accountClass.Name, accountClass.TypeID}, whereClauseArgs...)
	result, err = w.db.ExecContext(ctx, w.db.Rebind(updateQuery), args...)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateAccountClassFailed, "Update account class failed")
		return
	}

	return
}

func (w *writer) StoreAccountClass(ctx context.Context, accountClass *domain.AccountClass) (err error) {
	classType := classTypes[accountClass.TypeID]
	if classType.ID == 0 {
		err = errors.PropagateWithCode(err, EcodeStoreAccountClassFailed, "Type not valid")
		return
	}

	err = w.db.QueryRowContext(ctx, w.db.Rebind(`
		INSERT INTO account_classes (name, type_id, inactive)
		VALUES (?, ?, ?)
		RETURNING id
	`), accountClass.Name, accountClass.TypeID, accountClass.Inactive).Scan(&accountClass.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreAccountClassFailed, "Insert account class failed")
		return
	}

	return
}

func NewWriter(opt *Options, reader Reader) Writer {
	return &writer{opt.Logger, opt.MasterDB, reader}
}
