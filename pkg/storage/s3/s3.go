package s3

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

type StorageConfig struct {
	Region       string
	BucketName   string
	AppId        string
	SecretId     string
	SecretKey    string
	SessionToken string
}

type Storage struct {
	Conf   StorageConfig
	bucket *blob.Bucket
}

func NewStorage(ctx context.Context, conf StorageConfig) (*Storage, error) {
	s := &Storage{
		Conf: conf,
	}

	s3Config := aws.NewConfig().WithRegion(conf.Region).WithCredentials(credentials.NewStaticCredentials(conf.SecretId, conf.SecretKey, conf.SessionToken))

	sess, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}

	bucket, err := s3blob.OpenBucket(ctx, sess, conf.BucketName, &s3blob.Options{
		UseLegacyList: true,
	})
	if err != nil {
		return nil, err
	}
	s.bucket = bucket
	return s, nil
}

func ParseUrl(rawUrl string) (StorageConfig, error) {
	var conf StorageConfig
	s3Url, err := url.Parse(rawUrl)
	if err != nil {
		return conf, err
	}
	ss := strings.Split(s3Url.Host, ".")
	bucketName := ss[0]
	conf.Region = ss[2]

	// bucket-AppId
	idx := strings.LastIndex(bucketName, "-")
	if idx < 0 {
		return conf, fmt.Errorf("missed AppId in url")
	}
	//conf.BucketName = bucketName[:idx]
	conf.BucketName = bucketName
	conf.AppId = bucketName[idx+1:]

	fmt.Printf("conf: %v\n", conf)

	return conf, nil
}
