package sudoku

import (
	"fmt"
	"strconv"
	"time"
)

func NewPuzzle(values [9][9]int) *Puzzle {
	puzzle := &Puzzle{}
	cells := [9][9]*Cell{}

	for columnCount, row := range values {
		for rowCount, value := range row {
			cell := newCell(puzzle, value, rowCount, columnCount)
			cells[columnCount][rowCount] = &cell
		}
	}

	puzzle.cells = cells

	return puzzle
}

func NewDebugPuzzle(values [9][9]int) *Puzzle {
	puzzle := NewPuzzle(values)
	puzzle.debug = true

	puzzle.Print()

	return puzzle
}

type Puzzle struct {
	cells [9][9]*Cell

	debug bool
}

type Error struct {
	error string
}

func (error Error) Error() string {
	return error.error
}

type counter struct {
	count int
}

func printMsg(puzzle *Puzzle, msg string, args ...interface{}) {
	if puzzle.debug {
		fmt.Printf(msg + "\n", args...)
	}
}

func (puzzle *Puzzle) Solve() error {
	startNanos := time.Now().UnixNano()

	unsolved := puzzle.getUnsolved()
	printMsg(puzzle, "Starting with %v unsolved cells\n", len(unsolved))

	iteration := 0
	solvedCount := 0

	// Keep iterating through unsolved cells until
	// they're either all solvedCount or we go through
	// an iteration where no cells get solvedCount
	for true {
		iteration++
		itSolvedCount := 0

		solvedIndexes := make([]int, 0)

		for i, cell := range unsolved {
			solved, err := cell.trySolve()

			if err != nil {
				return err
			}

			if solved {
				solvedIndexes = append(solvedIndexes, i)

				itSolvedCount++
				printMsg(puzzle, "Solved cell! (%v, %v) = %v [It. %v]", cell.x, cell.y, cell.value, iteration)
			}
		}

		if itSolvedCount > 0 {
			printMsg(puzzle, "Solved %v in iteration %v\n", itSolvedCount, iteration)
		}

		solvedCount += itSolvedCount
		unsolved = getNewUnsolved(unsolved, solvedIndexes)

		if len(unsolved) == 0 || itSolvedCount == 0 {
			break
		}
	}

	printMsg(puzzle, "Solved %v cells", solvedCount)

	if len(unsolved) == 0 {
		printMsg(puzzle, "Puzzle solved! Took %v iterations", iteration)
	} else {
		printMsg(puzzle, "Stage 1 ended with %v cells left [It. %v]", len(unsolved), iteration)
		printMsg(puzzle, "\nAttempting recursive solution...")

		counter := &counter{}

		if _, err := unsolved[0].trySolveRecursive(unsolved, 0, counter); err == nil {
			printMsg(puzzle, "Solved recursively!")
		} else {
			return err
		}

		printMsg(puzzle, "Recursions: %v", counter.count)
	}

	nanosTaken := time.Now().UnixNano() - startNanos
	printMsg(puzzle, "Finished in %vms (%v nanoseconds)", nanosTaken / int64(1e6), nanosTaken)

	return nil
}

func (puzzle *Puzzle) getUnsolved() []*Cell {
	unsolved := make([]*Cell, 0)

	for _, row := range puzzle.cells {
		for _, cell := range row {
			if !cell.IsSolved() {
				unsolved = append(unsolved, cell)
			}
		}
	}

	return unsolved
}

func getNewUnsolved(cells []*Cell, solved []int) []*Cell {
	newUnsolved := make([]*Cell, 0)

	cellLoop: for i, cell := range cells {
		for j, index := range solved {
			if i == index {
				solved = append(solved[:j], solved[j + 1:]...)
				continue cellLoop
			}
		}

		newUnsolved = append(newUnsolved, cell)
	}

	return newUnsolved
}

func (puzzle *Puzzle) getRow(cell *Cell) []*Cell {
	rowCells := make([]*Cell, 0)

	for _, rowCell := range puzzle.cells[cell.y] {
		if rowCell == cell {
			continue
		}

		rowCells = append(rowCells, rowCell)
	}

	return rowCells
}

func (puzzle *Puzzle) getColumn(cell *Cell) []*Cell {
	columnCells := make([]*Cell, 0)

	for _, column := range puzzle.cells {
		columnCell := column[cell.x]

		if columnCell == cell {
			continue
		}

		columnCells = append(columnCells, columnCell)
	}

	return columnCells
}

func (puzzle *Puzzle) getBox(cell *Cell) []*Cell {
	boxCells := make([]*Cell, 0)

	boxX := (cell.x / 3) * 3
	boxY := (cell.y / 3) * 3

	for y := boxY; y <= boxY + 2; y++ {
		for x := boxX; x <= boxX + 2; x++ {
			boxCell := puzzle.cells[y][x]

			if boxCell == cell {
				continue
			}

			boxCells = append(boxCells, boxCell)
		}

	}

	return boxCells
}

func (puzzle Puzzle) Print() {
	gridString := ""

	for _, row := range puzzle.cells {
		rowString := ""

		for _, cell := range row {
			rowString += strconv.Itoa(cell.value) + " "
		}

		rowString = rowString[:len(rowString) - 1] + "\n"
		gridString += rowString
	}

	fmt.Println(gridString)
}