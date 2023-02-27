package date

import (
	"context"

	grpc_ "github.com/kaydxh/golang/go/net/grpc"
)

type Repository struct {
	grpc_.Repository[DateServiceClient]
}

func (r *Repository) Now(ctx context.Context, req *NowRequest) (resp *NowResponse, err error) {

	err = r.Call(ctx, func(ctx context.Context) error {
		nowResp, err := r.Client.Now(ctx, &NowRequest{})
		if err != nil {
			return err
		}

		resp = &NowResponse{
			RequestId: req.RequestId,
			Date:      nowResp.GetDate(),
		}
		return nil
	})

	return resp, err
}
