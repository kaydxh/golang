package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	reflect_ "github.com/kaydxh/golang/go/reflect"
)

// only get export Fields from arg
func HSetStruct(ctx context.Context, db *redis.Client, key string, arg interface{}) error {

	if db == nil {
		return fmt.Errorf("redis client is nil")
	}

	tagsValues := reflect_.AllTagsValues(arg, "db")
	_, err := db.HMSet(ctx, key, tagsValues).Result()

	return err
}
