package function

import "context"

type ConsumerFunc func(ctx context.Context, queue string) error
