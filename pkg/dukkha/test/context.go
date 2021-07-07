package dukkha_test

import (
	"context"

	"arhat.dev/dukkha/pkg/dukkha"
)

func NewTestContext(ctx context.Context) dukkha.Context {
	return dukkha.NewConfigResolvingContext(
		ctx, nil, nil, true, 1,
	)
}
