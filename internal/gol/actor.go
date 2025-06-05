package gol

import (
	"context"
)

func nextGeneration(curr GridState, nWorkers int) InitSetterFn {
	sem := make(chan struct{}, nWorkers) // semaphore to limit the number of goroutines
	return func(mX, mY int, matrix [][]uint8) {
		for row := 0; row < len(matrix); row++ {
			for col := 0; col < len(matrix[row]); col++ {
				sem <- struct{}{} // acquire semaphore
				go func() {
					defer func() { <-sem }() // release semaphore
					liveNeighbours := curr.getLiveNeighbours(row, col)
					currState := curr.matrix[row][col]
					switch currState {
					case live:
						if liveNeighbours == 2 || liveNeighbours == 3 {
							matrix[row][col] = live
							return
						} else {
							matrix[row][col] = dead
							return
						}
					case dead:
						if liveNeighbours == 3 {
							matrix[row][col] = live
						}
						return
					default:
						return
					}
				}()
			}
		}
	}
}

// RunGridActor manages the Game of Life states
// ctx context.Context parent context
// model *GridModel grid model
func RunGridActor(ctx context.Context, model GridModel, nextTick <-chan struct{}) {
	go func() {
		state := NewGridState(model, glider())
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-nextTick:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				default:
					state.Print()
					state = NewGridState(model, nextGeneration(*state, model.Concurrency)) // every goroutine uses the same copy of prev state
				}
			}
		}
	}()
}
