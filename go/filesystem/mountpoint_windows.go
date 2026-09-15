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
//go:build windows

package filesystem

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// DeviceNumber on Windows is the volume serial number — the closest
// equivalent of a unix device number for mount identity.
type DeviceNumber uint64

func (num DeviceNumber) String() string {
	return fmt.Sprintf("%d", uint64(num))
}

// Mount describes one Windows volume mount point.
type Mount struct {
	Path           string
	FilesystemType string
	Device         string
	DeviceNumber   DeviceNumber
	Subtree        string
	ReadOnly       bool
}

// FindMount resolves the volume containing path via GetVolumePathName —
// the Windows moral equivalent of walking /proc/self/mountinfo: the
// volume root is the mount point.
func FindMount(path string) (*Mount, error) {
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	var volume [windows.MAX_PATH + 1]uint16
	if err := windows.GetVolumePathName(p, &volume[0], uint32(len(volume))); err != nil {
		return nil, err
	}
	vol := windows.UTF16ToString(volume[:])

	var serial uint32
	volPtr, err := windows.UTF16PtrFromString(vol)
	// A serial we cannot read degrades to 0 — the mount still answers.
	if err == nil {
		if err := windows.GetVolumeInformation(
			volPtr, nil, 0, &serial, nil, nil, nil, 0); err != nil {
			serial = 0
		}
	}

	return &Mount{
		Path:           vol,
		FilesystemType: "NTFS",
		Device:         vol,
		DeviceNumber:   DeviceNumber(serial),
		Subtree:        vol,
	}, nil
}
