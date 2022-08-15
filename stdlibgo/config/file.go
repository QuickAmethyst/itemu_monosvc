package config

import (
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	"github.com/QuickAmethyst/monosvc/stdlibgo/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
	"syscall"
)

var once = sync.Once{}

type File interface {
	ReadAndWatch(dest interface{})
}

type file struct {
	log logger.Logger
	v   *viper.Viper
	opt Option
}

func (f *file) ReadAndWatch(dest interface{}) {
	if err := f.v.ReadInConfig(); err != nil {
		err := errors.PropagateWithCode(err, errors.ErrReadConfig, "Failed to read config")
		f.log.Fatal(err.Error())
	}

	if err := f.v.Unmarshal(dest, decoderConfig); err != nil {
		err := errors.PropagateWithCode(err, errors.ErrUnmarshal, "Failed to unmarshal config")
		f.log.Fatal(err.Error())
	}

	f.v.WatchConfig()
	f.v.OnConfigChange(func(e fsnotify.Event) {
		f.notifyStaticConfigChange(f.opt.RestartOnChange, e)
	})
}

// notifyStaticConfigChange send sighup signal if any configurations change
func (f *file) notifyStaticConfigChange(restartOnChange bool, e fsnotify.Event) {
	f.log.Info(
		"File configuration modified",
		zap.String("name", e.Name),
		zap.String("path", f.v.ConfigFileUsed()),
	)

	if restartOnChange {
		once.Do(func() {
			f.log.Info("[Restarting App]")
			_ = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		})
	}
}

func NewFile(options Option) (File, error) {
	if err := options.validate(); err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigType(options.Type)
	v.SetConfigFile(options.Path)

	return &file{log: options.Logger, v: v, opt: options}, nil
}
