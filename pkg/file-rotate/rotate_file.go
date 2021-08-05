package rotatefile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	os_ "github.com/kaydxh/golang/go/os"
	time_ "github.com/kaydxh/golang/go/time"
	cleanup_ "github.com/kaydxh/golang/pkg/file-cleanup"
)

type RotateFiler struct {
	file     *os.File
	filedir  string
	linkpath string
	seq      int64
	mu       sync.Mutex
	opts     struct {
		prefixName     string
		fileTimeLayout string //default "20060102150405" ,take effect if rotateInterval  > 0

		subfixName string
		//maxAge is the maximum number of time to retain old files, 0 is unlimited
		maxAge time.Duration
		//maxCount is the maximum number to retain old files, 0 is unlimited
		maxCount int64

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
	r.ApplyOptions(options...)

	if r.linkpath == "" {
		r.linkpath = filepath.Base(os.Args[0]) + ".log"
	}

	// if need rotate file with rotateInterval, set default timelayout
	if r.opts.rotateInterval > 0 {
		if r.opts.fileTimeLayout == "" {
			r.opts.fileTimeLayout = time_.ShortTimeFormat
		}
	}

	return r, nil
}

// /data/log/1%%%AA20160304 -> /data/log/1*A20160304*
func globFromFileTimeLayout(filepath string) string {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`%[%+A-Za-z]`),
		regexp.MustCompile(`\*+`),
	}

	for _, re := range regexps {
		filepath = re.ReplaceAllString(filepath, "*")
	}
	return filepath + "*"
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
		return time_.TruncateToUTCString(now, f.opts.rotateInterval, f.opts.fileTimeLayout)
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
	globPath := filepath

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
		filepath, err = f.generateNextSeqFilename(filepath)
		if err != nil {
			return nil, fmt.Errorf("failed to generate rotate file name by seq, err: %v", err)
		}

		rotated = true
	}

	if f.file == nil || rotated {
		fn, err := os_.OpenFile(filepath, true)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %v, err: %v", filepath, err)
		}

		if f.file != nil {
			f.file.Close()
		}
		f.file = fn

		os_.SymLink(filepath, f.linkpath)

		globFile := globFromFileTimeLayout(globPath)

		go cleanup_.FileCleanup(globFile, cleanup_.WithMaxAge(f.opts.maxAge), cleanup_.WithMaxCount(f.opts.maxCount))
	}

	return f.file, nil
}

//filename like foo foo.1 foo.2 ...
func (f *RotateFiler) generateNextSeqFilename(filepath string) (string, error) {

	var newFilepath string
	seq := f.seq

	for {
		if seq == 0 {
			newFilepath = filepath
		} else {
			newFilepath = fmt.Sprintf("%s.%d", filepath, seq)
		}

		_, err := os.Stat(newFilepath)
		if os.IsNotExist(err) {
			f.seq = seq
			return newFilepath, nil
		}
		if err != nil {
			return "", err
		}
		//file exist, need to get next seq filename
		seq++
	}

}
