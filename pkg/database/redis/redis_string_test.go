package redis_test

import (
	"sync"
	"testing"
	"time"

	"context"

	"github.com/go-redis/redis/v8"
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
			redisDB, err := db.GetRedis(context.Background())
			if err != nil {
				t.Fatalf("failed to get redis database: %v, got : %s", testCase.Addresses, err)
			}
			assert.NotNil(t, redisDB)
		})
	}

}

var (
	once sync.Once
	db   *redis.Client
	err  error
)

func GetDBOrDie() *redis.Client {
	once.Do(func() {
		cfgFile := "./redis.yaml"
		config := redis_.NewConfig(redis_.WithViper(viper_.GetViper(cfgFile, "database.redis")))

		db, err = config.Complete().New(context.Background())
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
	//defer db.Close()

	t.Logf("db: %#v", db)
}

// set string
// Redis `SET key value [expiration]` command.
//
// Use expiration for `SETEX`-like behavior.
// Zero expiration means the key has no expiration time.
func TestSet(t *testing.T) {

	db := GetDBOrDie()
	//defer db.Close()

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

	ctx := context.Background()
	for _, testCase := range testCases {

		result, err := db.Set(ctx, testCase.key, testCase.value, testCase.expire).Result()
		if err != nil {
			t.Fatalf("failed to set string, err: %v", err)
		}

		t.Logf("result of %v: %v", testCase.key, result)
	}

}

//get values with keys
func TestGetValues(t *testing.T) {

	db := GetDBOrDie()
	//	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	keys, err := db.Keys(ctx, "*").Result()
	if err != nil {
		t.Fatalf("failed to get all keys , err: %v", err)
	}

	values, err := redis_.GetValues(ctx, db, keys...)
	if err != nil {
		t.Fatalf("failed to get values, err: %v", err)
	}
	t.Logf("keys: %v, values: %v", keys, values)
}

//get range of value with key
func TestGetRange(t *testing.T) {

	db := GetDBOrDie()
	//	defer db.Close()

	testCases := []struct {
		key      string
		start    int64
		end      int64
		expected string
	}{
		{
			key:      "test1",
			start:    0,
			end:      -1, // get all range
			expected: "test1-1",
		},
		{
			key:      "test2",
			start:    0,
			end:      int64(len("test1-1")) - 1, //include end
			expected: "test2-1",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		value, err := db.GetRange(ctx, testCase.key, testCase.start, testCase.end).Result()
		if err != nil {
			t.Fatalf("failed to get range, err: %v", err)
		}
		t.Logf("key: %v, range [%d:%d] value: %v", testCase.key, testCase.start, testCase.end, value)
	}
}
