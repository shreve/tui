package tui_test

import (
	"github.com/shreve/tui"
	"testing"
)

func TestCursor(t *testing.T) {
	cursor := tui.NewCursor(1, 1)

	if cursor.Up() {
		t.Error("Went up from start")
	}

	if cursor.Left() {
		t.Error("Went left from start")
	}

	if cursor.Down() {
		t.Error("Went down a list of one")
	}

	if cursor.Right() {
		t.Error("Went right over a list of one")
	}

	cursor.SetSize(0, 0)

	if cursor.Down() {
		t.Error("Went down an empty list")
	}

	if cursor.Right() {
		t.Error("Went right over a list of one")
	}

	cursor.SetSize(-1, -1)
	if row, col := cursor.Size(); row != 0 || col != 0 {
		t.Error("Set size to a negative value")
	}

	cursor.SetSize(5, 5)
	cursor.SetPosition(5, 5)
	if row, col := cursor.Position(); row != 4 || col != 4 {
		t.Error("Set position beyond size of field")
	}

	cursor.SetSize(1, 1)
	if row, col := cursor.Position(); row != 0 || col != 0 {
		t.Error("Shrunk field beneath position")
	}
}
