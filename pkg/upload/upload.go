package upload

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	io_ "github.com/kaydxh/golang/go/io"
)

type UploadPartInput struct {
	PartId uint32
	FileId int64
	Offset int64
	Length int64
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

	err = io_.WriteReaderAt(filePath, srcFile, partInput.Offset, partInput.Length)
	if err != nil {
		return err
	}

	return nil
}
