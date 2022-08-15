package sql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresSQL interface {
	Master() DB
	Slave() DB
}

type postgres struct {
	master DB
	slave  DB
}

func (d *postgres) Master() DB {
	return d.master
}

func (d *postgres) Slave() DB {
	return d.slave
}

func NewPostgresSQL(opt PostgresSQLOptions) (PostgresSQL, error) {
	var (
		err      error
		dbMaster *sqlx.DB
		dbSlave  *sqlx.DB
	)

	if err := validateConf(&opt.Master); err != nil {
		return nil, err
	}

	if err := validateConf(&opt.Slave); err != nil {
		return nil, err
	}

	dbMaster, err = connect(opt.Master)
	if err != nil {
		return nil, err
	}

	dbSlave, err = connect(opt.Slave)
	if err != nil {
		return nil, err
	}

	return &postgres{
		master: NewDB(dbMaster),
		slave:  NewDB(dbSlave),
	}, nil
}

func connect(opt ConnectionOptions) (*sqlx.DB, error) {
	const driverName = "postgres"

	db, err := sql.Open(driverName, BuildPostgresUri(opt))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(opt.MaxOpen)
	db.SetMaxIdleConns(opt.MaxIdle)
	db.SetConnMaxLifetime(opt.MaxLifeTime)

	return sqlx.NewDb(db, driverName), nil
}