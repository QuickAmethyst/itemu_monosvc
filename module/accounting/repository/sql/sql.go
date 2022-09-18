package sql

import (
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
)

type Options struct {
	MasterDB sql.DB
	SlaveDB  sql.DB
	Logger   logger.Logger
}

type SQL interface {
	Reader
	Writer
}

func New(opt *Options) SQL {
	return &struct {
		Reader
		Writer
	} {
		Reader: NewReader(opt),
		Writer: NewWriter(opt),
	}
}
