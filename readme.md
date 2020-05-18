tui
===

A library for generating simple terminal user interface applications.

This tool primarily handles:
 * Reading input
 * Fast rendering of views
 * Configuring the terminal for easy ANSI code interaction

## Usage

```go
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
```

## Features

Beyond the basic io and rendering, this library includes several features to
help build out terminal based applications.

### Cursor

Provide dimensions of the space, then call Up, Down, Left, Right to move
around within the space.

```go
// Navigate a list with height of 3
cursor := tui.NewCursor(3, 1)
cursor.Size() // (3, 1)
cursor.Position() // (0, 0)
cursor.Up() // false
cursor.Down() // true
cursor.Down() // true
cursor.Down() // false
cursor.Up() // true
cursor.Left() // false
cursor.Right() // false
cursor.Position() // (1, 0)
cursor.Top()
cursor.Bottom()
cursor.Position() // (2, 0)
```

### Table

Table accepts slices and structs and draws their contents and makes them
searchable.

```go
import (
	"github.com/shreve/tui"
	"github.com/shreve/tui/ansi"
)

type Data struct {
    Field1 string
    Field2 int
}

var datum []Data

// Fill datum with elements

table := tui.Table{}
table.Height, table.Width := ansi.WindowSize()
table.Update(datum, []string{"Field1", "Field2"})
table.Draw()

// Use the Search method and the next Draw will be limited to rows
// case-insensitive matching the query.
table.Search("query")
```

### Input Helpers

Functions which cover typical input interactions. Currently, only movement is
covered for now, but any common use cases will be added.

```go
func (m *Mode) InputHandler(in string) {
    switch {

    // Use vi keys to move cursor
    case tui.InputMoveCursor(tui.ViCursor, in, &cursor):

    // Use arrow keys to move cursor
    case tui.InputMoveCursor(tui.ArrowCursor, in, &cursor):

    // Use wasd keys to move cursor
    case tui.InputMoveCursor(tui.WasdCursor, in, &cursor):

    default:
        switch in {
            // typical input handling
        }
    }
}
```

## Upcoming Features

These features are either in-progress or desired for the future

* Dedicated header and footer -- allow dedicated render funcs for top and bottom
  of app. This would limit the view funcs to the spaces in between
* Forms -- enable ease transformation from structs to string input fields and
  back for richer user-controlled state
