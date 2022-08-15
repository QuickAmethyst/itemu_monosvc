package config

import (
	"errors"
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
)

type Option struct {
	Path            string
	Type            string
	RestartOnChange bool
	Logger          logger.Logger
}

func (o *Option) validate() error {
	if o.Path == "" {
		return errors.New("config: Options.Path is required")
	}

	if o.Type == "" {
		return errors.New("config: Options.Type is required")
	}

	if o.Logger == nil {
		return errors.New("config: Options.Logger is required")
	}

	return nil
}
