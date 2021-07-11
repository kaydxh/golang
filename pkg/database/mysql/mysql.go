package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	time_ "github.com/kaydxh/golang/go/time"
)

var (
	sqlDB  SQLDB
	sqlDBs map[DBConfig]SQLDB
)

type DBConfig struct {
	Address  string
	DataName string
	UserName string
	Password string
}

type DB struct {
	//DSN        string
	Conf DBConfig
	db   *sqlx.DB

	opts struct {
		MaxConns        int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration
	}
}

func NewDB(conf DBConfig, opts ...DBOption) *DB {
	conn := &DB{
		Conf: conf,
	}

	conn.ApplyOptions(opts...)

	return conn
}

func Get() *sqlx.DB {
	return sqlDB.Load()
}

func (d *DB) GetDatabase() (*sqlx.DB, error) {
	if d.db != nil {
		return d.db, nil
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&loc=Local&parseTime=true",
		d.Conf.UserName,
		d.Conf.Password,
		d.Conf.Address,
		d.Conf.DataName,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(d.opts.MaxConns)
	db.SetMaxIdleConns(d.opts.MaxIdleConns)
	db.SetConnMaxLifetime(d.opts.ConnMaxLifetime)

	d.db = db
	sqlDB.Store(db)
	return d.db, nil
}

func (d *DB) GetDatabaseUntil(maxWaitInterval time.Duration, failAfter time.Duration) (*sqlx.DB, error) {

	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionMaxInterval(maxWaitInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(failAfter),
	)
	for {
		db, err := d.GetDatabase()
		if err == nil {
			return db, nil

		} else {
			actualInterval, ok := exp.NextBackOff()
			if !ok {
				return nil, fmt.Errorf("get database fail after: %v", failAfter)
			}

			time.Sleep(actualInterval)
		}
	}
}

func (d *DB) Close() error {
	if d.db == nil {
		return fmt.Errorf("no database pool")
	}
	return d.db.Close()
}
