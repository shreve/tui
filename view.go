package tui

import (
	"fmt"
	"github.com/shreve/tui/ansi"
)

type View []string
type Renderable interface {
	Render(int, int) View
}

// Draw all the lines in this view.
func (v View) Render() {
	for i := 0; i < len(v); i++ {
		v.drawLine(i)
	}
}

// Only draw lines that differ from a provided view. This is an important
// trade-off. This massively speeds up rendering in most cases, but may cause
// errors when terminal output has changed without our knowledge.
func (v View) RenderFrom(o View) {
	for i := 0; i < len(v); i++ {
		if i >= len(o) || i >= len(v) || v[i] != o[i] {
			v.drawLine(i)
		}
	}
}

// Draw the ith line of the view to the ith line on the screen
func (v View) drawLine(i int) {
	ansi.MoveCursor(i, 0)
	ansi.ClearLine()
	fmt.Print(v[i])
}
