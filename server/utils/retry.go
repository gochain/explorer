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
		if i >= (attempts - 1) {
			break
		}
		SleepCtx(ctx, sleep)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
