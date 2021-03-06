package mysql

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	time_ "github.com/kaydxh/golang/go/time"
)

var (
	sqlDB  SQLDB
	sqlDBs map[DBConfig]SQLDB
	mu     sync.Mutex
)

// Default values for Mysql.
const (
	DefaultMaxConns     = 100
	DefaultMaxIdleConns = 10
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
	conn.opts.MaxConns = DefaultMaxConns
	conn.opts.MaxIdleConns = DefaultMaxIdleConns

	conn.ApplyOptions(opts...)

	return conn
}

func GetDB() *sqlx.DB {
	return sqlDB.Load()
}

func GetTheDB(conf DBConfig) (*sqlx.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	sqlDB, ok := sqlDBs[conf]
	if !ok {
		return nil, fmt.Errorf("not found the db in cache")
	}
	return sqlDB.Load(), nil
}

func CloseDB() error {
	if sqlDB.Load() == nil {
		return nil
	}

	return sqlDB.Load().Close()
}

func CloseTheDB(conf DBConfig) error {
	mu.Lock()
	defer mu.Unlock()

	sqlDB, ok := sqlDBs[conf]
	if !ok {
		return fmt.Errorf("not found the db in cache")
	}
	if sqlDB.Load() == nil {
		return nil
	}

	err := sqlDB.Load().Close()
	if err != nil {
		return err
	}

	delete(sqlDBs, conf)
	return nil
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

	mu.Lock()
	defer mu.Unlock()
	if sqlDBs == nil {
		sqlDBs = make(map[DBConfig]SQLDB)
	}
	sqlDBs[d.Conf] = sqlDB
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
