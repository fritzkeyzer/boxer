package boxer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

//type Options struct {
//	Writer io.Writer
//}

//var DefaultOptions = Options{}

type BorderCfg struct {
	Top               string
	Bottom            string
	Left              string
	Right             string
	CornerTopLeft     string
	CornerTopRight    string
	CornerBottomLeft  string
	CornerBottomRight string

	ColumnSeperator  string
	ColumnJoinTop    string
	ColumnJoinBottom string
	ColumnJoin       string

	HeaderSeperator string
	HeaderJoinLeft  string
	HeaderJoinRight string
}

var DefaultBorderCfg = BorderCfg{
	Top:               "─",
	Bottom:            "─",
	Left:              "│",
	Right:             "│",
	CornerTopLeft:     "┌",
	CornerTopRight:    "┐",
	CornerBottomLeft:  "└",
	CornerBottomRight: "┘",
	ColumnSeperator:   "│",
	ColumnJoinTop:     "┬",
	ColumnJoinBottom:  "┴",
	ColumnJoin:        "┼",
	HeaderSeperator:   "─",
	HeaderJoinLeft:    "├",
	HeaderJoinRight:   "┤",
}

type row struct {
	cells []string
}

type Box struct {
	Writer io.Writer

	WrapLimit int

	ColumnsEnabled         bool
	HeaderSeperatorEnabled bool
	BorderEnabled          bool
	BorderCfg              BorderCfg
	Padding                int

	header       []string
	rows         []row
	columnWidths []int
	rowHeights   []int
}

// New returns a default box printer
//
//	id   name      value
//	1    foo       bar
//	2    example   multi
//	               line
func New() *Box {
	return &Box{
		Writer: os.Stdout,

		WrapLimit: 100,

		ColumnsEnabled:         false,
		HeaderSeperatorEnabled: false,
		BorderEnabled:          false,
		BorderCfg:              BorderCfg{},
		Padding:                3,
	}
}

// NewWithBorders
//
//	┌──┬───────┬─────┐
//	│id│name   │value│
//	├──┼───────┼─────┤
//	│1 │foo    │bar  │
//	│2 │example│multi│
//	│  │       │line │
//	└──┴───────┴─────┘
func NewWithBorders() *Box {
	return &Box{
		Writer: os.Stdout,

		WrapLimit: 100,

		ColumnsEnabled:         true,
		HeaderSeperatorEnabled: true,
		BorderEnabled:          true,

		BorderCfg: DefaultBorderCfg,
		Padding:   0,
	}
}

func (b *Box) AddHeader(values ...string) {
	for len(b.columnWidths) < len(values) {
		b.columnWidths = append(b.columnWidths, 0)
	}

	b.header = values

	for x, val := range values {
		vals := fmt.Sprint(val)

		// measure cell width
		split := strings.Split(vals, "\n")
		for _, s := range split {
			if len(s) > b.columnWidths[x] {
				b.columnWidths[x] = len(s)
			}
		}
	}
}

func (b *Box) AddLine(values ...any) {
	for len(b.columnWidths) < len(values) {
		b.columnWidths = append(b.columnWidths, 0)
	}

	var currentRow row
	for x, val := range values {
		vals := fmt.Sprint(val)
		currentRow.cells = append(currentRow.cells, vals)

		// measure cell width
		split := strings.Split(vals, "\n")
		for _, s := range split {
			if len(s) > b.columnWidths[x] {
				b.columnWidths[x] = len(s)
			}
		}
	}

	// add row to box
	b.rows = append(b.rows, currentRow)
}

func (b *Box) writeLine(tw *tabwriter.Writer, values ...string) {
	line := ""

	for i, valu := range values {
		left := ""
		if i > 0 {
			left = b.BorderCfg.ColumnSeperator
			left += strings.Repeat(" ", b.Padding)
		}

		line += fmt.Sprintf("%s%s\t", left, valu)
	}

	if b.BorderEnabled {
		line = b.BorderCfg.Left + line + b.BorderCfg.Right
	}

	fmt.Fprintln(tw, line)
}

