package grid

import (
	"github.com/ajstarks/svgo"
	"os"
	"strings"
)

type BoolGrid struct {
	walls []bool
	// number of rows and columns of cells
	rows, cols int
}

type Cell struct {
	Row, Col int
}

type Dir uint8

const (
	DirUp = iota
	DirDown
	DirLeft
	DirRight
)

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

func (grid BoolGrid) getIdx(c Cell, d Dir) int {
	switch d {
	case DirUp:
		return (c.Row * (grid.cols + 1)) + (c.Row * grid.cols) + c.Col
	case DirDown:
		return ((c.Row + 1) * (grid.cols + 1)) + ((c.Row + 1) * grid.cols) + c.Col
	case DirLeft:
		return (c.Row * (grid.cols + 1)) + ((c.Row + 1) * grid.cols) + c.Col
	case DirRight:
		return (c.Row * (grid.cols + 1)) + ((c.Row + 1) * grid.cols) + c.Col + 1
	default:
		panic("Unexpected direction.")
	}
}

func (grid *BoolGrid) AddWall(c Cell, d Dir) {
	grid.walls[grid.getIdx(c, d)] = true
}

func (grid *BoolGrid) RemoveWall(c Cell, d Dir) {
	grid.walls[grid.getIdx(c, d)] = false
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
		lineStyle := "stroke=\"black\" stroke-width=\"2\""
		canvas.Line(margin+x1, margin+y1, margin+x2, margin+y2, lineStyle)
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
