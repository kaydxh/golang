package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/kaydxh/golang/go/archive/option"
	os_ "github.com/kaydxh/golang/go/os"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type COPY_TYPE int

type Zip struct {
}

func (z *Zip) Extract(srcFile, destDir string) ([]*option.FileInfo, error) {
	r, err := zip.OpenReader(srcFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	err = os_.MakeDirAll(destDir)
	if err != nil {
		return nil, err
	}

	var extractedFiles []*option.FileInfo
	for _, f := range r.File {
		fileInfo, err := z.extractAndWriteFile(destDir, f)
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, fileInfo)
	}

	return extractedFiles, nil
}

/*
func (z *Zip) ExtractStream(srcFile, destDir string) ([]*option.FileInfo, error) {
	r, err := zip.OpenReader(srcFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	err = os_.MakeDirAll(destDir)
	if err != nil {
		return nil, err
	}

	var extractedFiles []*option.FileInfo
	for _, f := range r.File {
		fileInfo, err := z.extractAndWriteFile(destDir, f)
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, fileInfo)
	}

	return extractedFiles, nil
}
*/
func (z *Zip) extractAndWriteFile(
	destDir string,
	f *zip.File,
) (*option.FileInfo, error) {

	if f == nil {
		return nil, fmt.Errorf("invalid zip file")
	}

	decodeName := f.Name
	if !utf8.Valid([]byte(f.Name)) {
		i := bytes.NewReader([]byte(f.Name))
		decoder := transform.NewReader(
			i,
			simplifiedchinese.GB18030.NewDecoder(),
		)
		content, err := ioutil.ReadAll(decoder)
		if err != nil {
			return nil, err
		}
		decodeName = string(content)
	}

	baseName := filepath.Base(f.Name)
	if strings.HasPrefix(baseName, ".") {
		return nil, nil
	}

	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	cleanFunc := func() error {
		if err = rc.Close(); err != nil {
			return err
		}

		return nil
	}
	defer cleanFunc()

	path := filepath.Join(destDir, decodeName)
	if f.FileInfo().IsDir() {
		err = os_.MakeDirAll(path)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	fn, err := os_.OpenFile(path, false)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fn, rc)
	if err != nil {
		return nil, err
	}

	return &option.FileInfo{
		Path:     fn.Name(),
		FileInfo: f.FileInfo(),
	}, nil
}

/*

const (
	UNZIP_TYPE_NO_COPY uint = 1
	UNZIP_TYPE_COPY         = 2
)

type UnCompressMsg struct {
	FilePath string
	FileSize uint64
	Err      error
}

func Unzip(src, dest string, operationType uint) <-chan UnCompressMsg {
	unzipMsgCh := make(chan UnCompressMsg, 10000)

	go func() {
		defer close(unzipMsgCh)
		r, err := zip.OpenReader(src)
		if err != nil {
			unzipMsgCh <- UnCompressMsg{
				Err: err,
			}
			return
		}
		defer r.Close()

		var decodeName string
		for _, f := range r.File {
			if !utf8.Valid([]byte(f.Name)) {
				//`	if f.Flags == 0 {
				//如果标致位是0  则是默认的本地编码   默认为gbk
				i := bytes.NewReader([]byte(f.Name))
				decoder := transform.NewReader(
					i,
					simplifiedchinese.GB18030.NewDecoder(),
				)
				content, _ := ioutil.ReadAll(decoder)
				decodeName = string(content)
			} else {
				//如果标志为是 1 << 11也就是 2048  则是utf-8编码
				decodeName = f.Name
			}
			baseName := filepath.Base(f.Name)
			if strings.HasPrefix(baseName, ".") {
				continue
			}

			rc, err := f.Open()
			if err != nil {
				unzipMsgCh <- UnCompressMsg{
					Err: err,
				}
				return
			}

			path := filepath.Join(dest, decodeName)
			if f.FileInfo().IsDir() {
				//err = os.MkdirAll(path, 0755)
				err = os.MkdirAll(path, f.Mode())
				if err != nil {
					unzipMsgCh <- UnCompressMsg{
						Err: err,
					}
					return
				}
			} else {
				/*
					err = os.MkdirAll(filepath.Dir(path), 0755)
					if err != nil {
						unzipMsgCh <- UnCompressMsg{
							Err: err,
						}
						return
					}
*/
/*

				fn, err := os.OpenFile(
					path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					unzipMsgCh <- UnCompressMsg{
						//	FilePath: f.Name(),
						Err: err,
					}
					rc.Close()
					return
				}

				if operationType == UNZIP_TYPE_COPY {
					_, err = io.Copy(fn, rc)
					if err != nil {
						unzipMsgCh <- UnCompressMsg{
							//		FilePath: f.Name(),
							Err: err,
						}
						fn.Close()
						rc.Close()
						return
					}
				}

				unzipMsgCh <- UnCompressMsg{
					FilePath: fn.Name(),
					FileSize: uint64(f.FileInfo().Size()),
				}
				fn.Close()
				rc.Close()
			}
		}
	}()

	return unzipMsgCh
}
*/
