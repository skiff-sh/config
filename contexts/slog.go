package contexts

import (
	"context"
	"log/slog"
)

type slogKey struct{}

func WithLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, slogKey{}, log)
}

func GetLogger(ctx context.Context) *slog.Logger {
	v, _ := ctx.Value(slogKey{}).(*slog.Logger)
	if v != nil {
		return slog.Default()
	}
	return v
}
