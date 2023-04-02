/*
 *Copyright (c) 2023, kaydxh
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
package taskqueue

import "sync"

var (
	taskMu sync.RWMutex
	m      TaskerMap
)

// Register registers the tasker to the tasker map. b.Name will be
// used as the name registered with this builder.
//
func Register(b Tasker) {
	taskMu.Lock()
	defer taskMu.Unlock()

	if b == nil {
		panic("register tasker is nil")
	}

	tasker := Get(b.Scheme())
	if tasker != nil {
		panic("double register tasker " + b.Scheme())
	}
	m.Store(b.Scheme(), b)
}

// Get returns the tasker registered with the given scheme.
//
// If no tasker is register with the scheme, nil will be returned.
func Get(scheme string) Tasker {
	taskMu.Lock()
	defer taskMu.Unlock()

	b, has := m.Load(scheme)
	if !has {
		return nil
	}

	return b
}
