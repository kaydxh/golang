package rotatefile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	os_ "github.com/kaydxh/golang/go/os"
	time_ "github.com/kaydxh/golang/go/time"
)

const (
	defaultRotateSize                   = 100 * 1024 * 1024 //100MB
	defaultrotateInterval time.Duration = time.Hour
)

type RotateFiler struct {
	file     *os.File
	filename string
	mu       sync.Mutex
	opts     struct {
		prefixName string
		subfixName string
		//maxSize int64
		//maxAge is the maximum number of days to retain old files
		maxAge time.Duration

		//rotate file when file size larger than rotateSize
		rotateSize int64
		//rotate file in rotateInterval
		rotateInterval time.Duration
	}
}

func NewRotateFiler(filename string, options ...RotateFilerOption) (*RotateFiler, error) {
	r := &RotateFiler{
		filename: filename,
	}
	//	r.opts.rotateInterval = defaultrotateInterval

	r.ApplyOptions(options...)

	return r, nil
}

func (f *RotateFiler) Write(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	out, err := f.getWriterNolock(int64(len(p)))
	if err != nil {
		return 0, err
	}

	return out.Write(p)

}

func (f *RotateFiler) generateRotateFilename() string {
	now := time.Now()
	return time_.TruncateToUTCString(now, f.opts.rotateInterval, time_.ShortDashTimeHourFormat)
}

func (f *RotateFiler) getWriterNolock(length int64) (io.Writer, error) {
	basename := f.generateRotateFilename()
	filename := filepath.Join(f.filename, basename)
	fmt.Println("filename: ", filename)

	fi, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return os_.OpenFile(filename, false)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get file info, err: %v", err)
	}

	//rotate file by size, todo
	if f.opts.rotateSize > 0 && (f.opts.rotateSize+length) <= fi.Size() {
		return os_.OpenFile(filename, true)
	}

	//file exist, and rotate file by interval
	return os_.OpenFile(filename, true)
}
