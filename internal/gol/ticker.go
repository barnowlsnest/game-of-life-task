package gol

import (
	"context"
	"time"
)

// ScheduleTicks produces new ticks
// duration time.Duration configures the rate of ticker
// return unbuffered channel to sync consumer goroutines with ticker producer
func ScheduleTicks(ctx context.Context, duration time.Duration) <-chan struct{} {
	nextTick := make(chan struct{})
	go func() {
		defer close(nextTick)
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case <-ctx.Done():
				return
			case nextTick <- struct{}{}:
			}
		}
	}()
	return nextTick
}
