package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	time_ "github.com/kaydxh/golang/go/time"
)

type DBConfig struct {
	Address    string
	DataName   string
	UserName   string
	Password   string
	DriverName string
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

func NewDB(conf DBConfig, opts ...DBOption) (*DB, error) {
	conn := &DB{
		Conf: conf,
	}

	conn.ApplyOptions(opts...)

	return conn, nil
}

func (d *DB) GetDatabase() (*sqlx.DB, error) {
	if d.db != nil {
		return d.db, nil
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",
		d.Conf.UserName,
		d.Conf.Password,
		d.Conf.Address,
		d.Conf.DataName,
	)

	db, err := sqlx.Open(d.Conf.DriverName, dsn)
	if err != nil {
		return nil, err
	}
	if err := d.db.Ping(); err != nil {
		return nil, err
	}

	d.db.SetMaxOpenConns(d.opts.MaxConns)
	d.db.SetMaxIdleConns(d.opts.MaxIdleConns)
	d.db.SetConnMaxLifetime(d.opts.ConnMaxLifetime)

	d.db = db
	return d.db, nil
}

func (d *DB) GetDatabaseUntil(maxWaitInterval time.Duration, failAfter time.Duration) (*sqlx.DB, error) {

	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionMaxInterval(maxWaitInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(failAfter),
	)
	for {
		db, err := d.GetDatabase()
		if err != nil {
			continue
		}

		actualInterval, over := exp.NextBackOff()
		if over {
			return nil, fmt.Errorf("get datqabase fail after: %v", failAfter)
		}

		time.Sleep(actualInterval)
		return db, nil
	}
}