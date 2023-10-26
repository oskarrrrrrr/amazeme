package grid

import (
	"math/rand"
	"time"
)

type Grid interface {
	AddWall(c Cell, d Dir)
	RemoveWall(c Cell, d Dir)
	String() string
	Rows() int
	Cols() int
}

func FillBorder(grid Grid) {
	for col := 0; col < grid.Cols(); col++ {
		grid.AddWall(Cell{0, col}, DirUp)
		grid.AddWall(Cell{grid.Rows() - 1, col}, DirDown)
	}
	for row := 0; row < grid.Rows(); row++ {
		grid.AddWall(Cell{row, 0}, DirLeft)
		grid.AddWall(Cell{row, grid.Cols() - 1}, DirRight)
	}
}

func Fill(grid Grid) {
	for row := 0; row < grid.Rows(); row++ {
		for col := 0; col < grid.Cols(); col++ {
			grid.AddWall(Cell{row, col}, DirUp)
			grid.AddWall(Cell{row, col}, DirLeft)
			grid.AddWall(Cell{row, col}, DirRight)
			grid.AddWall(Cell{row, col}, DirDown)
		}
	}
}

func BinaryTreeGen(grid Grid) {
	rSrc := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rSrc)

	Fill(grid)
	for row := 0; row < grid.Rows(); row++ {
		for col := 0; col < grid.Cols(); col++ {
			if col == grid.Cols()-1 && row == grid.Rows()-1 {
				continue
			}
			var dir Dir
			if col == grid.Cols()-1 {
				dir = DirDown
			} else if row == grid.Rows()-1 {
				dir = DirRight
			} else {
				if r.Intn(2) == 0 {
					dir = DirDown
				} else {
					dir = DirRight
				}
			}
			grid.RemoveWall(Cell{row, col}, dir)
		}
	}
}

func SidewinderGen(grid Grid) {
	rSrc := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rSrc)

	Fill(grid)
	for row := 0; row < grid.Rows()-1; row++ {
		currSeg := 0
		for col := 0; col < grid.Cols(); col++ {
			if col != grid.Cols()-1 && r.Intn(3) < 2 {
				grid.RemoveWall(Cell{Row: row, Col: col}, DirRight)
				currSeg++
			} else {
				offset := 0
				if currSeg > 0 {
					offset = r.Intn(currSeg)
				}
				grid.RemoveWall(Cell{Row: row, Col: col - offset}, DirDown)
				currSeg = 0
			}
		}
		if currSeg > 0 {
			grid.RemoveWall(Cell{Row: row, Col: grid.Cols() - 1 - r.Intn(currSeg)}, DirDown)
			currSeg = 0
		}
	}

	for col := 0; col < grid.Cols()-1; col++ {
		grid.RemoveWall(Cell{Row: grid.Rows() - 1, Col: col}, DirRight)
	}
}
