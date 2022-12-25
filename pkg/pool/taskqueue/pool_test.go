package taskqueue_test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/go-redis/redis/v8"
	redis_ "github.com/kaydxh/golang/pkg/database/redis"
	taskq_ "github.com/kaydxh/golang/pkg/pool/taskqueue"
	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"
	redisq_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue/redis"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"golang.org/x/net/context"
)

type TaskA struct {
}

type TaskAArgs struct {
	Param1 string
	Param2 int
}

func (t TaskA) Scheme() string {
	return "taskA"
}

func (t TaskA) TaskHandler(message *queue_.Message) error {
	var args TaskAArgs
	err := json.Unmarshal([]byte(message.Args), &args)
	if err != nil {
		return err
	}

	fmt.Printf("get args: %v\n", args)
	fmt.Printf("process task: %v\n", message.Id)

	return nil
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

func TestTaskQueue(t *testing.T) {

	taskq_.Register(TaskA{})

	db := GetDBOrDie()

	redisq := redisq_.NewQueue(db, queue_.QueueOptions{
		Name: "redis",
	})

	pool := taskq_.NewPool(redisq)
	ctx := context.Background()

	args := TaskAArgs{
		Param1: "param1",
		Param2: 10,
	}

	data, err := json.Marshal(args)
	if err != nil {
		t.Errorf("failed to marshal args, err: %v", err)
		return
	}
	msg := &queue_.Message{
		Name:   "taskA-1",
		Scheme: "taskA",
		Args:   string(data),
	}
	err = pool.Publish(ctx, msg)
	if err != nil {
		t.Errorf("failed to pulibsh task, err: %v", err)
		return
	}

	err = pool.Consume(ctx)
	if err != nil {
		t.Errorf("failed to consume task, err: %v", err)
		return
	}

	select {}

}
