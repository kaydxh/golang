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
package task_test

import (
	"fmt"
	"testing"
	"time"

	pool_ "github.com/kaydxh/golang/pkg/pool/task"
)

func TestProcessOk(t *testing.T) {

	p := pool_.New(func(task interface{}) error {
		fmt.Printf("start task[%v]\n", task)
		time.Sleep(time.Second)
		fmt.Printf("finish task[%v]\n", task)
		return nil //must return nil
	}, pool_.WithBurst(5))

	//	defer p.Wait()
	for i := 0; i < 10; i++ {
		err := p.Put(fmt.Sprintf("%d", i))
		if err != nil {
			t.Errorf("put work err: %v", err)
			return
		}
		//p.Stop()
	}
	//	p.Stop()
	p.Wait()
	t.Logf("err: %v", p.Error())
}

func TestProcessFail(t *testing.T) {

	p := pool_.New(func(task interface{}) error {
		fmt.Printf("start task[%v]\n", task)
		time.Sleep(time.Second)
		fmt.Printf("finish task[%v]\n", task)
		return fmt.Errorf("failed to process task: %v\n", task)
	}, pool_.WithBurst(2), pool_.WithErrStop(true))

	//defer p.Wait()
	for i := 0; i < 4; i++ {
		err := p.Put(fmt.Sprintf("%d", i))
		if err != nil {
			t.Errorf("put work err: %v", err)
			return
		}
	}
	p.Stop()
	p.Wait()
	t.Logf("err: %v", p.Error())
}
