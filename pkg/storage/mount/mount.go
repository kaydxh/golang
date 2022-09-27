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
package mount

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	io_ "github.com/kaydxh/golang/go/io"

	os_ "github.com/kaydxh/golang/go/os"
	exec_ "github.com/kaydxh/golang/go/os/exec"
)

const (
	mountLabelFileName = "mount-label-file-75248"
)

func MountCeph(
	prefixMountPath, userName, address, keyring, workDir string, forceMount bool, timeout int,
) (mountPoint string, err error) {
	keyringFilePath, err := genCephKeyringConfFile(prefixMountPath, userName, keyring)
	if err != nil {
		return "", err
	}

	workDir = filepath.Join("/", workDir)

	cmd := fmt.Sprintf(
		`ceph-fuse -m %s -k %s -n client.%s -o rw.nonempty -r %s %s`,
		address,
		keyringFilePath,
		userName,
		workDir,
		mountPoint,
	)
	return doMountCmd(prefixMountPath, workDir, cmd, forceMount, timeout)
}

func MountCfs(
	prefixMountPath, address, workDir string,
	forceMount bool,
	timeout int,
) (mountPoint string, err error) {

	mountPoint = filepath.Join(prefixMountPath, workDir)
	workDir = filepath.Join("/", workDir)

	cmd := fmt.Sprintf(
		`mount -t nfs -o vers=4 %s:%s %s`,
		address,
		workDir,
		mountPoint,
	)
	return doMountCmd(prefixMountPath, workDir, cmd, forceMount, timeout)
}

func doMountCmd(
	prefixMountPath, workDir, cmd string,
	forceMount bool,
	timeout int,
) (mountPoint string, err error) {

	mountPoint = filepath.Join("/", prefixMountPath, workDir)
	if !forceMount {
		labelFilePath := filepath.Join(mountPoint, mountLabelFileName)
		exist, err := os_.PathExist(labelFilePath)
		if err == nil && exist {
			return mountPoint, nil
		}
	}

	err = os.MkdirAll(mountPoint, 0755)
	if err != nil {
		return "", err
	}

	_, _, err = exec_.Exec(time.Duration(timeout), cmd)
	if err != nil {
		return "", err
	}

	return mountPoint, nil
}

func genCephKeyringConfFile(prefixMountPath, userName, keyring string) (string, error) {
	keyringFileName := userName + ".keyring"
	keyringFilePath := filepath.Join(prefixMountPath, keyringFileName)

	line1 := "[client." + userName + "]"
	line2 := fmt.Sprintf("    key =  %s\n", keyring)
	lines := []string{line1, line2}

	err := io_.WriteFileLines(keyringFilePath, lines, false)
	if err != nil {
		return "", err
	}

	return keyringFilePath, nil
}
