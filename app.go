package tui

import (
	"fmt"
	"github.com/pkg/term"
	"github.com/shreve/tui/ansi"
	"os"
	"sync"
	"time"
)

type winSize struct {
	rows int
	cols int
}

type App struct {
	lock           sync.Mutex
	cond           sync.Cond
	running        bool
	term           *term.Term
	lastRender     View
	lastSize       winSize
	watchForResize bool
	modes          map[int]Mode
	mode           Mode

	Cursor   Cursor
	OnResize func(int, int)
}

var defaultOnResize = func(int, int) {}

func NewApp() *App {

	// Set up the app with non-zero defaults
	a := App{}
	a.cond = *sync.NewCond(&a.lock)
	a.running = true
	a.watchForResize = true
	a.modes = make(map[int]Mode)
	a.mode = &DefaultMode{&a}
	a.OnResize = defaultOnResize

	// Use term handle of stdin to set mode and read in bytes
	var err error
	a.term, err = term.Open("/dev/stdin")
	if err != nil {
		panic(err)
	}

	return &a
}

func (a *App) AddMode(id int, mode Mode) {
	a.modes[id] = mode
	if id == 0 {
		a.SetMode(0)
	}
}

func (a *App) SetMode(id int) {
	var ok bool
	a.mode, ok = a.modes[id]
	if ! ok {
		a.Panic("Set mode to a mode that doesn't exist.")
	}
	a.Redraw()
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

	if a.watchForResize {
		go a.resizeWatcher()
	}

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
	size := winSize{rows, cols}
	newRender := a.mode.Render(rows, cols)

	if size != a.lastSize {

		// If the window is a different size, re-draw everything
		a.lastSize = size
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

		a.mode.InputHandler(input)
		a.Redraw()
	}
}

func (a *App) resizeWatcher() {
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for {
		<-tick.C
		rows, cols := ansi.WindowSize()
		size := winSize{rows, cols}
		if size != a.lastSize {
			a.OnResize(rows, cols)
			a.render()
		}
	}
}
