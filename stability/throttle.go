package stability

import (
	"cloud_patters/stability/src"
	"context"
	"fmt"
	"sync"
	"time"
)

type effector = src.Circuit

func Throttle(e effector, max uint, refill uint, d time.Duration) effector {
	var tokens = max
	var once sync.Once

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return

					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
					}
				}
			}()
		})

		if tokens <= 0 {
			return "", fmt.Errorf("too many calls")
		}
		tokens--

		return e(ctx)
	}
}
