package taskqueue_test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

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

type TaskAResult struct {
	Result1 string
	Result2 int
}

func (t TaskA) Scheme() string {
	return "taskA"
}

func (t TaskA) TaskHandler(ctx context.Context, msg *queue_.Message) (*queue_.MessageResult, error) {

	time.Sleep(30 * time.Second)
	var args TaskAArgs
	err := json.Unmarshal([]byte(msg.Args), &args)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("get args: %v\n", args)
	fmt.Printf("process task id: %v, name: %v\n", msg.Id, msg.Name)

	result := &queue_.MessageResult{
		Id:      msg.Id,
		InnerId: msg.InnerId,
		Name:    msg.Name,
		Scheme:  msg.Scheme,
	}

	taskResult := &TaskAResult{
		Result1: "Result1",
		Result2: 2,
	}

	data, err := json.Marshal(taskResult)
	if err != nil {
		return nil, err
	}
	result.Result = string(data)
	return result, nil
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

	pool := taskq_.NewPool(redisq, taskq_.WithFetcherBurst(1), taskq_.WithWorkTimeout(-1), taskq_.WithResultCallbackFunc(func(ctx context.Context, result *queue_.MessageResult) {
		t.Logf("--> callback fetch result %v of msg", result)
	}))
	ctx := context.Background()
	err = pool.Consume(ctx)
	if err != nil {
		t.Errorf("failed to consume task, err: %v", err)
		return
	}

	args := TaskAArgs{
		Param1: "param1",
		Param2: 10,
	}

	data, err := json.Marshal(args)
	if err != nil {
		t.Errorf("failed to marshal args, err: %v", err)
		return
	}

	var (
		ids []string
		wg  sync.WaitGroup
	)
	for i := 0; i < 10; i++ {
		msg := &queue_.Message{
			Scheme: "taskA",
			Args:   string(data),
		}
		wg.Add(1)
		go func(i int, m *queue_.Message) {
			defer wg.Done()
			m.Name = fmt.Sprintf("taskA-%v", i)
			id, err := pool.Publish(ctx, m)
			if err != nil {
				t.Errorf("failed to pulibsh task, err: %v", err)
				return
			}
			t.Logf("pulibsh task, innerId: %v, msg: %v", id, m)
			ids = append(ids, id)
		}(i, msg)
	}

	wg.Wait()

	//	time.Sleep(5 * time.Second)

	/*
		for _, id := range ids {
			result, err := pool.FetchResult(ctx, id)
			if err != nil {
				t.Errorf("%v", err)
				continue
			}
			t.Logf("fetch result %v of msg %v", result, id)
		}
	*/
	select {}

}
