package gol

import (
	"fmt"
)

const (
	// dead cell
	dead uint8 = iota
	// live cell
	live
)

type GridModel struct {
	MaxX        int
	MaxY        int
	Concurrency int
}

// GridState represents the Game Of Life state
type GridState struct {
	model  GridModel
	matrix [][]uint8
}

// getLiveNeighbours finds adjacent cells with live value
func (gs *GridState) getLiveNeighbours(row, col int) int {
	neighbours := 0
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			nX, nY := row+x, col+y // offsets to neighbors cells
			if x == 0 && y == 0 {  // center of 3x3 boundary
				continue
			}
			if nX < 0 || nX >= gs.model.MaxX || nY < 0 || nY >= gs.model.MaxY { // out of bound
				continue
			}
			if gs.matrix[nX][nY] == live {
				neighbours++
			}
		}
	}
	return neighbours
}

func (gs *GridState) Print() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	for _, row := range gs.matrix {
		for _, col := range row {
			if col == live {
				fmt.Print(" @ ")
			} else {
				fmt.Print(" _ ")
			}
		}
		fmt.Print("\r\n")
	}
}

type InitSetterFn = func(mX, mY int, matrix [][]uint8)

// NewGridState create new Grid's state
// model GridModel
// initFn InitSetterFn
func NewGridState(model GridModel, initFn InitSetterFn) *GridState {
	nRows := model.MaxY
	nCols := model.MaxX
	matrix := make([][]uint8, nRows)
	for row := 0; row < nRows; row++ {
		matrix[row] = make([]uint8, nCols)
	}
	initFn(model.MaxX, model.MaxY, matrix)
	return &GridState{
		model, matrix,
	}
}

func glider() InitSetterFn {
	gliderPattern := [][2]int{
		{-0, -1}, {0, 1}, {1, 0}, {-1, 1}, {1, 1},
	}
	return func(mX, mY int, matrix [][]uint8) {
		centerRow, centerCol := mY/2, mX/2
		for _, offset := range gliderPattern {
			row, col := centerRow+offset[1], centerCol+offset[0]
			if row >= 0 && row < mY && col >= 0 && col < mX {
				matrix[row][col] = 1
			}
		}
	}
}
