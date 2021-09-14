package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

func GetValue(db *redis.Client, key string) ([]string, error) {

	if db == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	var values []string

	typ, err := db.Type(key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get type of key: %v, err: %v", key, err)
	}

	switch typ {
	case "string":
		data, err := db.Get(key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get value of key: %v, err: %v", key, err)
		}
		values = append(values, data)
	case "list":
		llen, err := db.LLen(key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get len of key: %v, err: %v", key, err)
		}

		var i int64
		for i = 0; i < llen; i++ {
			data, err := db.LIndex(key, i).Result()
			if err != nil {
				return nil, fmt.Errorf("failed to get list data of key: %v, err: %v", key, err)
			}
			values = append(values, data)
		}

	case "set":
		members, err := db.SMembers(key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get members of key: %v, err: %v", key, err)
		}
		values = append(values, members...)

	default:

	}

	return values, nil

}

func GetValues(db *redis.Client, keys ...string) ([][]string, error) {
	if db == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	var values [][]string

	for _, key := range keys {

		value, err := GetValue(db, key)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil

}
