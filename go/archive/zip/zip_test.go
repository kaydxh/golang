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
package zip_test

import (
	archivezip "archive/zip"
	"os"
	"path/filepath"
	"testing"

	zip_ "github.com/kaydxh/golang/go/archive/zip"
)

func TestExtractZip(t *testing.T) {
	root := t.TempDir()
	srcFile := filepath.Join(root, "新词词典.zip")
	file, err := os.Create(srcFile)
	if err != nil {
		t.Fatal(err)
	}
	writer := archivezip.NewWriter(file)
	entry, err := writer.Create("词典/内容.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte("content")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	destDir := filepath.Join(root, "unzip")
	fileInfos, err := zip_.Zip{}.Extract(srcFile, destDir)
	if err != nil {
		t.Fatalf("failed to Extract zip file: [%v], err: [%v]", srcFile, err)
	}
	if len(fileInfos) != 1 {
		t.Fatalf("file infos = %+v", fileInfos)
	}
	raw, err := os.ReadFile(filepath.Join(destDir, "词典", "内容.txt"))
	if err != nil || string(raw) != "content" {
		t.Fatalf("content = %q err=%v", raw, err)
	}
}
