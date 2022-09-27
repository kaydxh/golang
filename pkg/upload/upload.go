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
package upload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	md5_ "github.com/kaydxh/golang/go/crypto/md5"
	io_ "github.com/kaydxh/golang/go/io"
	"github.com/kaydxh/golang/go/sync/atomic"
	atomic_ "github.com/kaydxh/golang/go/sync/atomic"
)

type UploadPartInput struct {
	PartId uint32
	Offset int64
	Length int64
	Md5Sum string
}

const tmpFileSuffix = "_.tmp"

func UploadFile(
	r io.Reader,
	filePath string,
	md5Sum string,
) error {
	if filePath == "" {
		return fmt.Errorf("invalid filePath")
	}

	//check path
	absFilePath, err := filepath.Abs(filePath)
	if strings.Contains(absFilePath, "..") {
		err = fmt.Errorf("invalid filePath: %v", absFilePath)
		return err
	}

	var mu atomic_.FileLock = atomic.FileLock(filepath.Join("./", absFilePath))
	err = mu.TryLock()
	if err != nil {
		return err
	}

	unlockFun := func() error {
		err = mu.TryUnLock()
		if err != nil {
			return err
		}
		return nil
	}
	defer unlockFun()

	if md5Sum != "" {
		sum, err := md5_.SumReader(r)
		if err != nil {
			return err
		}
		gotSum := strings.ToLower(sum)
		expectSum := strings.ToLower(md5Sum)

		if gotSum != expectSum {
			return fmt.Errorf(
				"failed to check md5Sum, got: %v, expect: %v",
				gotSum,
				expectSum,
			)
		}

	}

	err = io_.WriteReader(filePath, r)
	if err != nil {
		return err
	}

	return nil
}

func UploadMultipart(
	r io.Reader,
	partInput *UploadPartInput,
	filePath string,
) error {
	if partInput == nil {
		return fmt.Errorf("invalid partInput")
	}
	if filePath == "" {
		return fmt.Errorf("invalid filePath")
	}

	//check path
	absFilePath, err := filepath.Abs(filePath)
	if strings.Contains(absFilePath, "..") {
		err = fmt.Errorf("invalid filePath: %v", absFilePath)
		return err
	}

	var mu atomic_.FileLock = atomic.FileLock(filepath.Join("./", absFilePath))
	err = mu.TryLock()
	if err != nil {
		return err
	}

	unlockFun := func() error {
		err = mu.TryUnLock()
		if err != nil {
			return err
		}
		return nil
	}
	defer unlockFun()

	if partInput.Md5Sum != "" {
		sum, err := md5_.SumReader(r)
		if err != nil {
			return err
		}
		gotSum := strings.ToLower(sum)
		expectSum := strings.ToLower(partInput.Md5Sum)

		if gotSum != expectSum {
			return fmt.Errorf(
				"failed to check md5Sum, got: %v, expect: %v",
				gotSum,
				expectSum,
			)
		}

	}

	tmpFilePath := filePath + tmpFileSuffix
	err = io_.WriteReaderAt(
		tmpFilePath,
		r,
		partInput.Offset,
		partInput.Length,
	)
	if err != nil {
		return err
	}

	return nil
}

func CompleteMultipartUpload(
	filePath string,
	md5Sum string) error {
	if filePath == "" {
		return fmt.Errorf("invalid filePath")
	}

	tmpFilePath := filePath + tmpFileSuffix
	if md5Sum != "" {
		sum, err := md5_.SumFile(tmpFilePath)
		if err != nil {
			return err
		}

		gotSum := strings.ToLower(sum)
		expectSum := strings.ToLower(md5Sum)

		if gotSum != expectSum {
			return fmt.Errorf(
				"failed to check md5Sum, got: %v, expect: %v",
				gotSum,
				expectSum,
			)
		}
	}

	_, err := os.Stat(tmpFilePath)
	if err != nil {
		return fmt.Errorf("file: %v is not existed", tmpFilePath)
	}

	err = os.Rename(tmpFilePath, filePath)
	if err != nil {
		return fmt.Errorf("failed to Rename: %v to %v", tmpFilePath, filePath)
	}

	return nil
}
