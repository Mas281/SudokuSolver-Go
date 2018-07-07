package sudoku

func newCell(puzzle *Puzzle, value, row, column int) Cell {
	return Cell{
		puzzle:         puzzle,
		value:          value,
		x:              row,
		y:              column,
		possibleValues: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
}

type Cell struct {
	puzzle *Puzzle

	value int

	x int
	y int

	possibleValues []int
}

func (cell *Cell) trySolve() (bool, error) {
	possibleValues := cell.getPossibleValues()

	if len(possibleValues) == 0 {
		return false, Error{"Zero possible values found for cell"}
	}

	cell.possibleValues = possibleValues

	if len(possibleValues) == 1 {
		cell.value = possibleValues[0]
		return true, nil
	}

	return false, nil
}

func (cell *Cell) trySolveRecursive(unsolved []*Cell, index int, counter *counter) (bool, error) {
	counter.count++

	possibleValues := cell.getPossibleValues()
	length := len(possibleValues)

	if length == 0 {
		return false, nil
	}

	if index == len(unsolved) - 1 {
		// Last cell
		if length == 1 {
			cell.value = possibleValues[0]
			return true, nil
		} else {
			return false, Error{"More than one possible value for last cell"}
		}
	} else {
		for _, value := range possibleValues {
			cell.value = value

			nextIndex := index + 1
			nextCell := unsolved[nextIndex]

			solved, err := nextCell.trySolveRecursive(unsolved, nextIndex, counter)

			if err != nil {
				return false, err
			}

			if solved {
				return true, nil
			}
		}

		cell.value = 0
		return false, nil
	}
}

func (cell *Cell) getPossibleValues() []int {
	return reducePossibleValues(cell.possibleValues, cell.getCellsToCheck())
}

func (cell *Cell) getCellsToCheck() []*Cell {
	puzzle := cell.puzzle

	row := puzzle.getRow(cell)
	column := puzzle.getColumn(cell)
	box := puzzle.getBox(cell)

	checkCells := make([]*Cell, 0)

	checkCells = append(checkCells, row...)
	checkCells = append(checkCells, column...)
	checkCells = append(checkCells, box...)

	checkCells = removeDuplicateValues(checkCells)

	return checkCells
}

func (cell Cell) IsSolved() bool {
	return cell.value != 0
}

func removeDuplicateValues(cells []*Cell) []*Cell {
	previous := make(map[int]bool)
	result := make([]*Cell, 0)

	for _, cell := range cells {
		if !cell.IsSolved() || previous[cell.value] {
			continue
		}

		previous[cell.value] = true
		result = append(result, cell)
	}

	return result
}

func reducePossibleValues(possibleValues []int, cells []*Cell) []int {
	newPossibleValues := make([]int, 0)

	valueLoop: for _, value := range possibleValues {
		for _, cell := range cells {
			if cell.value == value {
				continue valueLoop
			}
		}

		newPossibleValues = append(newPossibleValues, value)
	}

	return newPossibleValues
}