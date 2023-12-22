package etcd

import (
	"context"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func Watch(ctx context.Context, kv *clientv3.Client, key string, createCallbackFunc, deleteCallbackFunc EventCallbackFunc) {
	ch := kv.Watch(ctx, key, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
	go func() error {
		for resp := range ch {
			for _, event := range resp.Events {
				switch event.Type {
				case mvccpb.PUT:
					if createCallbackFunc != nil {
						createCallbackFunc(ctx, string(event.Kv.Key), string(event.Kv.Value))
					}
				case mvccpb.DELETE:
					if deleteCallbackFunc != nil {
						deleteCallbackFunc(ctx, string(event.Kv.Key), string(event.Kv.Value))

					}
				}
			}
		}

		return nil
	}()

	return
}
