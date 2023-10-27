package cli

import (
	"fmt"
	"github.com/oskarrrrrrr/amazeme/internal/grid"
	"os"
	"path"
)

// TODO: add some options
func Cli() {
	N := 10
	g := grid.NewBoolGrid(N, N)
	grid.SidewinderGen(&g)

	fmt.Printf("generating a maze of size %vx%v\n", N, N)

	outPath := "out/maze.svg"
	dir := path.Dir(outPath)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create the output directory: %v", dir)
		os.Exit(1)
	}
	g.ToSVG(outPath)
}
