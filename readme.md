tui
===

A library for generating simple terminal user interface applications.

This tool primarily handles:
 * Reading input
 * Fast rendering of views
 * Setting the output mode of the terminal for easy ANSI code interaction

There are more features in the works to quickly develop interface components.
Currently in dev is a table for rendering generic tabular data.

## Usage

```go
package main

import (
    "github.com/shreve/tui"
    "github.com/shreve/tui/ansi"
)

func handleInput(input []byte, app *tui.App) {
    switch input[0] {
    case 3, 113: // ctrl-c, q
        // Tell the app to gracefully stop
        app.Done()
    }

    // Tell the app to re-draw
    app.Redraw()
}

func indexView() tui.View {
    view := make(tui.View, 0)
    height, width := ansi.WindowSize()

    view[0] = "App Titlebar"

    for i := 1; i < height; i++ {
       view[i] = " * line content"
    }

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
