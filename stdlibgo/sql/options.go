package sql

import (
	"errors"
	"time"
)

type ConnectionOptions struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	SSL          bool
	MaxLifeTime  time.Duration
	MaxIdle      int
	MaxOpen      int
}

type PostgresSQLOptions struct {
	Master ConnectionOptions
	Slave  ConnectionOptions
}

func validateConf(conf *ConnectionOptions) error {
	if conf.Host == "" {
		return errors.New("db: host is required")
	}

	if conf.DatabaseName == "" {
		return errors.New("db: database name is required")
	}

	if conf.Password == "" {
		return errors.New("db: password is required")
	}

	if conf.Port == 0 {
		conf.Port = 5432
	}

	if conf.User == "" {
		return errors.New("db: user is required")
	}

	return nil
}
