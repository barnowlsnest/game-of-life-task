package main

import (
	"context"
	. "github.com/dshlychkou/game-of-life-task/internal/gol"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	rows = 25 // dimension X
	cols = 25 // dimension Y
)

/**
 * Potential improvements of the task:
 * - I would search for more optimal way how to work with the grid state: right now every tick performs O(n^2)
 * - Maybe it makes sense to handle states in segments not the whole grid as we do a lot of redundant CPU time.
 * - Maybe we can have the hashtable where the key is each cell and values are set of neighbours - it probably will provide +- O(1) access and O(n) traversal
 * - What to do when we reach the edge of the boundary. Right now app will consume resource and do nothing - detect and quite or try to reset state?
 */
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		defer close(done)
		<-sig
		cancel()
	}()

	nextTick := ScheduleTicks(ctx, time.Millisecond*100)
	RunGridActor(ctx, GridModel{
		MaxX:        rows,
		MaxY:        cols,
		Concurrency: rows,
	}, nextTick)

	<-done

	os.Exit(0)
}
