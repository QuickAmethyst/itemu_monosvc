package sql

import (
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	"github.com/QuickAmethyst/monosvc/stdlibgo/sql"
)

type SQL interface {
	Reader
	Writer
}

type Options struct {
	MasterDB sql.DB
	SlaveDB  sql.DB
	Logger   logger.Logger
}

func New(opt *Options) SQL {
	return struct {
		Reader
		Writer
	}{
		Reader: NewReader(opt),
		Writer: NewWriter(opt),
	}
}
