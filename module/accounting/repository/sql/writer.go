package sql

import (
	"context"
	goSql "database/sql"
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

	StoreTransaction(ctx context.Context, userID uuid.UUID, transactions Transaction) (journal *domain.Journal, err error)
	StoreTransactionTx(tx sql.Tx, ctx context.Context, userID uuid.UUID, transaction Transaction) (journal *domain.Journal, err error)
	VoidTransactionByID(ctx context.Context, journalID uuid.UUID) (err error)
	VoidTransactionByIDTx(tx sql.Tx, ctx context.Context, journalID uuid.UUID) (err error)

	UpdateGeneralLedgerPreferenceByID(ctx context.Context, id int64, preference *domain.GeneralLedgerPreference) (err error)
	UpdateGeneralLedgerPreferences(ctx context.Context, preferences []domain.GeneralLedgerPreference) (err error)

	StoreBankAccount(ctx context.Context, bankAccount *domain.BankAccount) (err error)
	UpdateBankAccountByID(ctx context.Context, id int64, bankAccount *domain.BankAccount) (err error)
	StoreBankDepositTransaction(ctx context.Context, userID uuid.UUID, transaction BankTransaction) (bankTransaction domain.BankTransaction, err error)

	StoreFiscalYear(ctx context.Context, fiscalYear *domain.FiscalYear) (err error)
	CloseFiscalYear(ctx context.Context, id int64, userID uuid.UUID) (err error)
}

type writer struct {
	logger logger.Logger
	db     sql.DB
	reader Reader
}

func (w *writer) storeBankTransaction(ctx context.Context, userID uuid.UUID, transaction BankTransaction) (bankTransaction domain.BankTransaction, err error) {
	var (
		totalAmount float64
		bankAccount domain.BankAccount
	)

	if transaction.bankTransactionType < Deposit || transaction.bankTransactionType > Deposit {
		err = errors.PropagateWithCode(fmt.Errorf("invalid bank transaction type"), EcodeBankTransactionTypeInvalid, "Invalid bank transaction type")
		return
	}

	if userID == uuid.Nil {
		err = errors.PropagateWithCode(fmt.Errorf("invalid user"), EcodeBankTransactionUserInvalid, "invalid user")
		return
	}

	if transaction.JournalID == uuid.Nil {
		transaction.JournalID = uuid.New()
	}

	if transaction.Date.IsZero() {
		transaction.Date = time.Now()
	}

	for _, row := range transaction.Data {
		var isBankAccount bool

		if row.Amount <= 0 {
			err = errors.PropagateWithCode(fmt.Errorf("amount must be greater than zero"), EcodeBankAccountDepositInvalidAmount, "Invalid bank account deposit amount")
			return
		}

		isBankAccount, err = w.reader.IsBankAccount(ctx, row.AccountID)
		if err != nil {
			return
		}

		if isBankAccount {
			err = errors.PropagateWithCode(err, EcodeBankAccountDepositInvalidAccount, "Bank account is prohibited")
			return
		}

		totalAmount += row.Amount
	}

	bankAccount, err = w.reader.GetBankAccountByID(ctx, transaction.BankAccountID)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetBankAccountFailed, "Failed on get bank account")
		return
	}

	// post bank account to gl
	transaction.Transaction.Data = append(transaction.Transaction.Data, TransactionRow{
		AccountID: bankAccount.AccountID,
		Amount:    totalAmount,
	})

	err = w.db.Transaction(ctx, nil, func(tx sql.Tx) error {
		if _, err := w.StoreTransaction(ctx, userID, transaction.Transaction); err != nil {
			err = errors.PropagateWithCode(err, EcodeStoreTransactionFailed, "Failed on store journal")
			return err
		}

		query := `
			INSERT INTO bank_transactions (journal_id, bank_account_id, amount, balance, memo, created_by, trans_date) VALUES
			(?, ?, ?, (SELECT COALESCE((SELECT balance FROM test_insert GROUP BY id  ORDER BY id DESC LIMIT 1), 0) + ?), ?, ?, ?);
		`

		err = w.db.QueryRowContext(
			ctx,
			w.db.Rebind(query),
			transaction.JournalID,
			transaction.BankAccountID,
			totalAmount,
			totalAmount,
			transaction.Memo,
			userID,
			transaction.Date,
		).Scan(&bankTransaction.ID)
		if err != nil {
			err = errors.PropagateWithCode(err, EcodeStoreTransactionFailed, "Failed on store transaction")
			return err
		}

		return nil
	})

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreBankTransactionFailed, "Failed on store bank transaction")
		return
	}

	return
}

