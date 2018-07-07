package sudoku

import (
	"testing"
	"fmt"
)

func TestPuzzle_Solve(t *testing.T) {
	puzzle := NewDebugPuzzle([9][9]int{
		{0, 0, 0, 6, 0, 0, 0, 1, 0},
		{0, 6, 9, 1, 0, 3, 0, 0, 0},
		{4, 0, 0, 0, 0, 0, 0, 0, 0},
		{5, 0, 0, 3, 0, 0, 0, 0, 0},
		{7, 0, 2, 0, 0, 0, 6, 0, 9},
		{0, 0, 0, 0, 8, 9, 0, 0, 4},
		{2, 3, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 9, 2, 0, 0, 0, 0},
		{0, 0, 5, 4, 0, 0, 0, 0, 8},
	})

	err := puzzle.Solve()

	if err == nil {
		fmt.Println()
		puzzle.Print()
	} else {
		panic(err)
	}
}