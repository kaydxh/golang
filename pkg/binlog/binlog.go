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

type MessageDecoderFunc func(ctx context.Context, data []byte) (interface{}, error)
type MessageKeyDecodeFunc func(ctx context.Context, data []byte) (ds_.MessageKey, error)

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

	msgDecodeFunc    MessageDecoderFunc
	msgKeyDecodeFunc MessageKeyDecodeFunc
}

type Channel struct {
	Name string
}

type BinlogService struct {
	consumers []mq_.Consumer

	dataStore ds_.DataStore

	taskq *taskq_.Pool

	opts BinlogOptions

	inShutdown atomic.Bool
	mu         sync.Mutex
	cancel     func()
}

func defaultBinlogServiceOptions() BinlogOptions {
	opts := BinlogOptions{
		prefixName:     "segment",
		suffixName:     "log",
		flushBatchSize: 1024,
		//flushBatchSize: 5,
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
	/*
		if taskq == nil {
			return nil, fmt.Errorf("taskq is empty")
		}
	*/
	if len(consumers) == 0 {
		return nil, fmt.Errorf("consumers is empty")
	}

	bs := &BinlogService{
		dataStore: dataStore,
		taskq:     taskq,
		consumers: consumers,
		opts:      defaultBinlogServiceOptions(),
	}
	bs.ApplyOptions(opts...)

	return bs, nil
}

// send archive path to topic
func (srv *BinlogService) rotateCallback(ctx context.Context, path string) {
	logger := srv.logger()

	if srv.taskq == nil {
		return
	}

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

func (srv *BinlogService) work(ctx context.Context, msgCh <-chan *ds_.Message) {
	logger := srv.logger()
	timer := time.NewTimer(srv.opts.flushInterval)
	defer timer.Stop()

	var (
		flushBatchData []interface{}
		lastMsgKey     ds_.MessageKey
	)

	operateFunc := func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				logger.Info("start to flush timer")
				if len(flushBatchData) > 0 {
					_, err := srv.dataStore.WriteData(ctx, flushBatchData, lastMsgKey)
					if err != nil {

					}
					flushBatchData = nil
					logger.Info("finished to flush timer")
				}
				timer.Reset(srv.opts.flushInterval)
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				var msgValue interface{}
				var err error
				if srv.opts.msgDecodeFunc != nil {
					msgValue, err = srv.opts.msgDecodeFunc(ctx, msg.Value)
					if err != nil {
						break
					}

				} else {
					msgValue = msg.Value
				}

				msgKey := ds_.MessageKey{
					Key: string(msg.Key),
				}
				if srv.opts.msgKeyDecodeFunc != nil {
					msgKey, err = srv.opts.msgKeyDecodeFunc(ctx, msg.Key)
					if err != nil {
						break
					}
				}

				if msgKey.MsgType != ds_.MsgType_Insert {
					if len(flushBatchData) > 0 {
						_, err = srv.dataStore.WriteData(ctx, flushBatchData, lastMsgKey)
						flushBatchData = nil
					}
					//todo do the msg
					continue

				} else {
					// insert type
					if len(flushBatchData) == 0 || msgKey.Equual(lastMsgKey) {
						flushBatchData = append(flushBatchData, msgValue)
						lastMsgKey = msgKey

					} else {
						logger.Infof("current msg key[%v] is not equal last msg key[%v]", msgKey, lastMsgKey)
						if len(flushBatchData) > 0 {
							_, err = srv.dataStore.WriteData(ctx, flushBatchData, lastMsgKey)
							flushBatchData = nil
						}
						//todo do the msg
						continue

					}
				}

				if len(flushBatchData) >= srv.opts.flushBatchSize {
					logger.Infof("flush batch data size[%v] >= flush batch size[%v]", len(flushBatchData), srv.opts.flushBatchSize)
					_, err = srv.dataStore.WriteData(ctx, flushBatchData[:srv.opts.flushBatchSize], lastMsgKey)
					flushBatchData = flushBatchData[srv.opts.flushBatchSize:]

					if err != nil {
						break
					}

					// https://github.com/golang/go/issues/27169
					// https://tonybai.com/2016/12/21/how-to-use-timer-reset-in-golang-correctly/
					if !timer.Stop() {
						select {
						case <-timer.C:
						default:
						}
					}
					timer.Reset(srv.opts.flushInterval)
				}
			}
		}

	}

	go operateFunc()
	return
}

func (srv *BinlogService) flush(ctx context.Context, consumer mq_.Consumer) error {
	logger := srv.logger()

	msgCh := make(chan *ds_.Message)
	srv.work(ctx, msgCh)
	for msg := range consumer.ReadStream(ctx) {
		logger.Infof("recv message key[%v] value[%v] from channel[%v]", string(msg.Key()), string(msg.Value()), consumer.Topic())
		if msg.Error() != nil {
			logger.WithError(msg.Error()).Errorf("faild to read stream %v", consumer.Topic())
			continue
		}

		msgWrap := &ds_.Message{
			Key:   msg.Key(),
			Value: msg.Value(),
		}
		msgCh <- msgWrap
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