func (w *writer) StoreBankDepositTransaction(ctx context.Context, userID uuid.UUID, transaction BankTransaction) (bankTransaction domain.BankTransaction, err error) {
	transaction.bankTransactionType = Deposit

	if bankTransaction, err = w.storeBankTransaction(ctx, userID, transaction); err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreBankAccountDepositFailed, "Failed on store bank account deposit")
		return
	}

	return
}

func (w *writer) VoidTransactionByID(ctx context.Context, journalID uuid.UUID) (err error) {
	err = w.db.Transaction(ctx, nil, func(tx sql.Tx) error {
		err = w.VoidTransactionByIDTx(tx, ctx, journalID)
		if err != nil {
			err = errors.PropagateWithCode(err, EcodeVoidTransactionByIDFailed, "Failed on void transaction")
			return err
		}

		return nil
	})

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeVoidTransactionByIDFailed, "Failed on void transaction")
		return err
	}

	return
}

func (w *writer) VoidTransactionByIDTx(tx sql.Tx, ctx context.Context, journalID uuid.UUID) (err error) {
	journal, err := w.reader.GetJournalByID(ctx, journalID)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeVoidTransactionByIDFailed, "Failed on get journal by id")
		return
	}

	activeFiscalYear, err := w.reader.GetActiveFiscalYear(ctx)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeVoidTransactionByIDFailed, "Failed on get active fiscal year")
		return
	}

	// check is journal closed
	if journal.TransDate.Before(activeFiscalYear.StartDate) {
		err = errors.PropagateWithCode(err, EcodeJournalAlreadyClosed, "Closing closed journal prohibited")
		return
	}

	query := "UPDATE journals SET deleted_at = ? WHERE journal_id = ?"
	if _, err = tx.ExecContext(ctx, tx.Rebind(query), time.Now(), journalID); err != nil {
		err = errors.PropagateWithCode(err, EcodeVoidTransactionByIDFailed, "Failed on void transaction")
		return
	}

	return
}

func (w *writer) StoreBankAccount(ctx context.Context, bankAccount *domain.BankAccount) (err error) {
	hasTransaction, err := w.reader.AccountHasTransaction(ctx, bankAccount.ID)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreBankAccountFailed, "Store bank account failed")
		return
	}

	if hasTransaction {
		err = errors.PropagateWithCode(fmt.Errorf("account already has transaction"), EcodeAccountHasTransaction, "Store bank account failed")
		return
	}

	query := "INSERT INTO bank_accounts (account_id, type_id, bank_number) VALUES (?, ?, ?) RETURNING id"
	err = w.db.QueryRowContext(
		ctx,
		w.db.Rebind(query),
		bankAccount.AccountID, bankAccount.TypeID, bankAccount.BankNumber,
	).Scan(&bankAccount.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreBankAccountFailed, "Store bank account failed")
		return
	}

	return
}

func (w *writer) UpdateBankAccountByID(ctx context.Context, id int64, bankAccount *domain.BankAccount) (err error) {
	prevBankAccount, err := w.reader.GetBankAccount(ctx, BankAccountStatement{ID: id})
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateBankAccountFailed, "Update bank account by id failed")
		return
	}

	if prevBankAccount.AccountID != bankAccount.AccountID {
		var hasTransaction bool

		hasTransaction, err = w.reader.AccountHasTransaction(ctx, bankAccount.ID)
		if err != nil {
			err = errors.PropagateWithCode(err, EcodeUpdateBankAccountFailed, "Update bank account by id failed")
			return
		}

		if hasTransaction {
			err = errors.PropagateWithCode(fmt.Errorf("account already has transaction"), EcodeAccountHasTransaction, "Store bank account failed")
			return
		}
	}

	bankAccount.ID = id
	if _, err = w.db.Updates(ctx, "bank_accounts", bankAccount, &BankAccountStatement{ID: id}); err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateBankAccountFailed, "Update bank account by id failed")
		return
	}

	return
}

