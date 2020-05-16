package tui

// Cursor is a tool for keeping track of state in a 2D array
// You need to provide a height and width so that the changes can be clamped
type Cursor struct {
	row, col, height, width int
}

func NewCursor(height, width int) Cursor {
	c := Cursor{}
	c.SetSize(height, width)
	return c
}

func (c *Cursor) Position() (int, int) {
	return c.row, c.col
}

func (c *Cursor) Size() (int, int) {
	return c.height, c.width
}

func (c *Cursor) Up() bool {
	if c.row > 0 {
		c.row--
		return true
	}
	return false
}

func (c *Cursor) Down() bool {
	if c.row < c.height {
		c.row++
		return true
	}
	return false
}

func (c *Cursor) Left() bool {
	if c.col > 0 {
		c.col--
		return true
	}
	return false
}

func (c *Cursor) Right() bool {
	if c.col < c.width {
		c.col++
		return true
	}
	return false
}

func (c *Cursor) Top() {
	c.row = 0
}

func (c *Cursor) Bottom() {
	c.row = c.height
}

func (c *Cursor) SetSize(height, width int) {
	if height > 0 {
		c.height = height - 1
	} else {
		c.height = 0
	}

	if c.row > c.height {
		c.row = c.height
	}

	if width > 0 {
		c.width = width - 1
	} else {
		c.width = 0
	}

	if c.col > c.width {
		c.col = c.width
	}
}

func (c *Cursor) SetPosition(row, col int) {
	if row < 0 {
		row = 0
	} else if row > c.height {
		row = c.height
	}

	c.row = row

	if col < 0 {
		col = 0
	} else if col > c.width {
		col = c.width
	}

	c.col = col
}
