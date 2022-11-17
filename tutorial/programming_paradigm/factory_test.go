package tutorial

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type FactoryConfigFunc func(c *FactoryConfig) error

type FactoryConfig struct {
	Addr    string
	Timeout time.Duration //接口处理超时时间
}

func (fc *FactoryConfig) ApplyOptions(configFuncs ...FactoryConfigFunc) error {

	for _, f := range configFuncs {
		err := f(fc)
		if err != nil {
			return fmt.Errorf("failed to apply factory config, err: %v", err)
		}
	}

	return nil
}

func (fc FactoryConfig) Validate() error {
	return nil
}

func (fc FactoryConfig) Do(ctx context.Context) error {
	fmt.Println("done")
	return nil
}

type Factory struct {
	fc FactoryConfig
}

func NewFactory(fc FactoryConfig, configFuncs ...FactoryConfigFunc) (Factory, error) {
	err := fc.ApplyOptions(configFuncs...)
	if err != nil {
		return Factory{}, err
	}

	err = fc.Validate()
	if err != nil {
		return Factory{}, err
	}

	return Factory{fc: fc}, nil
}

func (f Factory) NewRepository(ctx context.Context) (Repository, error) {
	//the Repository can also not FactoryConfig instance
	// the example use FactoryConfig
	return f.fc, nil
}

type Repository interface {
	Do(ctx context.Context) error
}

func (f Factory) Do(ctx context.Context) error {
	return f.fc.Do(ctx)
}

func TestFactroy(t *testing.T) {
	factory, err := NewFactory(FactoryConfig{
		Addr:    "localhost:10001",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Errorf("failed to new factory, err: %v", err)
	}

	ctx := context.Background()
	repository, err := factory.NewRepository(ctx)
	if err != nil {
		t.Errorf("failed to new repository, err: %v", err)
	}
	err = repository.Do(ctx)
	if err != nil {
		t.Errorf("failed to call Do method, err: %v", err)
	}
}
