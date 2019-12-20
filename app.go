package tui

import (
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

	InputHandler func([]byte, *App)
	CurrentView Renderable
}

func NewApp() *App {

	// Set up the app with non-zero defaults
	a := App{}
	a.cond = *sync.NewCond(&a.lock)
	a.running = true

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
	newRender := a.CurrentView()
	rows, cols := ansi.WindowSize()

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
		b = b[0:count]

		a.InputHandler(b, a)
	}
}
