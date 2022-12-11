package sql

const (
	EcodeStoreAccountClassFailed = iota + 1
	EcodeUpdateAccountClassFailed
	EcodeDeleteAccountClassFailed
	EcodeGetAccountClassListFailed
	EcodeGetAccountClassFailed
	EcodeStoreAccountGroupFailed
	EcodeUpdateAccountGroupFailed
	EcodeDeleteAccountGroupFailed
	EcodeGetAccountGroupFailed
	EcodeGetAllTopLevelAccountGroupFailed
	EcodeGetAllAccountGroupsFailed
	EcodeParentIDNotValid
	EcodeStoreAccountFailed
	EcodeUpdateAccountFailed
	EcodeDeleteAccountFailed
	EcodeGetAllAccountsFailed
	EcodeGetAccountFailed
	EcodeStoreJournalFailed
	EcodeStoreTransactionFailed
	EcodeStoreTransactionAtJournalFailed
	EcodeStoreTransactionAtGeneralLedgerFailed
	EcodeStoreTransactionCreatedByRequired
	EcodeTransactionNotBalance
	EcodeUpdateGeneralLedgerPreferenceFailed
	EcodeGetAllGeneralLedgerPreferencesFailed
	EcodeValidatePreferencesFailed
	EcodeStoreFiscalYearFailed
	EcodeValidateFiscalYearFailed
	EcodeGetFiscalYearListFailed
	EcodeGetFiscalYearFailed
	EcodeGetActiveFiscalYearFailed
	EcodeUpdateFiscalYearFailed
	EcodeCloseFiscalYearFailed
	EcodeGetBalanceSheetAmountFailed
	EcodeStoreTransactionProhibited
	EcodeGetGeneralLedgerPreferenceFailed
	EcodeStoreBankAccountFailed
	EcodeUpdateBankAccountFailed
	EcodeGetBankAccountListFailed
	EcodeGetBankAccountFailed
	EcodeGetGeneralLedgerFailed
	EcodeAccountHasTransaction
)
