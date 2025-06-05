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
		Concurrency: cols * rows, // in this case, the number of concurrent goroutines are up to Number of columns
	}, nextTick)

	<-done

	os.Exit(0)
}
