package tui

import (
	"fmt"
	"bytes"
	"reflect"
	"unicode/utf8"
	"github.com/shreve/tui/ansi"
)

// Table is a structure for drawing tabular data. Data is any slice of structs.
// Supply column names and widths to pull data out of records and draw.
type Table struct {
	Columns []string
	Widths []int
	Records []interface{}
	Height int
	Selected int
}

func NewTable(records interface{}) *Table {
	t := Table{}
	s := reflect.ValueOf(records)
	if s.Kind() != reflect.Slice {
		panic("Non-slice supplied to tui.NewTable")
	}
	t.Records = make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		t.Records[i] = s.Index(i).Interface()
	}
	return &t
}

func (t *Table) Draw() View {
	view := make(View, 0)
	view = append(view, t.Heading())
	return append(view, t.Body()...)
}

var titleDisplay = ansi.DisplayCode(ansi.NewDisplay(ansi.Black, ansi.White))
var highlightedDisplay = ansi.DisplayCode(ansi.NewDisplay(ansi.Black, ansi.Yellow))

// Draw the column names
func (t *Table) Heading() string {
	line := titleDisplay
	for i := 0; i < len(t.Columns); i++ {
		line += rightPad(t.Columns[i], t.Widths[i])
	}
	return line + ansi.DisplayResetCode
}

// Pull out the data from records based on column names
func (t *Table) Body() []string {
	out := make([]string, 0)
	h := t.Height
	if len(t.Records) < h { h = len(t.Records) }

	for i := 0; i < h;  i++ {
		line := ""
		if i == t.Selected {
			line += highlightedDisplay
		}
		for j := 0; j < len(t.Columns); j++ {
			value := reflect.ValueOf(t.Records[i]).FieldByName(t.Columns[j])
			width := t.Widths[j]
			line += rightPad(fmt.Sprintf("%v", value), width)
		}
		if i == t.Selected {
			line += ansi.DisplayResetCode
		}
		out = append(out, line)
	}

	return out
}

// Make a table-cell-style string out of an input to be a given total length
func rightPad(input string, length int) string {
	for utf8.RuneCountInString(input) < (length - 2) {
		input += " "
	}
	out := bytes.Buffer{}
	for utf8.RuneCountInString(out.String()) < length - 2 {
		r, size := utf8.DecodeRuneInString(input)
		out.WriteRune(r)
		input = input[size:]
	}
	return " " + out.String() + " "
}
