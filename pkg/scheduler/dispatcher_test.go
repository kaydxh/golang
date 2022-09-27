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
package scheduler_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	scheduler_ "github.com/kaydxh/golang/pkg/scheduler"
	task_ "github.com/kaydxh/golang/pkg/scheduler/task"
)

func TestDispatcher(t *testing.T) {
	burst := 1
	//	broker := memory_.NewBroker()
	dispatcher := scheduler_.NewDispatcher(burst)
	dispatcher.AddTask(&task_.Task{
		TaskId: uuid.New().String(),
		TaskFunc: func(a int, b string) error {
			fmt.Printf("do task: %v:%v\n", a, b)
			return nil
		},
		Args: []task_.TaskArgument{
			task_.TaskArgument{
				Type:  "int",
				Value: 8,
			},
			task_.TaskArgument{
				Type:  "int",
				Value: "hello",
			},
		},
	})
	dispatcher.Stop()
	fmt.Println("out")

}
