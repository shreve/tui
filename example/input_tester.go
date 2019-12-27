package main

import (
	"fmt"
	"github.com/shreve/tui"
)

type State struct {
	number int
	count uint
	input string
}

var state State

func inputHandler(input string, app *tui.App) {
	state.input = input
	state.count++
	switch input {
	case "q", tui.CtrlC:
		app.Done()
	case tui.KeyUp, tui.KeyRight:
		state.number++
	case tui.KeyDown, tui.KeyLeft:
		state.number--
	}
	app.Redraw()
}

func view(height, width int) tui.View {
	view := make(tui.View, height)
	view[0] = "Input Value Tester"
	view[1] = fmt.Sprintf("  Number: %d, Count: %d", state.number, state.count)
	view[2] = fmt.Sprintf("  Bytes: %v, String: %#v", []byte(state.input), state.input)
	return view
}

func main() {
	a := tui.NewApp()
	state.count = 0
	a.InputHandler = inputHandler
	a.CurrentView = view
	a.Run()
}
