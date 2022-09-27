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
package upload_test

import (
	"mime/multipart"
	"path/filepath"
	"testing"

	os_ "github.com/kaydxh/golang/go/os"
	upload_ "github.com/kaydxh/golang/pkg/upload"
)

func TestUploadMultipart(t *testing.T) {
	var file multipart.FileHeader
	srcFile, err := file.Open()
	if err != nil {
		return
	}
	defer srcFile.Close()

	workDir, _ := os_.Getwd()
	testFilePath := filepath.Join(workDir, "test-file-upload")
	partInput := &upload_.UploadPartInput{
		PartId: 1,
		Offset: 0,
		Length: 18,
	}
	err = upload_.UploadMultipart(srcFile,
		partInput,
		testFilePath)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
		return
	}
}

func TestCompleteMultipartUpload(t *testing.T) {
	workDir, _ := os_.Getwd()
	testFilePath := filepath.Join(workDir, "test-file-upload")
	err := upload_.CompleteMultipartUpload(
		testFilePath,
		"",
	)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
		return
	}
}
