package usecase

import "github.com/QuickAmethyst/monosvc/module/inventory/repository/sql"

type Usecase interface {
	Reader
	Writer
}

type Options struct {
	InventorySQL sql.SQL
}

func New(opt *Options) Usecase {
	return &struct {
		Reader
		Writer
	}{
		Reader: NewReader(opt),
		Writer: NewWriter(opt),
	}
}
