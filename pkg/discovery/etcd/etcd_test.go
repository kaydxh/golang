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
package etcd_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	etcd_ "github.com/kaydxh/golang/pkg/discovery/etcd"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetKV(t *testing.T) {
	testCases := []struct {
		Addresses []string
	}{
		{
			Addresses: []string{"127.0.0.1:3306"},
		},
	}

	ctx := context.Background()
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("name-%d", i), func(t *testing.T) {
			kv := etcd_.NewEtcdKV(etcd_.EtcdConfig{
				Addresses: testCase.Addresses,
			})
			etcdKV, err := kv.GetKV(ctx)
			if err != nil {
				t.Fatalf("failed to get kv: %v, got : %s", testCase.Addresses, err)
			}
			assert.NotNil(t, etcdKV)
		})
	}

}

func TestGetKVUntil(t *testing.T) {
	testCases := []struct {
		Addresses   []string
		DailTimeout time.Duration
	}{
		{
			Addresses:   []string{"127.0.0.1:3306"},
			DailTimeout: 3 * time.Second,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("name-%d", i), func(t *testing.T) {
			db := etcd_.NewEtcdKV(etcd_.EtcdConfig{
				Addresses: testCase.Addresses,
			}, etcd_.WithDialTimeout(testCase.DailTimeout))
			sqlDB, err := db.GetKVUntil(context.Background(), 5*time.Second, 20*time.Second)
			if err != nil {
				t.Fatalf("failed to get kv: %v, got : %s", testCase.Addresses, err)
			}
			assert.NotNil(t, sqlDB)
		})
	}

}

func TestNew(t *testing.T) {

	cfgFile := "./etcd.yaml"
	config := etcd_.NewConfig(etcd_.WithViper(viper_.GetViper(cfgFile, "discovery.etc")))

	kv, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	t.Logf("kv: %#v", kv)
}
