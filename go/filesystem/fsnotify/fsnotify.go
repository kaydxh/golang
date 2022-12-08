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
	"github.com/sirupsen/logrus"
)

type FsnotifyConfig struct {
	Dirs []string
}

type EventCallbackFunc func()

type FsnotifyOptions struct {
	WriteCallbackFunc EventCallbackFunc
}

type FsnotifyService struct {
	watcher *fsnotify.Watcher
	conf    FsnotifyConfig

	opts       FsnotifyOptions
	inShutdown atomic.Bool
	mu         sync.Mutex
	cancel     func()
}

func NewFsnotifyService(config FsnotifyConfig) (*FsnotifyService, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if len(config.Dirs) == 0 {
		return nil, fmt.Errorf("dirs is empty")
	}

	fs := &FsnotifyService{
		watcher: watcher,
		conf:    config,
	}

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

	err := srv.AddWatchDirs(srv.conf.Dirs...)
	if err != nil {
		logger.WithError(err).Errorf("failed to add watcher for dirs: %v", srv.conf.Dirs)
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
			}
			if ev.Op&fsnotify.Write != 0 {
				logger.Infof("%s happen write event", ev.Name)
			}
			if ev.Op&fsnotify.Remove != 0 {
				logger.Infof("%s happen remove event", ev.Name)
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
func (srv *FsnotifyService) AddWatchDir(path string) error {
	logger := srv.logger()
	return filepath.Walk(path, func(walkPath string, fi os.FileInfo, err error) error {
		if err != nil {
			logger.WithError(err).Errorf("failed to walk dir: %v", walkPath)
			return err
		}

		if fi.IsDir() {
			if err = srv.watcher.Add(walkPath); err != nil {
				logger.WithError(err).Errorf("failed to add watcher for dir: %v", walkPath)
				return err
			}
			logger.Infof("add watcher for dir: %v", walkPath)
		}

		return nil
	})
}

func (srv *FsnotifyService) AddWatchDirs(paths ...string) error {
	for _, path := range paths {
		err := srv.AddWatchDir(path)
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
