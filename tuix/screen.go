package tuix

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

type Screen struct {
	height int
	width  int
	out    io.Writer

	cells [][]Cell
	prev  [][]Cell

	dirty bool

	oldState *term.State

	// Physical terminal viewport, queried at Start. Bounds what Flush
	// can address with absolute cursor moves.
	termRows int
	termCols int

	// Absolute terminal row where this Screen's row 0 lives. Starts at
	// 1; EnsureRoom decreases it (possibly past 1) as the program grows
	// and older rows are pushed into scrollback.
	anchorRow int
}

type Cell struct {
	Rune  rune
	Style Style
	Wide  bool
}

func (s *Screen) GetCell(x int, y int) Cell {
	return s.cells[x][y]
}

func makeCellGrid(width, height int) [][]Cell {
	grid := make([][]Cell, width)
	for i := range grid {
		grid[i] = make([]Cell, height)
	}

	return grid
}

func NewScreenWriter(width int, height int, out io.Writer) *Screen {
	return &Screen{
		height: height,
		width:  width,
		out:    out,

		cells: makeCellGrid(width, height),
		prev:  makeCellGrid(width, height),
	}
}

func (s Screen) Width() int {
	return s.width
}

func (s Screen) Height() int {
	return s.height
}

func (s *Screen) Resize(width int, height int) {
	s.width = width
	s.height = height
	s.dirty = true

	s.cells = makeCellGrid(width, height)
	s.prev = makeCellGrid(width, height)
}

func (s *Screen) Start() {
	fmt.Print("\033[H\033[2J\033[3J") /* clear the screen */

	//enable raw mode
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	s.oldState = oldState

	if cols, rows, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		s.termCols = cols
		s.termRows = rows
	}
	s.anchorRow = 1

	//hide cursor
	fmt.Fprintf(s.out, "\033[?25l")
}

func (s Screen) Stop() {
	//show cursor
	fmt.Fprintf(s.out, "\033[?25h")

	//restore terminal old state
	term.Restore(int(os.Stdin.Fd()), s.oldState)
}

func (s *Screen) SetCell(x int, y int, value rune, style Style) {
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return
	}

	if runewidth.RuneWidth(value) == 2 {
		s.cells[x][y].Wide = true

		if x+1 < s.width {
			s.cells[x+1][y].Wide = true
		}
	}

	s.cells[x][y].Rune = value
	s.cells[x][y].Style = style
}

func (s *Screen) Flush() {
	for y := range s.height {
		absRow := s.anchorRow + y
		if absRow < 1 || absRow > s.termRows {
			// Row sits in scrollback or below the viewport — unrefreshable.
			continue
		}
		for x := range s.width {
			curr := s.cells[x][y]
			prev := s.prev[x][y]

			if curr == prev && !s.dirty {
				continue
			}

			fmt.Fprintf(s.out, "\033[%d;%dH", absRow, x+1)
			fmt.Fprintf(s.out, "%s%c\033[0m", curr.Style.ANSIPrefix(), curr.Rune)
			s.prev[x][y] = curr
		}
	}
	s.dirty = false
}

// EnsureRoom guarantees that the painted contentH rows fit within the
// physical terminal viewport. If the bottom of the program would sit
// below termRows, EnsureRoom writes the cells inline (with \r\n between
// rows) starting at the program's top — the terminal naturally scrolls
// older rows into scrollback as the cursor advances past the bottom row.
//
// MUST be called after the cell grid has been painted for this frame, as
// it reads from s.cells. After this call, prev is synced to cells for
// the rows that were just written, so a subsequent Flush is a no-op for
// them.
func (s *Screen) EnsureRoom(contentH int) {
	if s.termRows == 0 {
		// term.GetSize failed at Start; can't make scroll decisions.
		return
	}
	bottom := s.anchorRow + contentH - 1
	if bottom <= s.termRows {
		return
	}
	delta := bottom - s.termRows

	// Position cursor at the program's current top within the viewport
	// (or row 1 if the top is already in scrollback from a prior grow).
	topRow := s.anchorRow
	if topRow < 1 {
		topRow = 1
	}
	fmt.Fprintf(s.out, "\033[%d;1H", topRow)

	// Skip rows already in scrollback (negative absRow before this scroll).
	startY := 0
	if s.anchorRow < 1 {
		startY = 1 - s.anchorRow
	}

	for y := startY; y < contentH; y++ {
		for x := 0; x < s.width; x++ {
			cell := s.cells[x][y]
			r := cell.Rune
			if r == 0 {
				r = ' '
			}
			fmt.Fprintf(s.out, "%s%c\033[0m", cell.Style.ANSIPrefix(), r)
		}
		if y < contentH-1 {
			fmt.Fprint(s.out, "\r\n")
		}
	}

	// Sync prev = cells for the rows we just wrote so Flush won't re-emit them.
	for y := startY; y < contentH; y++ {
		for x := 0; x < s.width; x++ {
			s.prev[x][y] = s.cells[x][y]
		}
	}

	s.anchorRow -= delta
	s.dirty = false
}

func (s Screen) Clear() {
	for x := range s.width {
		for y := range s.height {
			s.cells[x][y].Rune = ' '
			s.cells[x][y].Wide = false
			s.cells[x][y].Style = Style{}
		}
	}
}

func RuneWidth(value rune) int {
	return runewidth.RuneWidth(value)
}

func (s *Screen) SetDimensions(width, height int) {
	s.width = width
	s.height = height
	s.cells = makeCellGrid(width, height)
	s.prev = makeCellGrid(width, height)
}

var StdOutScreen *Screen = &Screen{
	out: os.Stdout,
}
