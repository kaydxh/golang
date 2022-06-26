package instance_test

import (
	"context"
	"sync"
	"testing"
	"time"

	instance_ "github.com/kaydxh/golang/pkg/pool/instance"
	"github.com/sirupsen/logrus"
)

type InstancePool struct {
	pool *instance_.Pool
}

type Sdk struct {
}

func (s *Sdk) DoSdk() error {
	logrus.Infof("do DoSdk")
	time.Sleep(2 * time.Second)
	return nil
}

func (s *InstancePool) InvokeProcess(ctx context.Context) error {
	//logrus.Infof("do InvokeProcess")
	resp, err := s.pool.Invoke(ctx, func(ctx context.Context, instance interface{}) (interface{}, error) {
		return nil, instance.(*Sdk).DoSdk()
	})
	if err != nil {
		logrus.Errorf("do DoSdk err: %v", err)
	}
	_ = resp

	return err
}

func GlobalInit() error {
	logrus.Infof("do GlobaInit")
	return nil
}

func GlobalRelease() error {
	logrus.Infof("do GlobalRelease")
	return nil
}

func LocalInit() error {
	logrus.Infof("do LocalInit")
	return nil
}

func LocalRelease() error {
	logrus.Infof("do LocalRelease")
	return nil
}

func New() error {
	logrus.Infof("do New")
	return nil
}

func Delete() {
	logrus.Infof("do Delete")
}

func TestNewPool(t *testing.T) {
	ctx := context.Background()
	pool, err := instance_.NewPool(func() interface{} {
		New()
		return &Sdk{}
	},
		instance_.WithGpuIDs([]int64{0, 1, 2}),
		instance_.WithBatchSize(8),
		instance_.WithWaitTimeoutOnce(50*time.Millisecond),
		instance_.WithWaitTimeoutTotal(time.Second),
		instance_.WithName("test-instance"),
		instance_.WithEnabledPrintCostTime(true),
		//instance_.WithWaitTimeout(time.Millisecond),
		instance_.WithResevePoolSizePerGpu(1),
		instance_.WithCapacityPoolSizePerGpu(1),
		instance_.WithGlobalInitFunc(func() error {
			return GlobalInit()
		}),
		instance_.WithGlobalReleaseFunc(func() error {
			return GlobalRelease()
		}),
		instance_.WithLocalInitFunc(func(instace interface{}) error {
			return LocalInit()
		}),
		instance_.WithLocalReleaseFunc(func(instace interface{}) error {
			return LocalRelease()
		}),
		instance_.WithDeleteFunc(func(instace interface{}) {
			Delete()
		}),
	)

	err = pool.GlobalInit(ctx)
	if err != nil {
		t.Errorf("global init err: %v", err)
		return
	}

	defer pool.GlobalRelease(ctx)

	//do somthing
	instancePool := &InstancePool{
		pool: pool,
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = instancePool.InvokeProcess(ctx)
		}()
	}

	wg.Wait()

	return

}
