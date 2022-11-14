package date

import (
	"context"

	context_ "github.com/kaydxh/golang/go/context"
	grpc_ "github.com/kaydxh/golang/go/net/grpc"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	grpc_.Repository[DateServiceClient]
}

func (r *Repository) Now(ctx context.Context, req *NowRequest) (resp *NowResponse, err error) {

	ctx, cancel := context_.WithTimeout(ctx, r.Timeout)
	defer cancel()

	resp, err = r.Client.Now(ctx, req)
	if err != nil {
		logrus.Errorf("failed to call Now method, err: %v", err)
	}

	return resp, nil
}
