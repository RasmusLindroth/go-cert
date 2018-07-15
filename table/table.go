package table

import (
	"fmt"
	"io"
)

const (
	//AlignLeft aligns text left
	AlignLeft uint = 1 << iota
	//AlignRight aligns text right
	AlignRight
	//AlignCenter aligns text center
	AlignCenter
	//CenterHeader aligns first row center
	CenterHeader
)

// All format/color constants from https://github.com/fatih/color

//Attribute holds an int for formatting
type Attribute int

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

//InitTable returns a Table with alignment set
func InitTable(padding int, padChar string, alignment []uint) *Table {
	return &Table{Padding: padding, PadChar: padChar, Alignment: alignment}
}

//Table holds a table
type Table struct {
	Rows      []*Row
	Padding   int
	PadChar   string
	Alignment []uint
}

func (t *Table) AddRowStrings(s []string) {
	cols := []*Column{}
	for _, x := range s {
		cols = append(cols, &Column{Text: x})
	}
	t.Rows = append(t.Rows, &Row{Columns: cols})
}

func (t *Table) AddRow(c []*Column) {
	t.Rows = append(t.Rows, &Row{Columns: c})
}

func (t *Table) Print(w io.Writer) {
	width := t.columnWidths()
	lastIndex := len(width) - 1
	emptyColumn := &Column{Text: ""}
	for rowi, r := range t.Rows {
		for i, wh := range width {
			if i < len(r.Columns) {
				w.Write(t.formatColumn(i, wh, r.Columns[i], rowi == 0))
			} else {
				w.Write(t.formatColumn(i, wh, emptyColumn, rowi == 0))
			}

			if i < lastIndex {
				for x := 0; x < t.Padding; x++ {
					w.Write([]byte(t.PadChar))
				}
			}
		}
		w.Write([]byte("\n"))
	}
}

func (t *Table) formatColumn(columnIndex int, width int, c *Column, header bool) []byte {
	diff := width - c.length()
	align := 0

	if len(t.Alignment) > columnIndex {
		a := t.Alignment[columnIndex]

		if a&AlignLeft == AlignLeft {
			align = 0
		} else if a&AlignRight == AlignRight {
			align = 1
		} else if a&AlignCenter == AlignCenter {
			align = 3
		}

		if header && a&CenterHeader == CenterHeader {
			align = 3
		}
	}
	s := c.output()
	for i := 0; i < diff; i++ {
		if align == 0 {
			s = s + " "
		} else if align == 1 {
			s = " " + s
		} else if align == 3 {
			if i%2 == 0 {
				s = " " + s
			} else {
				s = s + " "
			}
		}
	}
	return []byte(s)
}

func (t *Table) columnWidths() []int {
	width := []int{}
	for _, r := range t.Rows {
		for i, c := range r.Columns {
			if len(width) < i+1 {
				width = append(width, c.length())
			} else if width[i] < c.length() {
				width[i] = c.length()
			}
		}
	}
	return width
}

//Row holds a row with columns
type Row struct {
	Columns []*Column
}

//Column holds an column
type Column struct {
	Text   string
	Format []Attribute
}

func (c *Column) length() int {
	return len(c.Text)
}

func (c *Column) output() string {
	if len(c.Format) == 0 {
		return c.Text
	}

	o := ""
	for _, a := range c.Format {
		o = o + addFormat(a)
	}
	o = o + c.Text
	o = o + addFormat(Reset)

	return o
}

func addFormat(format Attribute) string {
	return fmt.Sprintf("\033[%dm", format)
}
