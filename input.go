package tui

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

type InputHandler func(string, *App)
