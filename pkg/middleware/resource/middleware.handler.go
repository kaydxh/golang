package resource

import "context"

type HandlerWithContext[REQ any, RESP any] func(context.Context, REQ) (RESP, error)
