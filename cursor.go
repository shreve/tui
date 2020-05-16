package tui

// Cursor is a tool for keeping track of state in a 2D array
// You need to provide a height and width so that the changes can be clamped
type Cursor struct {
	Row, Col, Height, Width int
}

func NewCursor(height, width int) Cursor {
	c := Cursor{}
	c.SetSize(height, width)
	return c
}

func (c *Cursor) Up() bool {
	if c.Row > 0 {
		c.Row--
		return true
	}
	return false
}

func (c *Cursor) Down() bool {
	if c.Row < (c.Height - 1) {
		c.Row++
		return true
	}
	return false
}

func (c *Cursor) Left() bool {
	if c.Col > 0 {
		c.Col--
		return true
	}
	return false
}

func (c *Cursor) Right() bool {
	if c.Col < c.Width {
		c.Col++
		return true
	}
	return false
}

func (c *Cursor) SetSize(height, width int) {
	c.Height = height
	c.Width = width
}
