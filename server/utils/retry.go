package utils

import (
	"context"
	"fmt"
	"time"
)

func Retry(ctx context.Context, attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if i >= (attempts - 1) {
			break
		}
		if SleepCtx(ctx, sleep) != nil {
			return ctx.Err()
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
