package action

import (
	"context"
)

type key int

const actionKey key = 0

// NewContext -
func NewContext(ctx context.Context, action string) context.Context {
	return context.WithValue(ctx, actionKey, action)
}

// FromContext -
func FromContext(ctx context.Context) string {
	action, _ := ctx.Value(actionKey).(string)
	return action
}
