package date_test

import (
	"testing"
	"time"

	"context"

	grpc_ "github.com/kaydxh/golang/go/net/grpc"
	date_ "github.com/kaydxh/golang/go/net/grpc/example"
	"google.golang.org/grpc"
)

func TestNow(t *testing.T) {

	ctx := context.Background()
	/*
		repository, err := grpc_.NewRepository(ctx, grpc_.FactoryConfig[date_.DateServiceClient]{
			Addr:    "localhost:10001",
			Timeout: 5 * time.Second,
			NewServiceClient: func(c *grpc.ClientConn) date_.DateServiceClient {
				return date_.NewDateServiceClient(c)
			},
		},
		)
	*/

	factory, err := grpc_.NewFactory(grpc_.FactoryConfig[date_.DateServiceClient]{
		Addr:    "localhost:10001",
		Timeout: 5 * time.Second,
		NewServiceClient: func(c *grpc.ClientConn) date_.DateServiceClient {
			return date_.NewDateServiceClient(c)
		},
	},
	)
	if err != nil {
		t.Errorf("failed to new factory, err: %v", err)
	}
	repository, err := factory.NewClient(ctx)
	if err != nil {
		t.Errorf("failed to new respository client, err: %v", err)
	}

	respWrap := date_.Repository{
		Repository: repository,
	}
	resp, err := respWrap.Now(ctx, &date_.NowRequest{})
	/*
		var resp any
		err = respWrap.Call(ctx,
			func(ctx context.Context) error {
				// short connection
				newClient, conn, err := respWrap.NewConnect()
				if err != nil {
					return err
				}
				defer conn.Close()

				resp, err = newClient.Now(ctx, &date_.NowRequest{})
				//long connection
				//resp, err = respWrap.Now(ctx, &date_.NowRequest{})
				return err
			})
	*/
	if err != nil {
		t.Errorf("failed to call Now, err: %v", err)
		return
	}

	t.Logf("resp: %v", resp.RequestId)

}
