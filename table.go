//
// tui.Table
//
// Display, search, and select tabular data and render into a tui.View
//
// General data flow:
// 1. Update -> add new records
// 2. Fetch  -> capture string values from records
// 3. Search -> filter values into results
// 4. Draw   -> render results into table
//

// TODO: Extract style into config
// TODO: Extract widths generation to be recalled upon Width change

package tui

import (
	"bytes"
	"fmt"
	"github.com/shreve/tui/ansi"
	"reflect"
	"strings"
	"sync"
	"unicode/utf8"
)

// Table is a structure for drawing tabular data. Data is any slice of structs.
// Supply column names and widths to pull data out of records and draw.
type Table struct {

	// Values are extracted from records. This is done once to avoid using
	// reflection more than necessary.
	values []row

	// Results are indices of values returned from a search.
	results []int

	// Are we searching, and what for? Run a strings.Contains query on each
	// value for a record to select.
	searching bool
	query     string

	// Go structs supplied to the table. Panic if these aren't structs.
	records []interface{}

	// Generated widths for columns based on content length and Table width
	widths []int

	// General lock for multi-threaded weirdness
	lock sync.Mutex

	// Which row of the table is selected?
	Cursor Cursor

	// What are the names of the columns?
	Columns []string

	// At what size are we able to render this table?
	Height int
	Width  int
}

// row is a stringified record, which points back to its entry in records
type row struct {
	recordIndex int
	columns     []string
}

// Update the internal data of the table.
func (t *Table) Update(records interface{}, columns []string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.Columns = columns

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

	if len(t.records) == 0 {
		return
	}

	// Pull strings out of our []interface{} records

	// Reset the collection
	t.values = make([]row, len(t.records))

	lengths := make([]int, len(t.Columns))

	// For each row:
	for i := 0; i < len(t.records); i++ {

		// Save the record it came from.
		row := row{i, make([]string, len(t.Columns))}

		// For each column:
		for j := 0; j < len(t.Columns); j++ {

			// Get the value from the record by this column's name
			value := reflect.ValueOf(t.records[i]).FieldByName(t.Columns[j])

			// Cast the value to a string.
			row.columns[j] = fmt.Sprintf("%v", value)

			lengths[j] += len(row.columns[j])
		}

		// Add the found row to the list of values
		t.values[i] = row
	}

	total_length := 0
	for i := 0; i < len(lengths); i++ {
		lengths[i] /= len(t.records)
		total_length += lengths[i]
	}

	t.widths = make([]int, len(t.Columns))

	for i := 0; i < len(lengths); i++ {
		pct := float32(lengths[i]) / float32(total_length)
		t.widths[i] = int(pct * float32(t.Width))
	}

	for i := 0; sum(t.widths) <= (t.Width + 1); i++ {
		t.widths[i % len(t.widths)]++
	}

	// Bound our cursor to the potentially newly modified list
	t.Cursor.SetSize(len(t.values), 1)

	// Reset the filter
	t.resetResults()
}

func (t *Table) Search(query string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.searching = true
	t.query = query

	t.results = make([]int, 0)

	// Get our search query ready (lower to lower comparison)
	needle := strings.ToLower(t.query)

	for i := 0; i < len(t.values); i++ {

		matched := false

		for j := 0; j < len(t.Columns); j++ {
			haystack := strings.ToLower(t.values[i].columns[j])

			// If we are searching, see if this value matches the query.
			if strings.Contains(haystack, needle) {
				matched = true
			}
		}

		if matched {
			t.results = append(t.results, i)
		}
	}

	// Bound our cursor to the potentially newly modified list
	t.Cursor.SetSize(len(t.results), 1)
}

// Compute the table into a View
func (t *Table) Draw() (view View) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if len(t.records) == 0 {
		return make(View, t.Height)
	}

	view = append(view, t.Heading())
	return append(view, t.Body()...)
}

var titleDisplay = ansi.DisplayCode(ansi.NewDisplay(ansi.Black, ansi.White))
var highlightedDisplay = ansi.DisplayCode(ansi.NewDisplay(ansi.Black, ansi.Yellow))

// Draw the column names
func (t *Table) Heading() string {
	out := bytes.NewBufferString(titleDisplay)
	for i := 0; i < len(t.Columns); i++ {
		out.WriteString(rightPad(t.Columns[i], t.widths[i]))
	}
	out.WriteString(ansi.DisplayResetCode)
	return out.String()
}

// Pull out the data from records based on column names
func (t *Table) Body() View {
	out := make(View, t.Height)

	if len(t.records) == 0 {
		return out
	}

	// Provided height includes heading. Body height is one less.
	height := t.Height - 1

	// When searching, we append a line about the search, so the viewport is one
	// row shorter.
	if t.searching {
		height -= 1
	}

	// If we can't fill the whole space, only fill with what we have.
	if len(t.results) < height {
		height = len(t.results)
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
			value := t.values[t.results[index]].columns[j]
			line.WriteString(rightPad(fmt.Sprintf("%s", value), t.widths[j]))
		}

		// Reset the style and save the line
		line.WriteString(ansi.DisplayResetCode)
		out[i] = line.String()
	}

	// If we are currently searching, add info about the search to the bottom
	if t.searching {
		searchLine := bytes.NewBufferString(ansi.DisplayResetCode)
		searchLine.WriteString(titleDisplay)
		searchLine.WriteString(
			rightPad(fmt.Sprintf(" Searching For \"%s\"", t.query), t.Width))
		out[len(out)-1] = searchLine.String()
	}

	return out
}

// If we're keeping this table around, we need to be able to clear search mode.
func (t *Table) ClearSearch() {
	t.searching = false
}

// Returns the index of the record associated with the currently selected value.
func (t *Table) SelectedRecord() int {
	selected, _ := t.Cursor.Position()
	return t.values[t.results[selected]].recordIndex
}

func (t *Table) resetResults() {
	t.results = make([]int, len(t.values))
	for i := 0; i < len(t.values); i++ {
		t.results[i] = i
	}
}

// Make a table-cell-style string out of an input to be a given total length
// We use special utf8 functions rather than just len() here because multi-byte
// characters mess up alignment. We want to do our best to have `length` visible
// runes rather than just bytes.
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

func sum(nums []int) (n int) {
	for _, i := range nums {
		n += i
	}
	return
}
