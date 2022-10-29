package sql

type GeneralLedgerPreferenceID int64

const (
	RetainedEarnings GeneralLedgerPreferenceID = iota + 1
	ProfitLossYear   GeneralLedgerPreferenceID = iota
)
