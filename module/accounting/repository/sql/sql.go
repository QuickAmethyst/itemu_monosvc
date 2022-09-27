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
	reader := NewReader(opt)
	writer := NewWriter(opt, reader)

	return &struct {
		Reader
		Writer
	} {
		Reader: reader,
		Writer: writer,
	}
}
