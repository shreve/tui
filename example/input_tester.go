package main

import (
	"fmt"
	"github.com/shreve/tui"
)

type InputMode struct {
	number int64
	count  uint64
	input  string
}

func (m *InputMode) InputHandler(input string) {

	m.input = input
	m.count++

	switch input {

	case "q", tui.CtrlC:
		app.Done()

	case tui.KeyUp, tui.KeyRight:
		m.number++

	case tui.KeyDown, tui.KeyLeft:
		m.number--
	}
}

func (m *InputMode) Render(height, width int) tui.View {
	view := make(tui.View, height)
	view[0] = "Input Value Tester"
	view[1] = fmt.Sprintf("  Number: %d, Count: %d", m.number, m.count)
	view[2] = fmt.Sprintf("  Bytes: %v, String: %#v", []byte(m.input), m.input)
	return view
}

var app = tui.NewApp()

func main() {
	app.AddMode(0, &InputMode{})
	app.Run()
}
