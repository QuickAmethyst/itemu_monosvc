package usecase

import (
	"github.com/QuickAmethyst/monosvc/module/account/repository/sql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/auth"
)

type Usecase interface {
	Reader
}

type Options struct {
	AccountSQL sql.SQL
	Auth       auth.Auth
}

func New(opt *Options) Usecase {
	return &struct {
		Reader
	}{
		Reader: NewReader(opt),
	}
}
