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
	strings_ "github.com/kaydxh/golang/go/strings"
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
	DefaultCheckInterval       time.Duration = time.Minute
	DefaultbaseExpired         time.Duration = 72 * time.Hour
	DefalutRandomizationFactor               = 0.1
	DefalutMultiplier                        = 0.8
	DefalutMinInterval                       = time.Minute
)

type DiskCleanerConfig struct {
	checkInterval time.Duration
	baseExpired   time.Duration
	minExpired    time.Duration
}

// DiskCleanerSerivce ...
type DiskCleanerSerivce struct {
	// path:ExponentialBackOffMap
	epoByPath time_.ExponentialBackOffMap
	exts      []string
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

// NewDiskCleanerSerivce ...
func NewDiskCleanerSerivce(
	diskUsage float32,
	paths []string,
	exts []string,
	opts ...DiskCleanerConfigOption,
) (*DiskCleanerSerivce, error) {
	canPaths, ok := checkAndCanoicalzePaths(paths...)
	if !ok {
		return nil, fmt.Errorf("invalid paths for disk cheaner")
	}

	if diskUsage < 0 {
		diskUsage = 0
	}
	if diskUsage > 100 {
		diskUsage = 100
	}

	if len(exts) == 0 {
		return nil, fmt.Errorf("invalid exts for disk cleaner")
	}

	s := &DiskCleanerSerivce{
		diskUsage: diskUsage,
		exts:      exts,
	}
	s.opts.ApplyOptions(opts...)

	if s.opts.checkInterval == 0 {
		s.opts.checkInterval = DefaultCheckInterval
	}
	if s.opts.minExpired == 0 {
		s.opts.minExpired = DefalutMinInterval
	}
	if s.opts.baseExpired == 0 {
		s.opts.baseExpired = DefaultbaseExpired
	}

	for _, path := range canPaths {
		exp := time_.NewExponentialBackOff(
			time_.WithExponentialBackOffOptionInitialInterval(s.opts.baseExpired),
			time_.WithExponentialBackOffOptionRandomizationFactor(DefalutRandomizationFactor),
			time_.WithExponentialBackOffOptionMultiplier(DefalutMultiplier),
			time_.WithExponentialBackOffOptionMaxInterval(s.opts.baseExpired),
			time_.WithExponentialBackOffOptionMinInterval(s.opts.minExpired),
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
						if strings_.SliceContains(s.exts, filepath.Ext(filePath)) {
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
					}

					return nil
				})

			} else {
				// reset expired Time
				ebo.Reset()
				logger.Infof("reset disk path: %v expired time: %v", diskPath, ebo.GetCurrentInterval())
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
