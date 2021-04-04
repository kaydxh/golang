package upload_test

import (
	"mime/multipart"
	"path/filepath"
	"testing"

	os_ "github.com/kaydxh/golang/go/os"
	upload_ "github.com/kaydxh/golang/pkg/upload"
)

func TestUploadMultipart(t *testing.T) {
	var file *multipart.FileHeader
	workDir, _ := os_.Getwd()
	testFilePath := filepath.Join(workDir, "test-file-upload")
	partInput := &upload_.UploadPartInput{
		PartId: 1,
		FileId: 1,
		Offset: 0,
		Length: 18,
	}
	err := upload_.UploadMultipart(file,
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
