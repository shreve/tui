package tui

type Inputable interface {
	InputHandler(string)
}

type Mode interface {
	Renderable
	Inputable
}

var noop = func() {}

type DefaultMode struct {
	app *App
}

func (d *DefaultMode) Render(height, width int) View {
	view := make(View, 0)
	view[0] = "Hello! Thanks for using shreve/tui!"
	view[1] = "To get started, make a new mode to replace this one."
	view[3] = "Press `q` to quit."
	return view
}

func (d *DefaultMode) InputHandler(in string) {
	switch in {
	case "q", CtrlC:
		d.app.Done()
	}
}

func (d *DefaultMode) Init() { }
