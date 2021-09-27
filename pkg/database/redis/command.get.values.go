package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func GetValue(ctx context.Context, db *redis.Client, key string) ([]string, error) {

	if db == nil {
		return nil, fmt.Errorf("found unexpected nil redis client")
	}

	var values []string

	typ, err := db.Type(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get type of key: %v, err: %v", key, err)
	}

	switch typ {
	case "string":
		data, err := db.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get value of key: %v, err: %v", key, err)
		}
		values = append(values, data)
	case "list":
		llen, err := db.LLen(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get len of key: %v, err: %v", key, err)
		}

		var i int64
		for i = 0; i < llen; i++ {
			data, err := db.LIndex(ctx, key, i).Result()
			if err != nil {
				return nil, fmt.Errorf("failed to get list data of key: %v, err: %v", key, err)
			}
			values = append(values, data)
		}

	case "set":
		members, err := db.SMembers(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get members of key: %v, err: %v", key, err)
		}
		values = append(values, members...)

	default:

	}

	return values, nil

}

func GetValues(ctx context.Context, db *redis.Client, keys ...string) ([][]string, error) {
	if db == nil {
		return nil, fmt.Errorf("found unexpected nil redis client")
	}

	var values [][]string

	for _, key := range keys {

		value, err := GetValue(ctx, db, key)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil

}
