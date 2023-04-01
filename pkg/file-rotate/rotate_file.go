/*
 *Copyright (c) 2022, kaydxh
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
package rotatefile

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	os_ "github.com/kaydxh/golang/go/os"
	filepath_ "github.com/kaydxh/golang/go/path/filepath"
	time_ "github.com/kaydxh/golang/go/time"
	cleanup_ "github.com/kaydxh/golang/pkg/file-cleanup"
)

type EventCallbackFunc func(ctx context.Context, path string)

type RotateFiler struct {
	file        *os.File
	filedir     string
	curFilepath string
	seq         uint64
	linkpath    string
	mu          sync.Mutex
	opts        struct {
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
		rotateInterval     time.Duration
		syncInterval       time.Duration
		rotateCallbackFunc EventCallbackFunc
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

	if r.opts.rotateCallbackFunc != nil {
		if r.opts.syncInterval == 0 {
			r.opts.syncInterval = 30 * time.Second
		}
		go r.watch()
	}

	return r, nil
}

// /data/log/1%%%AA20160304 -> /data/log/1*A20160304*
func globFromFileTimeLayout(filePath string) string {
	regexps := []*regexp.Regexp{
		regexp.MustCompile(`%[%+A-Za-z]`),
		regexp.MustCompile(`\*+`),
	}

	for _, re := range regexps {
		filePath = re.ReplaceAllString(filePath, "*")
	}
	return filePath + "*"
}

func (f *RotateFiler) Write(p []byte) (file *os.File, n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	out, err := f.getWriterNolock(int64(len(p)))
	if err != nil {
		return nil, 0, err
	}

	n, err = out.Write(p)
	return f.file, n, err
}

func (f *RotateFiler) WriteBytesLine(p [][]byte) (file *os.File, n int, err error) {

	var data []byte
	for _, d := range p {
		data = append(data, d...)
		data = append(data, '\n')
	}
	return f.Write(data)
}

func (f *RotateFiler) generateRotateFilename() string {
	if f.opts.rotateInterval > 0 {
		now := time.Now()
		return time_.TruncateToUTCString(now, f.opts.rotateInterval, f.opts.fileTimeLayout)
	}
	return ""
}

func (f *RotateFiler) watch() {
	timer := time.NewTicker(f.opts.syncInterval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			func() {
				f.mu.Lock()
				defer f.mu.Unlock()
				f.getWriterNolock(0)
			}()
		}
	}
}

func (f *RotateFiler) getWriterNolock(length int64) (io.Writer, error) {
	basename := f.generateRotateFilename()
	filename := f.opts.prefixName + basename + f.opts.subfixName
	if filename == "" {
		filename = "default.log"
	}

	// first rotate log file, maybe /data/logs/logs.test20210917230000.log
	filePath := filepath.Join(f.filedir, filename)
	globPath := filepath.Join(filepath.Dir(filePath), f.opts.prefixName)

	// current log file, maybe /data/logs/logs.test20210917230000.log.1
	if f.curFilepath == "" {
		f.curFilepath, _ = f.getCurSeqFilename(globPath)
		f.seq = f.extractSeq(f.curFilepath)
	}

	// if curFilePath is different rotated time with filename, need reset curFilePath
	if !strings.Contains(f.curFilepath, filename) {
		f.curFilepath = filePath
		f.seq = 0
	}

	rotated := false

	fi, err := os.Stat(f.curFilepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to get file info, err: %v", err)
		}
		//file is not exist, think just like rotating file
		rotated = true
	}

	//rotate file by size
	if err == nil && f.opts.rotateSize > 0 && (fi.Size()+length) > f.opts.rotateSize {

		f.curFilepath, err = f.generateNextSeqFilename(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to generate rotate file name by seq, err: %v", err)
		}

		rotated = true
	}

	if f.file == nil || rotated {
		fn, err := os_.OpenFile(f.curFilepath, true)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %v, err: %v", f.curFilepath, err)
		}

		if f.file != nil {
			//callback
			if f.opts.rotateCallbackFunc != nil {
				f.opts.rotateCallbackFunc(context.Background(), f.file.Name())
			}
			f.file.Close()
		}
		f.file = fn

		f.seq = f.extractSeq(f.curFilepath)

		os_.SymLink(f.curFilepath, f.linkpath)

		globFile := globFromFileTimeLayout(globPath)

		go cleanup_.FileCleanup(globFile, cleanup_.WithMaxAge(f.opts.maxAge), cleanup_.WithMaxCount(f.opts.maxCount))
	}

	return f.file, nil
}

//filename like foo foo.1 foo.2 ...
func (f *RotateFiler) generateNextSeqFilename(filePath string) (string, error) {

	var newFilePath string
	seq := f.seq

	for {
		if seq == 0 {
			newFilePath = filePath
		} else {
			newFilePath = fmt.Sprintf("%s.%d", filePath, seq)
		}

		_, err := os.Stat(newFilePath)
		if os.IsNotExist(err) {
			f.seq = seq
			return newFilePath, nil
		}
		if err != nil {
			return "", err
		}
		//file exist, need to get next seq filename
		seq++
	}

}

// globPath: log/logs.test
// globFile: [log/logs.test20211008081908.log log/logs.test20211008081908.log.1 log/logs.test20211008081908.log.2]
func (f *RotateFiler) getCurSeqFilename(globPath string) (string, error) {

	globFile := globFromFileTimeLayout(globPath)
	matches, err := filepath_.Glob(globFile)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return globPath, nil
	}

	sort.Sort(cleanup_.RotatedFiles(matches))
	return matches[len(matches)-1], nil
}

func (f *RotateFiler) extractSeq(filePath string) uint64 {
	if filePath == "" {
		return 0
	}

	ext := filepath.Ext(filePath)
	if ext == "" {
		return 0
	}

	seq, err := strconv.ParseUint(ext[1:], 10, 64)
	if err != nil {
		return 0
	}

	return seq

}
