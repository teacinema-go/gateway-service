package context

import (
	"context"
)

type ctxKey[T any] struct{}

func WithValue[T any](ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, ctxKey[T]{}, value)
}

func GetValue[T any](ctx context.Context) (T, bool) {
	val, ok := ctx.Value(ctxKey[T]{}).(T)
	return val, ok
}
