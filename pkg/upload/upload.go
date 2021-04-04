package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	md5_ "github.com/kaydxh/golang/go/crypto/md5"
	io_ "github.com/kaydxh/golang/go/io"
)

type UploadPartInput struct {
	PartId uint32
	FileId int64
	Offset int64
	Length int64
	Md5Sum string
}

func UploadMultipart(
	file *multipart.FileHeader,
	partInput *UploadPartInput,
	filePath string,
) error {
	if file == nil {
		return fmt.Errorf("invalid multipart FileHeader")
	}
	if partInput == nil {
		return fmt.Errorf("invalid partInput")
	}
	if filePath == "" {
		return fmt.Errorf("invalid filePath")
	}

	//check path
	absFilePath, err := filepath.Abs(filePath)
	if strings.Contains(absFilePath, "..") {
		err = fmt.Errorf("invalid filePath: %v", absFilePath)
		return err
	}
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if partInput.Md5Sum != "" {
		sum, err := md5_.SumReader(srcFile)
		if err != nil {
			return err
		}
		gotSum := strings.ToLower(sum)
		expectSum := strings.ToLower(partInput.Md5Sum)

		if gotSum != expectSum {
			return fmt.Errorf("failed to check md5Sum, got: %v, expect: %v", gotSum, expectSum)
		}

	}

	err = io_.WriteReaderAt(filePath, srcFile, partInput.Offset, partInput.Length)
	if err != nil {
		return err
	}

	return nil
}

func CompleteMultipartUpload(
	filePath string,
	md5Sum string) error {
	if filePath == "" {
		return fmt.Errorf("invalid filePath")
	}

	if md5Sum != "" {
		sum, err := md5_.SumFile(filePath)
		if err != nil {
			return err
		}

		gotSum := strings.ToLower(sum)
		expectSum := strings.ToLower(md5Sum)

		if gotSum != expectSum {
			return fmt.Errorf("failed to check md5Sum, got: %v, expect: %v", gotSum, expectSum)
		}

		return nil
	}

	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file: %v is not existed", filePath)
	}

	return nil
}
