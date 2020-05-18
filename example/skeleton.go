package main

import (
	"github.com/shreve/tui"
)

type StartMode struct {
}

func (m *StartMode) InputHandler(in string) {
	switch in {
	case "q", tui.CtrlC:

		// Gracefully quit
		app.Done()
	}
}

func (m *StartMode) Render(height, width int) tui.View {
	out := make(tui.View, height)

	// Render your app into the lines of out

	return out
}

var start StartMode

func main() {
	app := tui.NewApp()
	app.AddMode(0, &start)
	app.Run()
}
