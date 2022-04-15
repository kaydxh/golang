package filecleanup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kaydxh/golang/go/errors"
	errors_ "github.com/kaydxh/golang/go/errors"
	filepath_ "github.com/kaydxh/golang/go/path/filepath"
	syscall_ "github.com/kaydxh/golang/go/syscall"
	time_ "github.com/kaydxh/golang/go/time"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

var (
	workDir string
)

func init() {
	workDir, _ = os.Getwd()
}

const (
	DefaultCheckInterval time.Duration = time.Minute
	DefaultbaseExpired   time.Duration = 72 * time.Hour
)

type DiskCleanerConfig struct {
	checkInterval time.Duration
	baseExpired   time.Duration
}

// DiskCleanerSerivce ...
type DiskCleanerSerivce struct {
	epoByPath time_.ExponentialBackOffMap
	//paths     []string
	//0-100
	diskUsage  float32
	inShutdown atomic.Bool // true when when server is in shutdown

	opts DiskCleanerConfig

	mu     sync.Mutex
	cancel func()
}

func checkAndCanoicalzePaths(paths ...string) ([]string, bool) {
	var canPaths []string
	for _, path := range paths {
		absPath, err := filepath_.CanonicalizePath(path)
		if err != nil {

			fmt.Printf("err: %v\n", err)
			return nil, false
		}

		if absPath == "" || absPath == "/" || absPath == workDir {
			return nil, false
		}
		canPaths = append(canPaths, absPath)
	}

	return canPaths, true

}

func defaultExponentialBackOff() time_.ExponentialBackOff {
	epo := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionInitialInterval(DefaultbaseExpired),
		time_.WithExponentialBackOffOptionRandomizationFactor(0.1),
		time_.WithExponentialBackOffOptionMultiplier(0.8),
		time_.WithExponentialBackOffOptionMaxInterval(DefaultbaseExpired),
		time_.WithExponentialBackOffOptionMinInterval(time.Minute),
		time_.WithExponentialBackOffOptionMaxElapsedTime(0),
	)
	return *epo
}

// NewDiskCleanerSerivce ...
func NewDiskCleanerSerivce(
	diskUsage float32,
	paths []string,
	opts ...DiskCleanerConfigOption,
) (*DiskCleanerSerivce, error) {
	canPaths, ok := checkAndCanoicalzePaths(paths...)
	if !ok {
		return nil, fmt.Errorf("invalid paths for disk")
	}
	if diskUsage < 0 {
		diskUsage = 0
	}
	if diskUsage > 100 {
		diskUsage = 100
	}

	s := &DiskCleanerSerivce{
		//	paths:     canPaths,
		diskUsage: diskUsage,
	}

	s.opts.checkInterval = DefaultCheckInterval
	s.opts.baseExpired = DefaultbaseExpired

	s.opts.ApplyOptions(opts...)

	for _, path := range canPaths {
		exp := time_.NewExponentialBackOff(
			time_.WithExponentialBackOffOptionInitialInterval(s.opts.baseExpired),
			time_.WithExponentialBackOffOptionRandomizationFactor(0.1),
			time_.WithExponentialBackOffOptionMultiplier(0.8),
			time_.WithExponentialBackOffOptionMaxInterval(s.opts.baseExpired),
			time_.WithExponentialBackOffOptionMinInterval(time.Minute),
			time_.WithExponentialBackOffOptionMaxElapsedTime(0),
		)

		s.epoByPath.Store(path, *exp)
	}
	fmt.Printf("s: %+v\n", s)
	return s, nil
}

func (s *DiskCleanerSerivce) getLogger() *logrus.Entry {
	return logrus.WithField("module", "DiskCleaner")
}

// Run will initialize the backend. It must not block, but may run go routines in the background.
func (s *DiskCleanerSerivce) Run(ctx context.Context) error {
	logger := s.getLogger()
	logger.Infoln("DiskCleaner Run")
	if s.inShutdown.Load() {
		logger.Infoln("DiskCleaner Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(s.Serve(ctx))
	}()
	return nil
}

// Serve ...
func (s *DiskCleanerSerivce) Serve(ctx context.Context) error {
	logger := s.getLogger()
	logger.Infoln("DiskCleaner Serve")

	if s.inShutdown.Load() {
		logger.Infoln("DiskCleaner Shutdown")
		return fmt.Errorf("server closed")
	}

	defer s.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.mu.Lock()
	s.cancel = cancel
	s.mu.Unlock()

	time_.UntilWithContxt(ctx, func(ctx context.Context) {
		err := s.clean(ctx)
		if err != nil {
			logger.WithError(err).Errorf("failed to clean")
			return
		}
	}, s.opts.checkInterval)
	if err := ctx.Err(); err != nil {
		logger.WithError(err).Errorf("stopped checking")
		return err
	}
	logger.Info("stopped checking")
	return nil
}

func (s *DiskCleanerSerivce) clean(ctx context.Context) error {

	var (
		wg   sync.WaitGroup
		errs []error
	)

	logger := s.getLogger()
	s.epoByPath.Range(func(path string, ebo time_.ExponentialBackOff) bool {
		wg.Add(1)
		go func(diskPath string, ebo time_.ExponentialBackOff) {
			defer wg.Done()
			du, err := syscall_.NewDiskUsage(diskPath)
			if err != nil {
				s.mu.Lock()
				errs = append(errs, err)
				s.mu.Unlock()
				return
			}

			if du.Usage() >= s.diskUsage {
				//clean
				logger.Infof("disk[%v] usage over %v, start to clean", diskPath, s.diskUsage)
				actualExpired, _ := ebo.NextBackOff()
				filepath.Walk(diskPath, func(filePath string, info os.FileInfo, err error) error {

					if !info.IsDir() {
						now := time.Now()
						if now.Sub(info.ModTime()) > actualExpired {
							logger.Infof(
								"file %v expired[%v], modify time: %v, now: %v",
								filePath,
								actualExpired,
								info.ModTime(),
								now,
							)
						} else {
							logger.Infof("file %v normal[%v], modify time: %v, now: %v", filePath, actualExpired, info.ModTime(), now)
						}
					}

					return nil
				})

			} else {
				// reset expired Time
				ebo.Reset()
				logger.Infof("disk path: %v reset expired time: %v", diskPath, ebo.GetCurrentInterval())
			}
			s.epoByPath.Store(diskPath, ebo)

		}(path, ebo)

		return true
	})
	wg.Wait()
	return errors_.NewAggregate(errs)
}

// Shutdown ...
func (s *DiskCleanerSerivce) Shutdown() {
	s.inShutdown.Store(true)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		s.cancel()
	}
}
