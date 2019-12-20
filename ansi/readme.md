tui/ansi
========

A sub-library of tui for generating ANSI escape sequences.

This is still in early development and the API is unstable.

```go
ansi.ClearScreen()

ansi.ClearLine()

ansi.MoveCursor(0, 0)

row, col := ansi.GetCursor()

height, width := ansi.WindowSize()

ansi.HideCursor()

defer ansi.ShowCursor()

ansi.SaveState()

defer ansi.RestoreState()

display := ansi.NewDisplay(ansi.Black, ansi.Yellow)

ansi.SetDisplay(display)

ansi.ResetDisplay()
```
