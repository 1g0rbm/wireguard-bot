package sessionctx

import "context"

type ctxKey string

const (
	ctxUsernameKey ctxKey = "username"
)

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, ctxUsernameKey, username)
}

func ExtractUsername(ctx context.Context) string {
	if username, ok := ctx.Value(ctxUsernameKey).(string); ok {
		return username
	}

	return ""
}
