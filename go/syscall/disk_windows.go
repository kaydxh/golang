//go:build windows

/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
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
package syscall

import (
	"golang.org/x/sys/windows"
)

// DiskUsage reports one volume's usage on Windows via GetDiskFreeSpaceEx —
// the moral equivalent of statfs(2): total, free, and free-to-the-caller
// in one call.
type DiskUsage struct {
	total uint64
	free  uint64
	avail uint64
}

// NewDiskUsage queries the volume containing path. The path may be any
// form Windows accepts (drive letter, UNC, relative).
func NewDiskUsage(path string) (*DiskUsage, error) {
	var (
		total, free, avail uint64
	)
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	if err := windows.GetDiskFreeSpaceEx(p, &avail, &total, &free); err != nil {
		return nil, err
	}
	return &DiskUsage{total: total, free: free, avail: avail}, nil
}

// Free returns total free bytes on the volume.
func (du *DiskUsage) Free() uint64 { return du.free }

// Avail returns the bytes available to the calling user (quotas applied).
func (du *DiskUsage) Avail() uint64 { return du.avail }

// Size returns total size of the volume.
func (du *DiskUsage) Size() uint64 { return du.total }

// Used returns total used bytes on the volume.
func (du *DiskUsage) Used() uint64 { return du.Size() - du.Free() }

// Usage returns the used percentage, matching the unix implementation's
// coreutils rounding (used against used+avail, rounded up on remainder).
func (du *DiskUsage) Usage() float32 {
	var deta float32
	u100 := du.Used() * 100
	nonrootTotal := du.Used() + du.Avail()
	if nonrootTotal != 0 && u100%nonrootTotal != 0 {
		deta = 1.0
	}
	if nonrootTotal == 0 {
		return 0
	}
	return float32(u100)/float32(nonrootTotal) + deta
}