func (w *writer) CloseFiscalYear(ctx context.Context, id int64, userID uuid.UUID) (err error) {
	// check if retained earnings account is valid
	retainedEarningsGLP, err := w.reader.GetGeneralLedgerPreferenceByID(ctx, GeneralLedgerPreferenceStatement{ID: int64(RetainedEarnings)})
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetGeneralLedgerPreferenceFailed, "Failed on get general ledger preference")
		return
	}

	if err = w.reader.ValidatePreferences(ctx, []domain.GeneralLedgerPreference{retainedEarningsGLP}); err != nil {
		err = errors.PropagateWithCode(err, EcodeCloseFiscalYearFailed, "Validate general ledger preferences failed")
		return
	}

	// check if close fiscal year is allowed
	activeFiscalYear, err := w.reader.GetActiveFiscalYear(ctx)
	if err != nil && err != sql.ErrNoRows {
		err = errors.PropagateWithCode(err, EcodeCloseFiscalYearFailed, "Failed on get active fiscal year")
		return
	}

	if err != sql.ErrNoRows && activeFiscalYear.ID != id {
		err = errors.PropagateWithCode(fmt.Errorf("close fiscal year prohibited"), EcodeCloseFiscalYearFailed, "Cannot close fiscal year while there are open fiscal years before")
		return
	}

	// get balance amount
	balance, err := w.reader.GetBalanceSheetAmount(ctx, activeFiscalYear.StartDate, activeFiscalYear.EndDate)
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeGetBalanceSheetAmountFailed, "Failed on get balance sheet amount")
		return
	}

	err = w.db.Transaction(ctx, nil, func(tx sql.Tx) error {
		_, err = w.StoreTransactionTx(tx, ctx, userID, Transaction{
			Memo: "Close Fiscal Year",
			Data: []TransactionRow{
				{retainedEarningsGLP.AccountID.Int64, -balance},
			},
		})

		if err != nil {
			err = errors.PropagateWithCode(err, EcodeStoreTransactionFailed, "Failed on store transaction")
			return err
		}

		activeFiscalYear.Closed = true
		if err = w.updateFiscalYearByIDTx(tx, ctx, id, &activeFiscalYear); err != nil {
			err = errors.PropagateWithCode(err, EcodeCloseFiscalYearFailed, "Update fiscal year failed")
			return err
		}

		return nil
	})

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeCloseFiscalYearFailed, "Failed on close fiscal year")
		return
	}

	return
}

func (w *writer) updateFiscalYearByIDTx(tx sql.Tx, ctx context.Context, id int64, fiscalYear *domain.FiscalYear) (err error) {
	if _, err = tx.Updates(ctx, "fiscal_years", fiscalYear, &FiscalYearStatement{ID: id}); err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateFiscalYearFailed, "Update fiscal year failed")
		return
	}

	return
}

func (w *writer) updateFiscalYearByID(ctx context.Context, id int64, fiscalYear *domain.FiscalYear) (err error) {
	err = w.db.Transaction(ctx, nil, func(tx sql.Tx) error {
		err = w.updateFiscalYearByIDTx(tx, ctx, id, fiscalYear)
		if err != nil {
			err = errors.PropagateWithCode(err, EcodeUpdateFiscalYearFailed, "Update fiscal year failed")
			return err
		}

		return nil
	})

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateFiscalYearFailed, "Update fiscal year failed")
		return err
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

	err = w.db.Transaction(ctx, nil, func(tx sql.Tx) error {
		for _, preference := range preferences {
			err = w.mustUpdateGeneralLedgerPreferenceByIDTx(tx, ctx, preference.ID, &preference)
			if err != nil {
				err = errors.PropagateWithCode(err, EcodeUpdateGeneralLedgerPreferenceFailed, "Update general ledger preferences failed")
				return err
			}
		}

		return nil
	})

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateGeneralLedgerPreferenceFailed, "Failed on update general ledger preferences")
	}

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

func (w *writer) mustUpdateGeneralLedgerPreferenceByIDTx(tx sql.Tx, ctx context.Context, id int64, preference *domain.GeneralLedgerPreference) (err error) {
	if _, err = tx.Updates(ctx, "general_ledger_preferences", preference, &GeneralLedgerPreferenceStatement{ID: id}); err != nil {
		err = errors.PropagateWithCode(goErr.New("update general ledger preference by id failed"), EcodeUpdateGeneralLedgerPreferenceFailed, "creator unknown")
		return
	}

	return
}

func (w *writer) mustUpdateGeneralLedgerPreferenceByID(ctx context.Context, id int64, preference *domain.GeneralLedgerPreference) (err error) {
	err = w.db.Transaction(ctx, nil, func(tx sql.Tx) error {
		err = w.mustUpdateGeneralLedgerPreferenceByIDTx(tx, ctx, id, preference)
		if err != nil {
			err = errors.PropagateWithCode(err, EcodeUpdateGeneralLedgerPreferenceFailed, "Failed on update general ledger preference")
			return err
		}

		return nil
	})

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeUpdateGeneralLedgerPreferenceFailed, "Failed on commit")
		return
	}

	return
}

