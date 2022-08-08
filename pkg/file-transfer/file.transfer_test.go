package filetransfer_test

import (
	"context"
	"fmt"
	"testing"

	filetransfer_ "github.com/kaydxh/golang/pkg/file-transfer"
	viper_ "github.com/kaydxh/golang/pkg/viper"
)

func TestDownload(t *testing.T) {
	cfgFile := "./ft.yaml"
	config := filetransfer_.NewConfig(filetransfer_.WithViper(viper_.GetViper(cfgFile, "filetransfer")))
	ft, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	testCases := []struct {
		Url string
	}{
		{
			Url: "http://127.0.0.1:3306",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			data, err := ft.Download(context.Background(), testCase.Url)
			if err != nil {
				t.Fatalf("failed to download: %v, got : %s", testCase.Url, err)
			}
			t.Logf("data len: %v", len(data))
		})
	}

}
