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
package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/kaydxh/golang/go/archive/option"
	os_ "github.com/kaydxh/golang/go/os"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Zip implements the Archiver interface for ZIP archive extraction
type Zip struct {
}

func (z Zip) Extract(srcFile, destDir string) ([]*option.FileInfo, error) {
	r, err := zip.OpenReader(srcFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file %s: %w", srcFile, err)
	}
	defer r.Close()

	err = os_.MakeDirAll(destDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	var extractedFiles []*option.FileInfo
	for _, f := range r.File {
		fileInfo, err := z.extractAndWriteFile(destDir, f)
		if err != nil {
			return nil, fmt.Errorf("failed to extract file %s: %w", f.Name, err)
		}

		if fileInfo != nil {
			extractedFiles = append(extractedFiles, fileInfo)
		}
	}

	return extractedFiles, nil
}

func (z Zip) ExtractStream(
	srcFile, destDir string,
) <-chan option.ExtractMsg {

	fileInfoCh := make(chan option.ExtractMsg, 1024)

	go func() {
		defer close(fileInfoCh)
		r, err := zip.OpenReader(srcFile)
		if err != nil {
			fileInfoCh <- option.ExtractMsg{
				Error: fmt.Errorf("failed to open zip file %s: %w", srcFile, err),
			}
			return
		}
		defer r.Close()

		err = os_.MakeDirAll(destDir)
		if err != nil {
			fileInfoCh <- option.ExtractMsg{
				Error: fmt.Errorf("failed to create destination directory %s: %w", destDir, err),
			}
			return
		}

		for _, f := range r.File {
			fileInfo, err := z.extractAndWriteFile(destDir, f)
			if err != nil {
				fileInfoCh <- option.ExtractMsg{
					Error: fmt.Errorf("failed to extract file %s: %w", f.Name, err),
				}
				// Continue processing other files instead of returning
				continue
			}

			if fileInfo != nil {
				fileInfoCh <- option.ExtractMsg{
					FileInfo: fileInfo,
					Error:    nil,
				}
			}
		}
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
		content, err := io.ReadAll(decoder)
		if err != nil {
			return nil, fmt.Errorf("failed to decode filename %s: %w", f.Name, err)
		}
		decodeName = string(content)
	}

	baseName := filepath.Base(f.Name)
	if strings.HasPrefix(baseName, ".") {
		return nil, nil
	}

	rc, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s from zip: %w", f.Name, err)
	}
	defer func() {
		if closeErr := rc.Close(); closeErr != nil {
			// If there's already an error, wrap it; otherwise set the close error
			if err != nil {
				err = fmt.Errorf("%w; close error: %v", err, closeErr)
			} else {
				err = fmt.Errorf("failed to close file %s: %w", f.Name, closeErr)
			}
		}
	}()

	// Security check: prevent zip slip attack by ensuring the resolved path is within destDir
	path := filepath.Join(destDir, decodeName)
	resolvedPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path %s: %w", path, err)
	}
	destDirAbs, err := filepath.Abs(destDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve destination directory %s: %w", destDir, err)
	}
	if !strings.HasPrefix(resolvedPath, destDirAbs+string(filepath.Separator)) && resolvedPath != destDirAbs {
		return nil, fmt.Errorf("zip slip detected: file path %s is outside destination directory %s", decodeName, destDir)
	}

	if f.FileInfo().IsDir() {
		err = os_.MakeDirAll(resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", resolvedPath, err)
		}

		return nil, nil
	}

	fn, err := os_.OpenFile(resolvedPath, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", resolvedPath, err)
	}
	defer fn.Close()

	_, err = io.Copy(fn, rc)
	if err != nil {
		return nil, fmt.Errorf("failed to write file %s: %w", resolvedPath, err)
	}

	return &option.FileInfo{
		Path:     fn.Name(),
		FileInfo: f.FileInfo(),
	}, nil
}
