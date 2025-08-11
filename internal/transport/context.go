package transport

import "context"

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

func withIsAuthenticated(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, isAuthenticatedContextKey, v)
}

func AuthenticatedFromContext(ctx context.Context) bool {
	v, _ := ctx.Value(isAuthenticatedContextKey).(bool)
	return v
}
