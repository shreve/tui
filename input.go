package tui

const (
	KeyEsc       = "\x1b"
	KeyUp        = "\x1b[A"
	KeyDown      = "\x1b[B"
	KeyLeft      = "\x1b[D"
	KeyRight     = "\x1b[C"
	KeyDelete    = "\x1b[3~"
	KeyBackspace = "\u007f"
	CtrlA        = "\x01"
	CtrlB        = "\x02"
	CtrlC        = "\x03"
	Enter        = "\r"
)

type Inputable interface {
	InputHandler(string)
}

const (
	WasdCursor = iota
	ArrowCursor
	ViCursor
)

func InputMoveCursor(mode int, input string, cursor *Cursor) bool {
	switch mode {
	case WasdCursor:
		switch input {
		case "w":
			cursor.Up()
		case "a":
			cursor.Left()
		case "s":
			cursor.Down()
		case "d":
			cursor.Right()
		default:
			return false
		}
	case ArrowCursor:
		switch input {
		case KeyUp:
			cursor.Up()
		case KeyLeft:
			cursor.Left()
		case KeyDown:
			cursor.Down()
		case KeyRight:
			cursor.Right()
		default:
			return false
		}
	case ViCursor:
		switch input {
		case "h":
			cursor.Left()
		case "j":
			cursor.Down()
		case "k":
			cursor.Up()
		case "l":
			cursor.Right()
		default:
			return false
		}

	default:
		return false
	}
	return true
}
