package redis_test

import (
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
	redis_ "github.com/kaydxh/golang/pkg/database/redis"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetDataBase(t *testing.T) {
	testCases := []struct {
		TestName  string
		Addresses []string
		DB        int
		UserName  string
		Password  string
	}{
		{
			TestName:  "test1",
			Addresses: []string{"9.135.232.102:6380"},
			DB:        0,
			UserName:  "root",
			Password:  "HXufW*3569FShs",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			db := redis_.NewRedisClient(redis_.DBConfig{
				Addresses: testCase.Addresses,
				UserName:  testCase.UserName,
				Password:  testCase.Password,
				DB:        testCase.DB,
			})
			redisDB, err := db.GetRedis()
			if err != nil {
				t.Fatalf("failed to get redis database: %v, got : %s", testCase.Addresses, err)
			}
			assert.NotNil(t, redisDB)
		})
	}

}

func GetDBOrDie() *redis.Client {

	var (
		once sync.Once
		db   *redis.Client
		err  error
	)

	once.Do(func() {
		cfgFile := "./redis.yaml"
		config := redis_.NewConfig(redis_.WithViper(viper_.GetViper(cfgFile, "database.redis")))

		db, err = config.Complete().New()
		if err != nil {
			panic(err)
		}
		if db == nil {
			panic("db is not enable")
		}
	})

	return db

}

func TestNew(t *testing.T) {
	db := GetDBOrDie()
	defer db.Close()

	t.Logf("db: %#v", db)
}

// set string
// Redis `SET key value [expiration]` command.
//
// Use expiration for `SETEX`-like behavior.
// Zero expiration means the key has no expiration time.
func TestSet(t *testing.T) {

	db := GetDBOrDie()
	defer db.Close()

	testCases := []struct {
		key      string
		value    string
		expire   time.Duration
		expected string
	}{
		{
			key:      "test1",
			value:    "test1-1, test1-2",
			expected: "test1",
		},
		{
			key:      "test2",
			value:    "test2-1, test2-2",
			expected: "test2",
		},

		{
			key:      "test3-tmp",
			value:    "test3-1, test3-2",
			expire:   time.Minute,
			expected: "test3",
		},
	}

	for _, testCase := range testCases {

		result, err := db.Set(testCase.key, testCase.value, testCase.expire).Result()
		if err != nil {
			t.Fatalf("failed to set string, err: %v", err)
		}

		t.Logf("result of %v: %v", testCase.key, result)
	}

}

func TestKeys(t *testing.T) {

	db := GetDBOrDie()
	defer db.Close()

	keys, err := db.Keys("*").Result()
	if err != nil {
		t.Fatalf("failed to get all keys , err: %v", err)
	}

	for _, key := range keys {
		typ, err := db.Type(key).Result()
		if err != nil {
			t.Fatalf("failed to get type of key: %v, err: %v", key, err)
		}

		if typ == "string" {
			data, err := db.Get(key).Result()
			if err != nil {
				t.Fatalf("failed to get value of key: %v, err: %v", key, err)
			}

			t.Logf(" key %v, value: %v ", key, data)
		}

	}
}

/*
func TestGetDatabaseUntil(t *testing.T) {
	testCases := []struct {
		Address  string
		DataName string
		UserName string
		Password string
	}{
		{
			Address:  "127.0.0.1:3306",
			DataName: "test",
			UserName: "root",
			Password: "123456",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Address+testCase.DataName, func(t *testing.T) {
			db := mysql_.NewDB(mysql_.DBConfig{
				Address:  testCase.Address,
				DataName: testCase.DataName,
				UserName: testCase.UserName,
				Password: testCase.Password,
			})
			sqlDB, err := db.GetDatabaseUntil(5*time.Second, 20*time.Second)
			if err != nil {
				t.Fatalf("failed to get database: %v, got : %s", testCase.DataName, err)
			}
			assert.NotNil(t, sqlDB)
		})
	}

}

func TestGetTheDBAndClose(t *testing.T) {
	testCases := []struct {
		Address  string
		DataName string
		UserName string
		Password string
	}{
		{
			Address:  "127.0.0.1:3306",
			DataName: "test",
			UserName: "root",
			Password: "123456",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Address+testCase.DataName, func(t *testing.T) {
			conf := mysql_.DBConfig{
				Address:  testCase.Address,
				DataName: testCase.DataName,
				UserName: testCase.UserName,
				Password: testCase.Password,
			}
			db := mysql_.NewDB(conf)
			sqlDB, err := db.GetDatabaseUntil(5*time.Second, 20*time.Second)
			if err != nil {
				t.Fatalf("failed to get database: %v, got : %s", testCase.DataName, err)
			}
			assert.NotNil(t, sqlDB)

			theDB, err := mysql_.GetTheDB(conf)
			assert.Nil(t, err)

			assert.Equal(t, sqlDB, mysql_.GetDB())
			t.Logf("GetDB got sqlDB: %v, expect %v", sqlDB, mysql_.GetDB())
			t.Logf("GetTheDB got sqlDB: %v, expect %v", sqlDB, theDB)
			err = mysql_.CloseTheDB(conf)
			assert.Nil(t, err)
			err = mysql_.CloseDB()
			assert.Nil(t, err)
		})
	}

}
*/
