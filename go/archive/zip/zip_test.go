package zip_test

import (
	"testing"

	zip_ "github.com/kaydxh/golang/go/archive/zip"
)

func TestExtractZip(t *testing.T) {
	srcFile := "./新词词典.zip"
	destDir := "unzip"
	fileInfos, err := zip_.Zip{}.Extract(srcFile, destDir)
	if err != nil {
		t.Errorf("failed to Extract zip file: [%v], err: [%v]", srcFile, err)
	}

	t.Logf("extract file: [%v], result: [%+v]", srcFile, fileInfos)
	/*

		for extractMsg := range zip_.Zip{}.Extract(srcFile, destDir) {
			if extractMsg.Err != nil {
				fmt.Println(msg.Error)
			} else {
			  fmt.Printf("fileInfo:[%v]", *msg.FileInfo)
			}
		}
	*/
}
