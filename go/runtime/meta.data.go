package runtime

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func GetMetadata(ctx context.Context, key string) []string {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok || md.HeaderMD == nil {
		return nil
	}

	return md.HeaderMD.Get(key)
}
