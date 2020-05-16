tui
===

A library for generating simple terminal user interface applications.

This tool primarily handles:
 * Reading input
 * Fast rendering of views
 * Setting the output mode of the terminal for easy ANSI code interaction

There are more features in the works to quickly develop interface components.

## Usage

```go
package main

import (
    "github.com/shreve/tui"
)

func handleInput(input string, app *tui.App) {
    switch input[0] {
    case "q", tui.CtrlC:
        // Tell the app to gracefully stop
        app.Done()
    }

    // Tell the app to re-draw
    app.Redraw()
}

func indexView(height, width int) tui.View {
    view := make(tui.View, height)

    // Draw your app into `view`

    return view
}

func main() {
    app := tui.NewApp()
    app.InputHandler = handleInput
    app.CurrentView = indexView
    app.Run()
}
```

This library also includes a simple tool for keeping track of cursor state.
Provide dimensions of the space, then call Up, Down, Left, Right to move
around within the space.

```go
cursor := tui.NewCursor(3, 1)
cursor.Up() // false
cursor.Down() // true
cursor.Down() // true
cursor.Down() // false
cursor.Up() // true
cursor.Left() // false
cursor.Right() // false
```

## Upcoming Features

These features are either in-progress or desired for the future

* Tables -- working on a type that can take a slice of structs and turn it into
  a list view, including scrolling via the cursor
* Dedicated header and footer -- allow dedicated render funcs for top and bottom
  of app. This would limit the view funcs to the spaces in between
* Better state structure -- more guidance in place for managing app state {
  id, view, input handler, transition event }