func (w *writer) StoreTransactionTx(tx sql.Tx, ctx context.Context, userID uuid.UUID, transaction Transaction) (journal *domain.Journal, err error) {
	var (
		gls           []domain.GeneralLedger
		memo          goSql.NullString
		journalAmount float64
		balanceAmount float64
	)

	if userID == uuid.Nil {
		err = errors.PropagateWithCode(goErr.New("creator unknown"), EcodeStoreTransactionCreatedByRequired, "creator unknown")
		return
	}

	if transaction.Memo != "" {
		if err = memo.Scan(transaction.Memo); err != nil {
			err = errors.PropagateWithCode(err, EcodeStoreTransactionFailed, "Failed on scan memo value")
			return
		}
	}

	if transaction.JournalID == uuid.Nil {
		transaction.JournalID = uuid.New()
	}

	now := time.Now()

	if transaction.Date.IsZero() {
		transaction.Date = now
	}

	// startDate <= transDate <= endDate
	_, err = w.reader.GetFiscalYear(ctx, FiscalYearStatement{
		StartDateLTE: transaction.Date,
		EndDateGTE:   transaction.Date,
		ClosedNotEQ:  true,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.PropagateWithCode(err, EcodeStoreTransactionProhibited, "No active fiscal year for transaction date")
			return
		}

		err = errors.PropagateWithCode(err, EcodeGetFiscalYearFailed, "Failed on get fiscal year")
		return
	}

	for _, row := range transaction.Data {
		if row.Amount == 0 {
			continue
		}

		gls = append(gls, domain.GeneralLedger{
			ID:        uuid.New(),
			JournalID: transaction.JournalID,
			AccountID: row.AccountID,
			CreatedBy: userID,
			Amount:    row.Amount,
		})

		balanceAmount += row.Amount
		if row.Amount > 0 {
			journalAmount += row.Amount
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
		ID:        transaction.JournalID,
		Amount:    journalAmount,
		TransDate: transaction.Date,
		Memo:      memo,
		CreatedAt: now,
	}

	if err = w.StoreJournalTx(tx, ctx, journal); err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreTransactionAtJournalFailed, "Store transaction failed")
		return
	}

	if err = w.StoreGeneralLedgersTx(tx, ctx, gls); err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreTransactionAtGeneralLedgerFailed, "Store transaction failed")
		return
	}

	return
}

func (w *writer) StoreTransaction(ctx context.Context, userID uuid.UUID, transaction Transaction) (journal *domain.Journal, err error) {
	err = w.db.Transaction(ctx, nil, func(tx sql.Tx) error {
		journal, err = w.StoreTransactionTx(tx, ctx, userID, transaction)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		err = errors.PropagateWithCode(err, errors.GetCode(err), "Failed on store transaction")
		return
	}

	return
}

func (w *writer) StoreJournalTx(tx sql.Tx, ctx context.Context, journal *domain.Journal) (err error) {
	now := time.Now()

	if journal.TransDate.IsZero() {
		journal.TransDate = now
	}

	query := "INSERT INTO journals (id, amount, trans_date, memo) VALUES (?, ?, ?, ?) RETURNING id"
	err = tx.QueryRowContext(
		ctx,
		w.db.Rebind(query),
		journal.ID, journal.Amount, journal.TransDate, journal.Memo,
	).Scan(&journal.ID)

	if err != nil {
		err = errors.PropagateWithCode(err, EcodeStoreJournalFailed, "Store journal failed")
		return
	}

	return
}

func (w *writer) StoreGeneralLedgersTx(tx sql.Tx, ctx context.Context, gls []domain.GeneralLedger) (err error) {
	var params []interface{}

	query := `INSERT INTO general_ledgers (id, journal_id, account_id, amount, created_by) VALUES `

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
	if accountGroup.ParentID.Valid && accountGroup.ParentID.Int64 != 0 {
		if accountGroup.ParentID.Int64 == id {
			err = errors.PropagateWithCode(goErr.New("invalid parent id"), EcodeParentIDNotValid, "cannot set parent with the same account group")
			return
		}

		var parentAccountGroup domain.AccountGroup
		parentAccountGroup, err = w.reader.GetAccountGroupByID(ctx, accountGroup.ParentID.Int64)
		if err != nil {
			err = errors.PropagateWithCode(err, errors.GetCode(err), "Failed on get parent account group")
			return
		}

		accountGroup.ClassID = parentAccountGroup.ClassID
	}

	whereClause, whereClauseArgs, err := qb.NewWhereClause(AccountClassStatement{ID: id})
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

	_, err = w.db.ExecContext(ctx, w.db.Rebind(updateQuery), args...)
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
