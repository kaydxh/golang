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
 */package syscall_test

import (
	"fmt"
	"testing"
	"time"

	syscall_ "github.com/kaydxh/golang/go/syscall"
	"github.com/shirou/gopsutil/mem"
)

func TestSysTotalMemory(t *testing.T) {
	total := syscall_.MemoryUsage{}.SysTotalMemory()
	t.Logf("total: %vG", total/1024/1024/1024)

	free := syscall_.MemoryUsage{}.SysFreeMemory()
	t.Logf("free: %vG", free/1024/1024/1024)

	for {
		time.Sleep(time.Second)
		v, _ := mem.VirtualMemory()

		fmt.Printf(
			"Total: %v, Free:%v, UsedPercent:%f%%, Used:%v, Available: %v\n",
			v.Total,
			v.Free,
			v.UsedPercent,
			v.Used,
			v.Available,
		)

	}

}
