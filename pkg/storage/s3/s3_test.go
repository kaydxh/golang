package s3_test

import (
	"fmt"
	"testing"
	"time"

	s3_ "github.com/kaydxh/golang/pkg/storage/s3"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"golang.org/x/net/context"
)

func TestS3Upload(t *testing.T) {
	cfgFile := "./s3.yaml"
	config := s3_.NewConfig(s3_.WithViper(viper_.GetViper(cfgFile, "storage.s3")))

	s3Bucket, err := config.Complete().New(context.Background())
	if err != nil || s3Bucket == nil {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			err = s3Bucket.WriteAll(ctx, testCase.key, testCase.data, nil)
			if err != nil {
				t.Fatalf("failed to WriteAll, got : %s", err)
			}
		})
	}

}

func TestS3Down(t *testing.T) {
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
			key: "test_data/dongfangmingzhu.jpeg",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			data, err := s3Bucket.ReadAll(ctx, testCase.key)
			if err != nil {
				t.Fatalf("failed to WriteAll, got : %s", err)
			}
			t.Logf("data len: %v", len(data))
		})
	}

}
