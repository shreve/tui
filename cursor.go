package tui

// Cursor is a tool for keeping track of state in a 2D array
// You need to provide a height and width so that the changes can be clamped
type Cursor struct {
	Row, Col, Height, Width int
}

func (c *Cursor) Up() {
	if c.Row > 0 {
		c.Row--;
	}
}

func (c *Cursor) Down() {
	if c.Row < c.Height {
		c.Row++;
	}
}

func (c *Cursor) Left() {
	if c.Col > 0 {
		c.Col--;
	}
}

func (c *Cursor) Right() {
	if c.Col < c.Width {
		c.Col++;
	}
}
