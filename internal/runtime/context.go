package jsruntime

import (
	"context"
	"syscall/js"
)

type (
	contextKeyTriggerObj struct{}
)

func New(ctx context.Context, triggerObj js.Value) context.Context {
	ctx = context.WithValue(ctx, contextKeyTriggerObj{}, triggerObj)
	return ctx
}

// MustExtractTriggerObj extracts trigger object from context.
// This function panics when trigger object was not found.
func MustExtractTriggerObj(ctx context.Context) js.Value {
	v, ok := ctx.Value(contextKeyTriggerObj{}).(js.Value)
	if !ok {
		panic("trigger object was not found")
	}
	return v
}
