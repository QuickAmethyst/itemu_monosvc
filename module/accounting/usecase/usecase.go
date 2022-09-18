package usecase

import "github.com/QuickAmethyst/monosvc/module/accounting/repository/sql"

type Options struct {
	AccountingSQL sql.SQL
}

type Usecase interface {
	Reader
	Writer
}

func New(opt *Options) Usecase {
	return &struct {
		Reader
		Writer
	} {
		Reader: NewReader(opt),
		Writer: NewWriter(opt),
	}
}