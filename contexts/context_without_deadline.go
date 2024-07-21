package contexts

import (
	"context"
	"time"
)

type ContextWithoutDeadline struct {
	ctx context.Context
}

func (*ContextWithoutDeadline) Deadline() (time.Time, bool) { return time.Time{}, false }
func (*ContextWithoutDeadline) Done() <-chan struct{}       { return nil }
func (*ContextWithoutDeadline) Err() error                  { return nil }

func (l *ContextWithoutDeadline) Value(key interface{}) interface{} {
	return l.ctx.Value(key)
}

func Copy(ctx context.Context) *ContextWithoutDeadline {
	return &ContextWithoutDeadline{
		ctx: ctx,
	}
}
