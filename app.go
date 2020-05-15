package tui

import (
	"os"
	"fmt"
	"sync"
	"github.com/pkg/term"
	"github.com/shreve/tui/ansi"
)

type App struct {
	lock sync.Mutex
	cond sync.Cond
	running bool
	term *term.Term
	lastRender View
	lastSize int

	InputHandler func(string, *App)
	CurrentView Renderable
}

const (
	KeyEsc = "\x1b"
	KeyUp = "\x1b[A"
	KeyDown = "\x1b[B"
	KeyLeft = "\x1b[D"
	KeyRight = "\x1b[C"
	KeyDelete = "\x1b[3~"
	KeyBackspace = "\u007f"
	CtrlA = "\x01"
	CtrlB = "\x02"
	CtrlC = "\x03"
	Enter = "\r"
)

var emptyInputHandler = func(input string, app *App) {
	switch input {
	case "q", CtrlC:
		app.Done()
	}
}
var emptyView = func(int, int) View { return make(View, 0) }

func NewApp() *App {

	// Set up the app with non-zero defaults
	a := App{}
	a.cond = *sync.NewCond(&a.lock)
	a.running = true
	a.InputHandler = emptyInputHandler
	a.CurrentView = emptyView

	// Use term handle of stdin to set mode and read in bytes
	var err error
	a.term, err = term.Open("/dev/stdin")
	if err != nil { panic(err) }

	return &a
}

// Finish execution by closing render and input loops
func (a *App) Done() {
	a.running = false
}

func (a *App) Panic(msg string) {
	ansi.RestoreState()
	ansi.ShowCursor()
	a.term.Restore()
	fmt.Println(msg)
	os.Exit(1)
}

// Signal renderer
func (a *App) Redraw() {
	a.cond.Signal()
}

// Set up the app and run the loops
func (a *App) Run() {

	// Save the previous term state and restore it on close
	ansi.SaveState()
	defer ansi.RestoreState()

	// Hide the terminal cursor and restore it on close
	ansi.HideCursor()
	defer ansi.ShowCursor()

	// Set the terminal into raw mode and restore on close
	a.term.SetRaw()
	a.term.SetCbreak()
	defer a.term.Restore()

	go a.renderLoop()
	a.inputLoop()
}

// Wrap rendering in a condition variable so we can signal at will
func (a *App) renderLoop() {
	a.lock.Lock()
	for a.running {
		a.render()
		a.cond.Wait()
	}
	a.lock.Unlock()
}

// Perform the render
func (a *App) render() {
	rows, cols := ansi.WindowSize()
	newRender := a.CurrentView(rows, cols)

	if rows * cols != a.lastSize {

		// If the window is a different size, re-draw everything
		a.lastSize = rows * cols
		newRender.Render()
	} else {

		// Otherwise, do a diff render based on the last draw
		newRender.RenderFrom(a.lastRender)
	}
	a.lastRender = newRender
}

// Read in inputs one key at a time and pass off to user handler
func (a *App) inputLoop() {
	for a.running {
		b := make([]byte, 15)
		count, err := a.term.Read(b)
		if err != nil {
			continue
		}
		input := string(b[0:count])

		a.InputHandler(input, a)
	}
}
