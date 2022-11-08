package s3_test

import (
	"fmt"
	"testing"

	s3_ "github.com/kaydxh/golang/pkg/storage/s3"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"golang.org/x/net/context"
)

func TestS3Upload(t *testing.T) {
	cfgFile := "./s3.yaml"
	config := s3_.NewConfig(s3_.WithViper(viper_.GetViper(cfgFile, "storage.s3")))

	s3Bucket, err := config.Complete().New(context.Background())
	if err != nil {
		t.Fatalf("failed to new config err: %v", err)
	}

	testCases := []struct {
		key  string
		data []byte
	}{
		{
			key:  "keytest",
			data: []byte("123"),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err = s3Bucket.WriteAll(context.Background(), testCase.key, testCase.data, nil)
			if err != nil {
				t.Fatalf("failed to WriteAll, got : %s", err)
			}
		})
	}

}
