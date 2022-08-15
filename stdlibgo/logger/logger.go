package logger

import (
	"go.uber.org/zap"
)

type Logger = *zap.Logger

type Option struct {
	Development bool
}

func New(option Option) (Logger, error) {
	var (
		log *zap.Logger
		err error
	)

	if option.Development {
		if log, err = zap.NewDevelopment(); err != nil {
			return nil, err
		}
	} else {
		if log, err = zap.NewProduction(); err != nil {
			return nil, err
		}
	}

	defer log.Sync()

	return log, nil
}
