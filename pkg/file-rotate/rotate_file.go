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
	defaultRotateSize = 100 * 1024 * 1024 //100MB
//	defaultrotateInterval time.Duration = time.Hour
//defaultSuffixName     string        = ".log"
)

type RotateFiler struct {
	file    *os.File
	filedir string
	mu      sync.Mutex
	opts    struct {
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

func NewRotateFiler(filedir string, options ...RotateFilerOption) (*RotateFiler, error) {
	r := &RotateFiler{
		filedir: filedir,
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
	if f.opts.rotateInterval > 0 {
		now := time.Now()
		return time_.TruncateToUTCString(now, f.opts.rotateInterval, time_.ShortTimeFormat)
	}
	return ""
}

func (f *RotateFiler) getWriterNolock(length int64) (io.Writer, error) {
	basename := f.generateRotateFilename()
	filename := f.opts.prefixName + basename + f.opts.subfixName
	if filename == "" {
		filename = "default.log"
	}
	filepath := filepath.Join(f.filedir, filename)

	var err error
	rotated := false

	fi, err := os.Stat(filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to get file info, err: %v", err)
		}
		//file is not exist, think just like rotating file
		rotated = true
	}

	//rotate file by size
	if err == nil && f.opts.rotateSize > 0 && (fi.Size()+length) > f.opts.rotateSize {
		/*
			fmt.Printf(
				"generateNextSeqFilename ffi.Size()+length: %v, rotateSize: %v\n",
				fi.Size()+length,
				f.opts.rotateSize,
			)
		*/
		filepath, err = f.generateNextSeqFilename(filepath)
		if err != nil {
			return nil, fmt.Errorf("failed to generate rotate file name by seq, err: %v", err)
		}

		rotated = true
	}

	if rotated {
		fn, err := os_.OpenFile(filepath, true)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %v, err: %v", filepath, err)
		}

		if f.file != nil {
			f.file.Close()
		}
		f.file = fn
	}

	return f.file, nil
}

//filename like foo foo.1 foo.2 ...
func (f *RotateFiler) generateNextSeqFilename(filepath string) (string, error) {

	var newFilepath string
	seq := 0

	for {
		if seq == 0 {
			newFilepath = filepath
		} else {
			newFilepath = fmt.Sprintf("%s.%d", filepath, seq)
		}

		_, err := os.Stat(newFilepath)
		if os.IsNotExist(err) {
			return newFilepath, nil
		}
		if err != nil {
			return "", err
		}
		//file exist, need to get next seq filename
		seq++
	}

}
