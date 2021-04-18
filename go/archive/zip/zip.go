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

func (z Zip) Extract(srcFile, destDir string) ([]*option.FileInfo, error) {
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

func (z Zip) ExtractStream(
	srcFile, destDir string,
) <-chan option.ExtractMsg {

	fileInfoCh := make(chan option.ExtractMsg, 1024)

	go func() error {
		defer close(fileInfoCh)
		r, err := zip.OpenReader(srcFile)
		if err != nil {
			return err
		}
		defer r.Close()

		err = os_.MakeDirAll(destDir)
		if err != nil {
			return err
		}

		for _, f := range r.File {
			fileInfo, err := z.extractAndWriteFile(destDir, f)
			if err != nil {
				fileInfoCh <- option.ExtractMsg{
					Error: err,
				}
				return err
			}

			if fileInfo != nil {
				fileInfoCh <- option.ExtractMsg{
					FileInfo: fileInfo,
					Error:    err,
				}
				return nil
			}

		}
		return nil
	}()

	return fileInfoCh
}

func (z Zip) extractAndWriteFile(
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
	defer fn.Close()

	_, err = io.Copy(fn, rc)
	if err != nil {
		return nil, err
	}

	return &option.FileInfo{
		Path:     fn.Name(),
		FileInfo: f.FileInfo(),
	}, nil
}
