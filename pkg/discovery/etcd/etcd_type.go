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
	"sync/atomic"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ETCDKV atomic.Value

//check type ETCDKV* whether to implement interface of atomic.Value
//var _ atomic.Value = (*ETCDKV)(nil)

//check type ETCDKV whether to implement interface of atomic.Value
//var _ atomic.Value = ETCDKV{}

func (m *ETCDKV) Store(value *clientv3.Client) {
	(*atomic.Value)(m).Store(value)
}

func (m *ETCDKV) Load() *clientv3.Client {
	value := (*atomic.Value)(m).Load()
	if value == nil {
		return nil
	}
	return value.(*clientv3.Client)
}
