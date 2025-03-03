package processid

import (
	"context"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
)

type key int

const processIDKey key = 0

// NewContext -
func NewContext(ctx context.Context) context.Context {
	uuid, _ := common.UUID()
	return context.WithValue(ctx, processIDKey, uuid)
}

// FromContext -
func FromContext(ctx context.Context) string {
	procID, _ := ctx.Value(processIDKey).(string)
	return procID
}
