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
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/kaydxh/golang/go/errors"
	ds_ "github.com/kaydxh/golang/pkg/binlog/datastore"
	mq_ "github.com/kaydxh/golang/pkg/mq"
	taskq_ "github.com/kaydxh/golang/pkg/pool/taskqueue"
	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type BinlogOptions struct {
	rootPath       string
	prefixName     string
	suffixName     string
	flushBatchSize int
	flushInterval  time.Duration
	rotateInterval time.Duration
	rotateSize     int64

	remotePrefixPath string
	archive          bool
}

type Channel struct {
	Name string
}

type BinlogService struct {
	consumers []mq_.Consumer

	//rotateFiler *rotate_.RotateFiler
	dataStore ds_.DataStore
	//rotateFilers   map[string]*rotate_.RotateFiler //message key -> rotateFilter
	//dataStores     map[string]DataStore //message key -> rotateFilter
	//rotateFilersMu sync.Mutex

	taskq *taskq_.Pool

	opts BinlogOptions

	inShutdown atomic.Bool
	mu         sync.Mutex
	cancel     func()
}

func defaultBinlogServiceOptions() BinlogOptions {
	opts := BinlogOptions{
		prefixName: "segment",
		suffixName: "log",
		//flushBatchSize: 1024,
		flushBatchSize: 1,
		flushInterval:  time.Second, // 1s
		rotateInterval: time.Hour,
		rotateSize:     100 * 1024 * 1024, //100M
	}
	path, err := os.Getwd()
	if err != nil {
		path = "/"
	}
	opts.rootPath = path
	return opts
}

func NewBinlogService(dataStore ds_.DataStore, taskq *taskq_.Pool, consumers []mq_.Consumer, opts ...BinlogServiceOption) (*BinlogService, error) {
	if taskq == nil {
		return nil, fmt.Errorf("taskq is empty")
	}
	if len(consumers) == 0 {
		return nil, fmt.Errorf("consumers is empty")
	}

	bs := &BinlogService{
		dataStore: dataStore,
		//rotateFilers: make(map[string]*rotate_.RotateFiler, 0),
		//dataStores: make(map[string]DataStore, 0),
		taskq:     taskq,
		consumers: consumers,
		opts:      defaultBinlogServiceOptions(),
	}
	bs.ApplyOptions(opts...)

	/*
		rotateFiler, _ := rotate_.NewRotateFiler(
			bs.opts.rootPath,
			rotate_.WithRotateSize(bs.opts.rotateSize),
			rotate_.WithRotateInterval(bs.opts.rotateInterval),
			rotate_.WithSuffixName(bs.opts.suffixName),
			rotate_.WithPrefixName(bs.opts.prefixName),
			rotate_.WithRotateCallback(bs.rotateCallback),
		)
		bs.rotateFiler = rotateFiler
	*/

	return bs, nil
}

func (srv *BinlogService) rotateCallback(ctx context.Context, path string) {
	logger := srv.logger()

	if srv.opts.archive {
		args := &ArchiveTaskArgs{
			LocalFilePath:  path,
			RemoteRootPath: filepath.Join(srv.opts.remotePrefixPath, filepath.Base(path)),
		}
		data, err := json.Marshal(args)
		if err != nil {
			logger.WithError(err).Errorf("failed to marshal args: %v", args)
			return
		}

		id := uuid.NewString()
		msg := &queue_.Message{
			Id:     id,
			Name:   id,
			Scheme: ArchiveTaskScheme,
			Args:   string(data),
		}
		_, err = srv.taskq.Publish(ctx, msg)
		if err != nil {
			logger.WithError(err).Errorf("failed to publish msg: %v", msg)
			return
		}
	}
}

func (srv *BinlogService) logger() logrus.FieldLogger {
	return logrus.WithField("module", "BinlogService")
}

func (srv *BinlogService) Run(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("BinlogService Run")
	if srv.inShutdown.Load() {
		logger.Infoln("BinlogService Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(srv.Serve(ctx))
	}()
	return nil
}

/*
func (srv *BinlogService) getOrCreateRotateFilers(ctx context.Context, key string) *rotate_.RotateFiler {
	if key == "" {
		return srv.rotateFiler
	}

	srv.rotateFilersMu.Lock()
	defer srv.rotateFilersMu.Unlock()
	rotateFiler, ok := srv.rotateFilers[key]
	if !ok {
		rotateFiler, _ = rotate_.NewRotateFiler(
			filepath.Join(srv.opts.rootPath, key),
			rotate_.WithRotateSize(srv.opts.rotateSize),
			rotate_.WithRotateInterval(srv.opts.rotateInterval),
			rotate_.WithSuffixName(srv.opts.suffixName),
			rotate_.WithPrefixName(srv.opts.prefixName),
			rotate_.WithRotateCallback(srv.rotateCallback),
		)
		srv.rotateFilers[key] = rotateFiler
	}

	return rotateFiler
}
*/

func (srv *BinlogService) flush(ctx context.Context, consumer mq_.Consumer) error {
	logger := srv.logger()
	timer := time.NewTicker(srv.opts.flushInterval)
	defer timer.Stop()

	var flushBatchData [][]byte
	for msg := range consumer.ReadStream(ctx) {
		if msg.Error() != nil {
			logger.WithError(msg.Error()).Errorf("faild to read stream %v", consumer.Channel())
			continue
		}

		//rotateFiler := srv.getOrCreateRotateFilers(ctx, string(msg.Key()))
		flushFunc := func(ctx context.Context, data []byte) (err error) {
			flushBatchData = append(flushBatchData, data)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				if len(flushBatchData) > 0 {
					_, _, err = srv.dataStore.WriteData(ctx, "", nil, string(msg.Key()), flushBatchData)
					flushBatchData = nil
				}
			default:
				if len(flushBatchData) >= srv.opts.flushBatchSize {
					_, _, err = srv.dataStore.WriteData(ctx, "", nil, string(msg.Key()), flushBatchData)
					flushBatchData = nil
					return err
				}
				if err != nil {
					return err
				}
			}
			return nil
		}
		err := flushFunc(ctx, msg.Value())
		if err != nil {
			logger.WithError(msg.Error()).Errorf("faild to flush message for channel[%v]", consumer.Channel())
			return err
		}
	}

	return nil
}

func (srv *BinlogService) Serve(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("ServiceBinlog Serve")

	if srv.inShutdown.Load() {
		err := fmt.Errorf("server closed")
		logger.WithError(err).Errorf("ServiceBinlog Serve canceled")
		return err
	}

	defer srv.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	srv.mu.Lock()
	srv.cancel = cancel
	srv.mu.Unlock()

	g, gCtx := errgroup.WithContext(ctx)
	for _, consumer := range srv.consumers {
		consumer := consumer
		g.Go(func() error {
			return srv.flush(gCtx, consumer)
		})
	}
	err := g.Wait()
	if err != nil {
		logger.WithError(err).Errorf("wait flush worker")
		return err
	}
	logger.Info("stopped binlog service")
	return nil
}

func (srv *BinlogService) Shutdown() {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	defer srv.mu.Unlock()

	for _, consumer := range srv.consumers {
		if consumer != nil {
			consumer.Close()
		}
	}
	if srv.cancel != nil {
		srv.cancel()
	}
}
