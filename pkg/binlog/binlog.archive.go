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
package binlog

import (
	"context"
	"encoding/json"

	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"

	io_ "github.com/kaydxh/golang/go/io"
	s3_ "github.com/kaydxh/golang/pkg/storage/s3"
	"github.com/sirupsen/logrus"
)

const (
	ArchiveTaskScheme = "ArchiveTaskTask"
)

type ArchiveTaskArgs struct {
	LocalFilePath  string
	RemoteRootPath string
}

type ArchiveTask struct {
	bucket *s3_.Storage
}

func (t ArchiveTask) Scheme() string {
	return ArchiveTaskScheme
}

func (t ArchiveTask) TaskHandler(ctx context.Context, msg *queue_.Message) (*queue_.MessageResult, error) {
	logger := logrus.WithField("message_id", msg.Id).
		WithField("message_inner_id", msg.InnerId).
		WithField("module", "TaskHandler")

	result := &queue_.MessageResult{
		Id:      msg.Id,
		InnerId: msg.InnerId,
		Name:    msg.Name,
		Scheme:  msg.Scheme,
	}

	var args ArchiveTaskArgs
	err := json.Unmarshal([]byte(msg.Args), &args)
	if err != nil {
		logger.WithError(err).Errorf("failed to unmarshal msg: %v", msg)
		return result, err
	}
	data, err := io_.ReadFile(args.LocalFilePath)
	if err != nil {
		logger.WithError(err).Errorf("failed to read binlog[%v]", args.LocalFilePath)
		return result, err
	}
	if len(data) == 0 {
		logger.Infof("binlog[%v] is empty, not need to archive", args.LocalFilePath)
		return result, nil
	}

	s3Path := args.RemoteRootPath
	err = t.bucket.WriteAll(ctx, s3Path, data, nil)
	if err != nil {
		logger.WithError(err).Errorf("failed to upload data size[%v] to s3[%v]", len(data), s3Path)
		return result, err
	}

	return result, nil
}
