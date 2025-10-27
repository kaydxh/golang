package s3

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	http_ "github.com/kaydxh/golang/go/net/http"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

// url: http://examplebucket-1250000000.cos.ap-beijing.myqcloud.com
type StorageConfig struct {
	// 机房分布区域, cos的数据存放在这些地域的存储桶中
	Region string
	// 存储桶名字,存储桶是对象的载体,可理解为存放对象的容器
	BucketName string
	// 对象被存放到存储桶中,用户可通过访问域名访问和下载对象
	Endpoint string
	// 与账户对应的ID
	AppId string
	// 密钥
	SecretId     string
	SecretKey    string
	SessionToken string

	DisableSSL bool
}

type Storage struct {
	Conf StorageConfig
	*blob.Bucket
	opts struct {

		// 前缀路径,一般以"/"结尾, 操作子目录
		// the PrefixPath should end with "/", so that the resulting operates in a subfoleder
		PrefixPath string
	}
}

func NewStorage(ctx context.Context, conf StorageConfig, opts ...StorageOption) (*Storage, error) {
	s := &Storage{
		Conf: conf,
	}
	s.ApplyOptions(opts...)

	client, _ := http_.NewClient()
	s3Config := aws.NewConfig().
		WithRegion(conf.Region).
		WithCredentials(credentials.NewStaticCredentials(conf.SecretId, conf.SecretKey, conf.SessionToken)).
		WithDisableSSL(conf.DisableSSL).
		WithHTTPClient(&client.Client).
		WithEndpoint(conf.Endpoint)

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

	if s.opts.PrefixPath != "" {
		bucket = blob.PrefixedBucket(bucket, s.opts.PrefixPath)
	}

	s.Bucket = bucket
	return s, nil
}

// url: http://examplebucket-1250000000.cos.ap-beijing.myqcloud.com
func ParseUrl(rawUrl string) (*StorageConfig, error) {
	var conf StorageConfig
	s3Url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	if s3Url.Scheme != "https" {
		conf.DisableSSL = true
	}
	ss := strings.Split(s3Url.Host, ".")
	if len(ss) < 3 {
		return nil, fmt.Errorf("the number of dot in %v too less", s3Url.Host)
	}
	bucketName := ss[0]
	conf.Region = ss[2]

	// len(ss) >= 3 so strings.Index(s3Url.Host, ".") must >= 0
	conf.Endpoint = s3Url.Host[strings.Index(s3Url.Host, ".")+1:]

	// bucket-AppId
	idx := strings.LastIndex(bucketName, "-")
	if idx < 0 {
		return nil, fmt.Errorf("missed - in url")
	}
	conf.BucketName = bucketName
	conf.AppId = bucketName[idx+1:]

	return &conf, nil
}
