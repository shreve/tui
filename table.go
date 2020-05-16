//
// tui.Table
//
// Display, search, and select tabular data and render into a tui.View
//

// TODO: Replace Widths with Width, and auto-size columns
// TODO: Only collect values once when the list of records changes
// TODO: Create more natural "pushing" scrolling
// TODO: Extract style into config

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

	// Values are extracted from records. This is done once to avoid using
	// reflection more than necessary.
	values    []row

	// Are we searching, and what for? Run a strings.Contains query on each
	// value for a record to select.
	searching bool
	query     string

	// Go structs supplied to the table. Panic if these aren't structs.
	records  []interface{}

	// Which row of the table is selected?
	Cursor   Cursor

	// What are the names of the columns and how wide are they?
	Columns  []string
	Widths   []int

	// At what size are we able to render this table?
	Height   int
	// Width    int
}

// row is a stringified record, which points back to its entry in records
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

// Update the internal data of the table.
func (t *Table) Update(records interface{}) {

	// The supplied records must be a slice of structs
	s := reflect.ValueOf(records)
	if s.Kind() != reflect.Slice {
		panic("Non-slice supplied to tui.Table.Update")
	}
	t.records = make([]interface{}, s.Len())
	if s.Len() > 0 {
		if s.Index(0).Kind() != reflect.Struct {
			panic("Slice of non-structs was supplied to tui.Table.Update")
		}
	}
	for i := 0; i < s.Len(); i++ {
		t.records[i] = s.Index(i).Interface()
	}
}

// Compute the table into a View
func (t *Table) Draw() (view View) {
	t.Fetch()
	view = append(view, t.Heading())
	return append(view, t.Body()...)
}

var titleDisplay = ansi.DisplayCode(ansi.NewDisplay(ansi.Black, ansi.White))
var highlightedDisplay = ansi.DisplayCode(ansi.NewDisplay(ansi.Black, ansi.Yellow))

// Draw the column names
func (t *Table) Heading() string {
	out := bytes.NewBufferString(titleDisplay)
	for i := 0; i < len(t.Columns); i++ {
		out.WriteString(rightPad(t.Columns[i], t.Widths[i]))
	}
	out.WriteString(ansi.DisplayResetCode)
	return out.String()
}

// Pull out the data from records based on column names
func (t *Table) Body() View {
	out := make(View, t.Height)

	// Provided height includes heading. Body height is one less.
	height := t.Height - 1

	// When searching, we append a line about the search, so the viewport is one
	// row shorter.
	if t.searching {
		height -= 1
	}

	// If we can't fill the whole space, only fill with what we have.
	if len(t.values) < height {
		height = len(t.values)
	}

	// If the current selection is beyond the height of our viewport, we need to
	// use an offset to shift our contents so the selection is in view.
	selected, _ := t.Cursor.Position()
	offset := 0
	if selected > (height - 2) {
		offset = selected - height + 1
	}

	// For the height of our viewport:
	for i := 0; i < height; i++ {

		// Pick a row of data.
		index := i + offset
		line := bytes.NewBufferString(ansi.DisplayResetCode)

		// If it's selected, highlight it.
		if index == selected {
			line.WriteString(highlightedDisplay)
		}

		// Write out the content for each column.
		for j := 0; j < len(t.Columns); j++ {
			line.WriteString(
				rightPad(
					fmt.Sprintf("%s", t.values[index].columns[j]),
					t.Widths[j]))
		}

		// Reset the style and save the line
		line.WriteString(ansi.DisplayResetCode)
		out[i] = line.String()
	}

	// If we are currently searching, add info about the search to the bottom
	if t.searching {
		out[len(out) - 1] = ansi.DisplayResetCode +
			fmt.Sprintf(" Searching For \"%s\"", t.query)
	}

	return out
}

// Pull strings out of our []interface{} records and perform our search
func (t *Table) Fetch() {

	// Reset the collection
	t.values = make([]row, 0)

	// Get our search query ready (lower to lower comparison)
	needle := strings.ToLower(t.query)

	// For each row:
	for i := 0; i < len(t.records); i++ {

		// Save the record it came from.
		row := row{i, make([]string, len(t.Columns))}

		matched := false

		// For each column:
		for j := 0; j < len(t.Columns); j++ {

			// Get the value from the record by this column's name
			value := reflect.ValueOf(t.records[i]).FieldByName(t.Columns[j])

			// Cast the value to a string.
			row.columns[j] = fmt.Sprintf("%v", value)

			// If we are searching, see if this value matches the query.
			if !t.searching ||
				strings.Contains(strings.ToLower(row.columns[j]), needle) {
				matched = true
			}
		}

		// Add the found row to the list of values
		if matched {
			t.values = append(t.values, row)
		}
	}

	// Bound our cursor to the potentially newly modified list
	t.Cursor.SetSize(len(t.values), 1)
}

// The search is performed on Draw(), but this informs us we need to perform it.
func (t *Table) Search(query string) {
	t.searching = true
	t.query = query
}

// If we're keeping this table around, we need to be able to clear search mode.
func (t *Table) ClearSearch() {
	t.searching = false
}

// Returns the index of the record associated with the currently selected value.
func (t *Table) SelectedRecord() int {
	selected, _ := t.Cursor.Position()
	return t.values[selected].recordIndex
}

// Make a table-cell-style string out of an input to be a given total length
func rightPad(input string, length int) string {
	for utf8.RuneCountInString(input) < (length - 2) {
		input += " "
	}
	out := bytes.NewBufferString(" ")
	for utf8.RuneCountInString(out.String()) < length-2 {
		r, size := utf8.DecodeRuneInString(input)
		out.WriteRune(r)
		input = input[size:]
	}
	out.WriteString(" ")
	return out.String()
}
