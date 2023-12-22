/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package etcd

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	time_ "github.com/kaydxh/golang/go/time"
	"github.com/sirupsen/logrus"
)

type EventCallbackFunc func(ctx context.Context, key, value string)

var (
	etcdKV ETCDKV
	mu     sync.Mutex
)

// Default values for Etcd.
const (
	DefaultDialTimeout = 5 * time.Second
)

type EtcdConfig struct {
	Addresses []string
	UserName  string
	Password  string
}

type EtcdKVOptions struct {
	DialTimeout        time.Duration
	MaxCallSendMsgSize int
	MaxCallRecvMsgSize int
	AutoSyncInterval   time.Duration

	WatchPaths         []string
	CreateCallbackFunc EventCallbackFunc
	DeleteCallbackFunc EventCallbackFunc
}

type EtcdKV struct {
	Conf EtcdConfig
	kv   *clientv3.Client

	opts EtcdKVOptions
}

func NewEtcdKV(conf EtcdConfig, opts ...EtcdKVOption) *EtcdKV {
	kv := &EtcdKV{
		Conf: conf,
	}
	kv.opts.DialTimeout = DefaultDialTimeout
	kv.ApplyOptions(opts...)

	return kv
}

func GetKV() *clientv3.Client {
	return etcdKV.Load()
}

func CloseKV() error {
	if etcdKV.Load() == nil {
		return nil
	}

	return etcdKV.Load().Close()
}

func (d *EtcdKV) GetKV(ctx context.Context) (*clientv3.Client, error) {
	if d.kv != nil {
		return d.kv, nil
	}

	if len(d.Conf.Addresses) == 0 {
		return nil, fmt.Errorf("invalid etcd address")
	}
	logrus.Infof("dialTimeout: %v", d.opts.DialTimeout)

	kv, err := clientv3.New(clientv3.Config{
		Endpoints:          d.Conf.Addresses,
		Context:            ctx,
		Username:           d.Conf.UserName,
		Password:           d.Conf.Password,
		MaxCallRecvMsgSize: d.opts.MaxCallRecvMsgSize,
		MaxCallSendMsgSize: d.opts.MaxCallSendMsgSize,
		AutoSyncInterval:   d.opts.AutoSyncInterval,
		DialTimeout:        d.opts.DialTimeout,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, d.opts.DialTimeout)
	defer cancel()
	status, err := kv.Status(ctx, d.Conf.Addresses[0])
	if err != nil {
		return nil, err
	}
	logrus.Infof("kv status: %v", status)

	d.kv = kv
	etcdKV.Store(kv)

	return kv, nil
}

func (d *EtcdKV) GetKVUntil(
	ctx context.Context,
	maxWaitInterval time.Duration,
	failAfter time.Duration,
) (*clientv3.Client, error) {
	var kv *clientv3.Client
	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionMaxInterval(maxWaitInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(failAfter),
	)
	err := time_.BackOffUntilWithContext(ctx, func(ctx context.Context) (err_ error) {
		kv, err_ = d.GetKV(ctx)
		if err_ != nil {
			return err_
		}
		return nil
	}, exp, true, false)
	if err != nil {
		return nil, fmt.Errorf("get etcd fail after: %v", failAfter)
	}

	return kv, nil

}

func (d *EtcdKV) Watch(ctx context.Context) {
	for _, path := range d.opts.WatchPaths {
		fmt.Printf("watch path: %v\n", path)
		Watch(ctx, d.kv, path, d.opts.CreateCallbackFunc, d.opts.DeleteCallbackFunc)
	}
}

func (d *EtcdKV) Close() error {
	if d.kv == nil {
		return fmt.Errorf("no etcd pool")
	}
	return d.kv.Close()
}
