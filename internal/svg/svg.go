package svg

import (
	"fmt"
	"io"
	"strings"
)

type SVG struct {
	Writer io.Writer
}

func New(w io.Writer) SVG {
	return SVG{w}
}

func (svg SVG) Start(w, h int) {
	svg.printf("<?xml version=\"1.0\"?>\n")
	svg.printf(`<svg width="%d" height="%d"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">`, w, h)
	svg.printf("\n")
}

func (svg SVG) End() {
	svg.printf("</svg>")
}

func (svg SVG) Line(x1 int, y1 int, x2 int, y2 int, s ...string) {
	svg.printf(
		`<line x1="%d" y1="%d" x2="%d" y2="%d" %s/>`,
		x1, y1, x2, y2, strings.Join(s, " "),
	)
	svg.printf("\n")
}

func (svg SVG) Use(id string, x, y int) {
	svg.printf(`<use href="%s" x="%d" y="%d"/>`, id, x, y)
}

func Attr(name, value string) string {
	return fmt.Sprintf(`%s="%s"`, name, value)
}

func (svg SVG) printf(format string, a ...any) (n int, errno error) {
	return fmt.Fprintf(svg.Writer, format, a...)
}
