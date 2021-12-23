package mysql

import (
	"context"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

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
		maxConns     int
		maxIdleConns int
		dialTimeout  time.Duration
		readTimeout  time.Duration
		writeTimeout time.Duration
		// connection reused time, 0 means never expired
		connMaxLifetime   time.Duration
		interpolateParams bool
	}
}

func NewDB(conf DBConfig, opts ...DBOption) *DB {
	conn := &DB{
		Conf: conf,
	}
	conn.opts.maxConns = DefaultMaxConns
	conn.opts.maxIdleConns = DefaultMaxIdleConns

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

	dsnFull := fmt.Sprintf("%s%s", dsn, func() string {
		var params string
		if d.opts.dialTimeout > 0 {
			params += fmt.Sprintf("&timeout=%fs", d.opts.dialTimeout.Seconds())
		}

		if d.opts.readTimeout > 0 {
			params += fmt.Sprintf("&readTimeout=%fs", d.opts.readTimeout.Seconds())
		}

		if d.opts.writeTimeout > 0 {
			params += fmt.Sprintf("&writeTimeout=%fs", d.opts.writeTimeout.Seconds())
		}

		if d.opts.interpolateParams {
			params += "&interpolateParams=true"
		} else {
			params += "&interpolateParams=false"
		}

		return params

	}())

	db, err := sqlx.Open("mysql", dsnFull)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(d.opts.maxConns)
	db.SetMaxIdleConns(d.opts.maxIdleConns)
	db.SetConnMaxLifetime(d.opts.connMaxLifetime)

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

func (d *DB) GetDatabaseUntil(
	ctx context.Context,
	maxWaitInterval time.Duration,
	failAfter time.Duration,
) (*sqlx.DB, error) {

	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionMaxInterval(maxWaitInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(failAfter),
	)
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("cancel get database: %v", ctx.Err())

		default:
			db, err := d.GetDatabase()
			if err == nil {
				return db, nil

			} else {
				actualInterval, ok := exp.NextBackOff()
				if !ok {
					return nil, fmt.Errorf("get database fail after: %v", failAfter)
				}

				logrus.Infof("delay %v, try again", actualInterval)
				time.Sleep(actualInterval)
			}
		}
	}
}

func (d *DB) Close() error {
	if d.db == nil {
		return fmt.Errorf("no database pool")
	}
	return d.db.Close()
}
