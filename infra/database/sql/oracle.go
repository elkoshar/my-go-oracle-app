package sql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
)

// DB defines common interface for both master & slave DB
type DB interface {
	Rebind(query string) string
	Ping() error
	Close() error
}

// MasterDB defines interface for master database
// which will be used for write operation
type MasterDB interface {
	DB
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Beginx() (*sqlx.Tx, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PreparexContext(ctx context.Context, query string) (MasterStatement, error)
	// QueryRowxContext placed in MasterDB in case you need to return something when perform INSERT or UPDATE query
	// never use this on SELECT query except on retry condition
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	// QueryContext placed in MasterDB in case you need to return rows when perform INSERT or UPDATE query
	// never use this on SELECT query except on retry condition
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// SlaveDB defines interface for slave database
// which will be used for read operation
type SlaveDB interface {
	DB
	PreparexContext(ctx context.Context, query string) (SlaveStatement, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

// MasterStatement is statement interface supported by master database
// it only contains write operation
type MasterStatement interface {
	Statement
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, args ...interface{}) error
	QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
}

// SlaveStatement is statement interface supported by slave database
// it only contains read operation
type SlaveStatement interface {
	Statement
	QueryRowxContext(ctx context.Context, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, args ...interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error
}

type TransactionDB interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	MustExec(query string, args ...interface{}) sql.Result
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Statement is common interface for both MasterStatement and SlaveStatement
type Statement interface {
	Close() error
}
type slaveDBImpl struct {
	*sqlx.DB
}
type masterDBImpl struct {
	*sqlx.DB
}

// OpenMasterDB open a master database from the given driver name and
// data source nama (DSN)
func OpenMasterDB(driver, dsn string, maxOpenConn, maxIdleConn int, connMaxIdleTime, connMaxLifeTime time.Duration) (MasterDB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(maxIdleConn)
	db.SetMaxOpenConns(maxOpenConn)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetConnMaxLifetime(connMaxLifeTime)
	return &masterDBImpl{
		DB: db,
	}, db.Ping()
}

// NewMasterDB creates new MasterDB object from existing sql.DB object
func NewMasterDB(db *sql.DB, driverName string) MasterDB {
	return &masterDBImpl{
		DB: sqlx.NewDb(db, driverName),
	}
}

func (mdb *masterDBImpl) PreparexContext(ctx context.Context, query string) (MasterStatement, error) {
	return mdb.DB.PreparexContext(ctx, query)
}

// NewSlaveDB creates new SlaveDB object from existing sql.DB object
func NewSlaveDB(db *sql.DB, driverName string) SlaveDB {
	return &slaveDBImpl{
		DB: sqlx.NewDb(db, driverName),
	}
}

// OpenSlaveDB open a slave database from the given driver name and
// data source name
func OpenSlaveDB(driver, dsn string, maxOpenConn, maxIdleConn int, connMaxIdleTime, connMaxLifeTime time.Duration) (SlaveDB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(maxIdleConn)
	db.SetMaxOpenConns(maxOpenConn)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetConnMaxLifetime(connMaxLifeTime)
	return &slaveDBImpl{
		DB: db,
	}, db.Ping()
}

func (sd *slaveDBImpl) PreparexContext(ctx context.Context, query string) (SlaveStatement, error) {
	return sd.DB.PreparexContext(ctx, query)
}
