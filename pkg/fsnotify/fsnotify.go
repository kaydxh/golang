package fsnotify

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/kaydxh/golang/go/errors"
	os_ "github.com/kaydxh/golang/go/os"
	"github.com/sirupsen/logrus"
)

type EventCallbackFunc func(ctx context.Context, path string)

type FsnotifyOptions struct {
	CreateCallbackFunc EventCallbackFunc
	WriteCallbackFunc  EventCallbackFunc
	RemoveCallbackFunc EventCallbackFunc
}

type FsnotifyService struct {
	watcher *fsnotify.Watcher
	paths   []string

	opts       FsnotifyOptions
	inShutdown atomic.Bool
	mu         sync.Mutex
	cancel     func()
}

// paths can also be dir or file or both of them
func NewFsnotifyService(paths []string, opts ...FsnotifyOption) (*FsnotifyService, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("paths is empty")
	}

	for _, path := range paths {
		ok, err := os_.PathExist(path)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("%s is not exist", path)
		}
	}

	fs := &FsnotifyService{
		watcher: watcher,
		paths:   paths,
	}
	fs.ApplyOptions(opts...)

	return fs, nil
}

func (srv *FsnotifyService) logger() logrus.FieldLogger {
	return logrus.WithField("module", "FsnotifyService")
}

func (srv *FsnotifyService) Run(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("FsnotifyService Run")
	if srv.inShutdown.Load() {
		logger.Infoln("FsnotifyService Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(srv.Serve(ctx))
	}()
	return nil
}

func (srv *FsnotifyService) Serve(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("FsnotifyService Serve")

	if srv.inShutdown.Load() {
		err := fmt.Errorf("server closed")
		logger.WithError(err).Errorf("FsnotifyService Serve canceled")
		return err
	}

	defer srv.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	srv.mu.Lock()
	srv.cancel = cancel
	srv.mu.Unlock()

	defer func() {
		err := srv.watcher.Close()
		if err != nil {
			logger.WithError(err).Errorf("failed to close FsnotifyService watcher")
		}
	}()

	err := srv.AddWatchPaths(false, srv.paths...)
	if err != nil {
		logger.WithError(err).Errorf("failed to add watcher for path: %v", srv.paths)
		return err
	}

	for {
		select {
		case ev, ok := <-srv.watcher.Events:
			if !ok {
				continue
			}

			if ev.Op&fsnotify.Create != 0 {
				logger.Infof("%s happen create event", ev.Name)
				srv.AddWatchPaths(false, ev.Name)
				if srv.opts.CreateCallbackFunc != nil {
					srv.opts.CreateCallbackFunc(ctx, ev.Name)
				}
			}
			if ev.Op&fsnotify.Write != 0 {
				logger.Infof("%s happen write event", ev.Name)
				srv.AddWatchPaths(false, ev.Name)
				if srv.opts.WriteCallbackFunc != nil {
					srv.opts.WriteCallbackFunc(ctx, ev.Name)
				}
			}
			if ev.Op&fsnotify.Remove != 0 {
				logger.Infof("%s happen remove event", ev.Name)
				srv.AddWatchPaths(true, ev.Name)
				if srv.opts.RemoveCallbackFunc != nil {
					srv.opts.RemoveCallbackFunc(ctx, ev.Name)
				}
			}
			if ev.Op&fsnotify.Rename != 0 {
				logger.Infof("%s happen rename event", ev.Name)
			}
			if ev.Op&fsnotify.Chmod != 0 {
				logger.Infof("%s happen chmod event", ev.Name)
			}

		case <-ctx.Done():
			logger.WithError(ctx.Err()).Errorf("server canceld")
			return ctx.Err()
		}
	}
}

// Add starts watching the named directory (support recursively).
func (srv *FsnotifyService) AddWatchPath(unWatch bool, path string) error {
	logger := srv.logger()

	ok, err := os_.IsDir(path)
	if err != nil {
		return err
	}

	if ok {
		return filepath.Walk(path, func(walkPath string, fi os.FileInfo, err error) error {
			if err != nil {
				logger.WithError(err).Errorf("failed to walk dir: %v", walkPath)
				return err
			}

			if fi.IsDir() {
				return srv.Add(unWatch, walkPath)
			}
			return nil
		})
	} else {
		return srv.Add(unWatch, path)
	}
}

// Add starts watching the named directory (non-recursively).
func (srv *FsnotifyService) Add(unWatch bool, path string) (err error) {
	logger := srv.logger()
	if unWatch {
		err = srv.watcher.Remove(path)
	} else {
		err = srv.watcher.Add(path)
	}
	if err != nil {
		logger.WithError(err).Errorf("failed to add watcher for path: %v, ", path)
		return err
	}
	logger.Infof("add watcher for path: %v", path)

	return nil
}

func (srv *FsnotifyService) AddWatchPaths(unWatch bool, paths ...string) error {
	for _, path := range paths {
		err := srv.AddWatchPath(unWatch, path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *FsnotifyService) Shutdown() {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.cancel != nil {
		srv.cancel()
	}
}