// Print buffered content of Box and flush buffer.
func (b *Box) Print() {
	tw := tabwriter.NewWriter(b.Writer, 1, 2*b.WrapLimit, 0, ' ', 0)

	// calculate widths
	columnWidths := make([]int, len(b.columnWidths))
	//log.Println(columnWidths)
	totalWidth := 0
	for i, w := range b.columnWidths {
		totalWidth += w
		columnWidths[i] = w
	}
	if totalWidth > b.WrapLimit {
		columnWidths = make([]int, len(b.columnWidths))
		lim := b.WrapLimit
		i := 0
		for lim > 0 {
			if columnWidths[i]+1 <= b.columnWidths[i] {
				columnWidths[i]++
				lim--
			}

			i++
			if i >= len(b.columnWidths) {
				i = 0
			}
		}
	}

	// process rows
	var lines [][]string
	for i := range b.rows {
		type cell struct {
			lines []string
		}

		var cells []cell
		rowHeight := 1
		for x := range b.columnWidths {
			if len(b.rows[i].cells) > x {
				wrappedLines := wrapLine(b.rows[i].cells[x], columnWidths[x])
				cells = append(cells, cell{lines: wrappedLines})
				if len(wrappedLines) > rowHeight {
					rowHeight = len(wrappedLines)
				}
			} else {
				// append empty cell to end of row
				cells = append(cells, cell{})
			}
		}

		for y := 0; y < rowHeight; y++ {
			var line []string
			for x := range cells {
				text := ""
				if len(cells[x].lines) > y {
					text = cells[x].lines[y]
				}

				line = append(line, text)
			}
			lines = append(lines, line)
		}
	}

	// print top border
	if b.BorderEnabled {
		prevColSep := b.BorderCfg.ColumnSeperator
		prevBorderLeft := b.BorderCfg.Left
		prevBorderRight := b.BorderCfg.Right

		//b.ColumnSeperator = strings.Repeat(b.Top, len(prevColSep))
		b.BorderCfg.ColumnSeperator = b.BorderCfg.ColumnJoinTop
		b.BorderCfg.Left = b.BorderCfg.CornerTopLeft
		b.BorderCfg.Right = b.BorderCfg.CornerTopRight

		// write header
		var vals []string
		for i := 0; i < len(columnWidths); i++ {
			borderPiece := strings.Repeat(b.BorderCfg.Top, columnWidths[i])
			vals = append(vals, borderPiece)
		}
		b.writeLine(tw, vals...)

		// reset borders
		b.BorderCfg.ColumnSeperator = prevColSep
		b.BorderCfg.Left = prevBorderLeft
		b.BorderCfg.Right = prevBorderRight
	}

	// print header
	if len(b.header) > 0 {
		type cell struct {
			lines []string
		}

		var cells []cell
		rowHeight := 1
		for x := range columnWidths {

			if len(b.header) > x {
				wrappedLines := wrapLine(b.header[x], columnWidths[x])
				cells = append(cells, cell{lines: wrappedLines})
				if len(wrappedLines) > rowHeight {
					rowHeight = len(wrappedLines)
				}
			} else {
				// append empty cell to end of row
				cells = append(cells, cell{})
			}
		}

		for y := 0; y < rowHeight; y++ {
			var line []string
			for x := range cells {
				text := ""
				if len(cells[x].lines) > y {
					text = cells[x].lines[y]
				}

				line = append(line, text)
			}
			b.writeLine(tw, line...)
		}

		prevColSep := b.BorderCfg.ColumnSeperator
		prevBorderLeft := b.BorderCfg.Left
		prevBorderRight := b.BorderCfg.Right

		b.BorderCfg.ColumnSeperator = b.BorderCfg.ColumnJoin
		b.BorderCfg.Left = b.BorderCfg.HeaderJoinLeft
		b.BorderCfg.Right = b.BorderCfg.HeaderJoinRight

		if b.HeaderSeperatorEnabled {
			var vals []string
			for i := 0; i < len(columnWidths); i++ {
				borderPiece := strings.Repeat(b.BorderCfg.HeaderSeperator, columnWidths[i])
				vals = append(vals, borderPiece)
			}
			b.writeLine(tw, vals...)
		}

		b.BorderCfg.ColumnSeperator = prevColSep
		b.BorderCfg.Left = prevBorderLeft
		b.BorderCfg.Right = prevBorderRight
	}

	// print lines
	for _, line := range lines {
		b.writeLine(tw, line...)
	}

	// print bottom border
	if b.BorderEnabled {
		prevColSep := b.BorderCfg.ColumnSeperator
		prevBorderLeft := b.BorderCfg.Left
		prevBorderRight := b.BorderCfg.Right

		b.BorderCfg.ColumnSeperator = b.BorderCfg.ColumnJoinBottom
		b.BorderCfg.Left = b.BorderCfg.CornerBottomLeft
		b.BorderCfg.Right = b.BorderCfg.CornerBottomRight

		var vals []string
		for i := 0; i < len(columnWidths); i++ {
			vals = append(vals, strings.Repeat(b.BorderCfg.Bottom, columnWidths[i]))
		}
		b.writeLine(tw, vals...)

		// reset borders
		b.BorderCfg.ColumnSeperator = prevColSep
		b.BorderCfg.Left = prevBorderLeft
		b.BorderCfg.Right = prevBorderRight
	}

	tw.Flush()
}

func wrapLine(in string, limit int) []string {
	split := strings.Split(in, "\n")
	var out []string
	for i := range split {
		if len(split[i]) <= limit {
			out = append(out, split[i])
			continue
		}

		// chop line
		chopIndex := strings.LastIndex(split[i][:limit], " ")
		if chopIndex == -1 {
			chopIndex = limit
		}

		if chopIndex == 0 {
			out = append(out, split[i])
			continue
		}

		if chopIndex > len(split[i]) {
			chopIndex = len(split[i])
		}

		out = append(out, split[i][:chopIndex])
		remaining := split[i][chopIndex:]
		remaining = strings.TrimPrefix(remaining, " ")

		if remaining == "" {
			continue
		}

		otherLines := wrapLine(remaining, limit)
		out = append(out, otherLines...)
	}

	return out
}
