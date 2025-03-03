package auth

import "context"

type key int

const authKey key = 0

// NewContext -
func NewContext(ctx context.Context, auth string) context.Context {
	return context.WithValue(ctx, authKey, auth)
}

// FromContext -
func FromContext(ctx context.Context) string {
	auth, _ := ctx.Value(authKey).(string)
	return auth
}
