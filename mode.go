package tui

import "github.com/shreve/tui/ansi"

type Mode struct {
	View Renderable
	InputHandler InputHandler
	Init func()
	Cursor Cursor
}

var noop = func() { }

func NewMode(v Renderable, i InputHandler) *Mode {
	m := Mode{}
	m.View = v
	m.InputHandler = i
	m.Init = noop
	m.Cursor = NewCursor(ansi.WindowSize())
	return &m
}

var DefaultInputHandler = func(input string, app *App, cursor *Cursor) {
	switch input {
	case "q", CtrlC:
		app.Done()
	}
}

func defaultView(height, width int, cursor *Cursor) View {
	view := make(View, 0)
	view[0] = "Hello! Thanks for using shreve/tui!"
	view[1] = "To get started, make a new mode to replace this one."
	view[3] = "Press `q` to quit."
	return view
}

var defaultMode = NewMode(defaultView, DefaultInputHandler)
