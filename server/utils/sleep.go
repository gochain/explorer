package utils

import (
	"context"
	"time"
)

// SleepCtx is like time.Sleep(dur), but returns an error if the ctx is Done before then.
func SleepCtx(ctx context.Context, dur time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(dur):
		return nil
	}
}
