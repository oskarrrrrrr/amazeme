package grid

import (
	"github.com/oskarrrrrrr/amazeme/internal/svg"
	"os"
	"strings"
)

type BoolGrid struct {
	walls []bool
	// number of rows and columns of cells
	rows, cols int
}

func NewBoolGrid(rows, cols int) BoolGrid {
	horizWalls := (rows + 1) * cols
	vertWalls := rows * (cols + 1)
	walls := vertWalls + horizWalls
	board := BoolGrid{
		walls: make([]bool, walls),
		rows:  rows,
		cols:  cols,
	}
	return board
}

func (grid BoolGrid) Rows() int {
	return grid.rows
}

func (grid BoolGrid) Cols() int {
	return grid.cols
}

func (grid BoolGrid) getIdx(row, col int, d Dir) int {
	switch d {
	case DirUp:
		return (row * (grid.cols + 1)) + (row * grid.cols) + col
	case DirDown:
		return ((row + 1) * (grid.cols + 1)) + ((row + 1) * grid.cols) + col
	case DirLeft:
		return (row * (grid.cols + 1)) + ((row + 1) * grid.cols) + col
	case DirRight:
		return (row * (grid.cols + 1)) + ((row + 1) * grid.cols) + col + 1
	default:
		panic("Unexpected direction.")
	}
}

func (grid *BoolGrid) AddWall(row, col int, d Dir) {
	grid.walls[grid.getIdx(row, col, d)] = true
}

func (grid *BoolGrid) RemoveWall(row, col int, d Dir) {
	grid.walls[grid.getIdx(row, col, d)] = false
}

func (grid BoolGrid) String() string {
	sb := strings.Builder{}
	idx := 0
	for row := 0; row < 2*grid.Rows()+1; row++ {
		if row%2 == 0 {
			sb.WriteByte(' ')
			for col := 0; col < grid.cols; col++ {
				if grid.walls[idx] {
					sb.WriteByte('-')
					sb.WriteByte('-')
				} else {
					sb.WriteByte(' ')
					sb.WriteByte(' ')
				}
				sb.WriteByte(' ')
				idx++
			}
			sb.WriteByte('\n')
		} else {
			for col := 0; col < grid.cols+1; col++ {
				if grid.walls[idx] {
					sb.WriteByte('|')
					sb.WriteByte(' ')
				} else {
					sb.WriteByte(' ')
					sb.WriteByte(' ')
				}
				sb.WriteByte(' ')
				idx++
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func (grid BoolGrid) ToSVG(file_name string) {
	file, err := os.Create(file_name)
	if err != nil {
		panic(err)
	}
	canvas := svg.New(file)
	margin := 10
	cellHeight, cellWidth := 50, 50
	width := (2 * margin) + (grid.Rows() * cellHeight)
	height := (2 * margin) + (grid.Cols() * cellWidth)
	canvas.Start(width, height)

	drawLine := func(x1, y1, x2, y2 int) {
		canvas.Line(
			margin+x1, margin+y1, margin+x2, margin+y2,
			svg.Attr("stroke", "black"),
			svg.Attr("stroke-width", "2"),
		)
	}

	idx := 0
	length := 0
	for row := 0; row <= grid.Rows(); row++ {
		for col := 0; col < grid.Cols(); col++ {
			if grid.walls[idx] {
				length++
			} else if length > 0 {
				drawLine(
					(col-length)*cellWidth,
					row*cellHeight,
					col*cellWidth,
					row*cellHeight,
				)
				length = 0
			}
			idx++
		}
		if length > 0 {
			drawLine(
				(grid.Cols()-length)*cellWidth,
				row*cellHeight,
				grid.Cols()*cellWidth,
				row*cellHeight,
			)
			length = 0
		}
		idx += grid.Cols() + 1
	}

	idx = grid.Cols()
	heights := make([]int, grid.Cols()+1)
	for row := 0; row < grid.Rows(); row++ {
		for col := 0; col <= grid.Cols(); col++ {
			if grid.walls[idx] {
				heights[col]++
			} else if heights[col] > 0 {
				drawLine(
					col*cellWidth,
					(row-heights[col])*cellHeight,
					col*cellWidth,
					row*cellHeight,
				)
				heights[col] = 0
			}
			idx++
		}
		idx += grid.Rows()
	}
	for col, h := range heights {
		if h > 0 {
			drawLine(
				col*cellWidth,
				(grid.Rows()-heights[col])*cellHeight,
				col*cellWidth,
				grid.Rows()*cellHeight,
			)
		}
	}
	canvas.End()
}
