package tui

import (
	"bytes"
	"fmt"
	"github.com/shreve/tui/ansi"
	"reflect"
	"unicode/utf8"
	"strings"
)

// Table is a structure for drawing tabular data. Data is any slice of structs.
// Supply column names and widths to pull data out of records and draw.
type Table struct {
	values    []row
	searching bool
	query     string

	Cursor   Cursor
	Columns  []string
	Widths   []int
	Records  []interface{}
	Height   int
	Width    int
}

type row struct {
	recordIndex int
	columns     []string
}

func NewTable(records interface{}, height int) *Table {
	t := Table{}
	t.Height = height
	t.Update(records)
	return &t
}

func (t *Table) Update(records interface{}) {
	s := reflect.ValueOf(records)
	if s.Kind() != reflect.Slice {
		panic("Non-slice supplied to tui.Table.Update")
	}
	t.Records = make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		t.Records[i] = s.Index(i).Interface()
	}
}

func (t *Table) Draw() View {
	t.Fetch()
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
func (t *Table) Body() View {
	out := make(View, t.Height)

	// Provided height includes heading
	h := t.Height - 1
	if len(t.values) < h {
		h = len(t.values)
	}

	selected, _ := t.Cursor.Position()

	offset := 0
	if selected > (h - 2) {
		offset = selected - h + 1
	}

	for i := 0; i < h; i++ {
		line := ansi.DisplayResetCode
		if i + offset == selected {
			line += highlightedDisplay
		}
		for j := 0; j < len(t.Columns); j++ {
			line += rightPad(fmt.Sprintf("%s", t.values[i + offset].columns[j]), t.Widths[j])
		}
		if i == selected {
			line += ansi.DisplayResetCode
		}
		out[i] = line
	}

	if t.searching {
		out[len(out) - 1] = ansi.DisplayResetCode + fmt.Sprintf(" Searching For \"%s\"", t.query)
	}

	return out
}

func (t *Table) Fetch() {
	t.values = make([]row, 0)
	for i := 0; i < len(t.Records); i++ {
		row := row{i, make([]string, len(t.Columns))}
		matched := false
		needle := strings.ToLower(t.query)

		for j := 0; j < len(t.Columns); j++ {
			value := reflect.ValueOf(t.Records[i]).FieldByName(t.Columns[j])
			row.columns[j] = fmt.Sprintf("%v", value)

			if !t.searching ||
				strings.Contains(strings.ToLower(row.columns[j]), needle) {
				matched = true
			}
		}

		if matched {
			t.values = append(t.values, row)
		}
	}

	t.Cursor.SetSize(len(t.values), 1)
}

func (t *Table) Search(query string) {
	t.searching = true
	t.query = query
}

func (t *Table) ClearSearch() {
	t.searching = false
}

func (t *Table) SelectedRecord() int {
	selected, _ := t.Cursor.Position()
	return t.values[selected].recordIndex
}

// Make a table-cell-style string out of an input to be a given total length
func rightPad(input string, length int) string {
	for utf8.RuneCountInString(input) < (length - 2) {
		input += " "
	}
	out := bytes.Buffer{}
	for utf8.RuneCountInString(out.String()) < length-2 {
		r, size := utf8.DecodeRuneInString(input)
		out.WriteRune(r)
		input = input[size:]
	}
	return " " + out.String() + " "
}
