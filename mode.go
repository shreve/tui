package tui

type Mode struct {
	View Renderable
	InputHandler InputHandler
	Init func()
}

var DefaultInputHandler = func(input string, app *App) {
	switch input {
	case "q", CtrlC:
		app.Done()
	}
}

func defaultView(height, width int) View {
	view := make(View, 0)
	view[0] = "Hello! Thanks for using shreve/tui!"
	view[1] = "To get started, make a new mode to replace this one."
	view[3] = "Press `q` to quit."
	return view
}

var defaultMode = Mode{
	defaultView,
	DefaultInputHandler,
	func() { },
}
