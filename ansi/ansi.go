package ansi

import (
	"C"
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const (
	Black = 30
	Red = 31
	Green = 32
	Yellow = 33
	Blue = 34
	Magenta = 35
	Cyan = 36
	White = 37
)

const DisplayResetCode = "\033[0m"

const (
	// Clear entire screen
	clearScreen = "\033[2J"

	// Erase entire line regardless of cursor position
	clearLine = "\033[2K"

	// Erase from cursor, left
	clearLineLeft = "\033[1K"

	// Set cursor position (row, column), 1-indexed
	setCursorPos = "\033[%d;%dH"

	// Print the cursor position back into stdin
	getCursorPos = "\033[6n"

	// Stop displaying cursor
	showCursor = "\033[?25h"

	// Start showing cursor if stopped
	hideCursor = "\033[?25l"
)

func ClearScreen() {
	fmt.Print(clearScreen)
}

func ClearLine() {
	fmt.Print(clearLine)
}

func ClearRestOfLine() {
	fmt.Print(clearLineLeft)
}

func HideCursor() {
	fmt.Print(hideCursor);
}

func ShowCursor() {
	fmt.Print(showCursor);
}

// Set cursor position. If beyond size of terminal, behavior is undefined.
func MoveCursor(row, col int) {
	fmt.Printf(setCursorPos, row + 1, col + 1)
}

// Ask terminal for current cursor position
func GetCursor() (int, int) {

	// Print the query sequence
	os.Stdout.Write([]byte(getCursorPos))

	// Read the response sequence
	b := make([]byte, 15)
	os.Stdin.Read(b)

	// b is now \e[{ROW};{COL}R
	split := bytes.Index(b, []byte(";"))
	end := bytes.Index(b, []byte("R"))
	if split < 2 || end < 3 {
		return 0, 0
	}

	// Convert the read strings to integers
	row, _ := strconv.Atoi(string(b[2:split]))
	col, _ := strconv.Atoi(string(b[split+1:end]))
	return row, col
}

type winsize struct {
	Row, Col, Xpixel, Ypixel uint16
}

// Use ioctl to ask for size of terminal window
func WindowSize() (int, int) {
	win := winsize{}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, 0, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&win)))
	if err != 0 {
		panic(err)
	}
	return int(win.Row), int(win.Col)
}

type Display struct {
	Fg, Bg int
	Bright, Dim, Underscore, Blink, Reverse, Hidden bool
}

// Generate a display based on foreground and background colors.
func NewDisplay(fg, bg int) Display {
	d := Display{}
	d.Fg = fg
	d.Bg = bg
	return d
}

func (d Display) Code() string {
	return DisplayCode(d)
}

// Generate the escape sequence for a given display configuration
func DisplayCode(d Display) string {
	attrs := make([]string, 0)
	if d.Bright { attrs = append(attrs, "1") }
	if d.Dim { attrs = append(attrs, "2") }
	if d.Underscore { attrs = append(attrs, "4") }
	if d.Blink { attrs = append(attrs, "5") }
	if d.Reverse { attrs = append(attrs, "7") }
	if d.Hidden { attrs = append(attrs, "8") }
	if d.Fg != 0 { attrs = append(attrs, strconv.Itoa(d.Fg)) }
	if d.Bg != 0 { attrs = append(attrs, strconv.Itoa(d.Bg + 10)) }

	out := "\033["
	out += strings.Join(attrs, ";")
	out += "m"
	return out
}

// Print the escape sequence for a given display configuration
func SetDisplay(d Display) {
	fmt.Print(DisplayCode(d))
}

// Clear all output formatting
func ResetDisplay() {
	fmt.Print(DisplayResetCode)
}

// Tell the terminal to save the current output.
func SaveState() {
	cmd := exec.Command("tput", "smcup")
	cmd.Stdout = os.Stdin
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

// Tell the terminal to restore the previously saved output. This is useful for
// a full-window app that doesn't want to leave the terminal with a dead window
// upon close/exit.
func RestoreState() {
	cmd := exec.Command("tput", "rmcup")
	cmd.Stdout = os.Stdin
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
