package usecase

import "github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"

type Options struct {
	AccountingSQL sql.SQL
}

type Usecase interface {
	Writer
}

func New(opt *Options) Usecase {
	return &struct {
		Writer
	} {
		Writer: NewWriter(opt),
	}
}