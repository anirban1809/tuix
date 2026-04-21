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
	out    *io.Writer

	cells [][]Cell
	prev  [][]Cell

	dirty bool

	oldState *term.State
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
		out:    &out,

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
	//enable raw mode
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	s.oldState = oldState

	//hide cursor
	fmt.Fprintf(*s.out, "\033[?25l")
}

func (s Screen) Stop() {
	//show cursor
	fmt.Fprintf(*s.out, "\033[?25h")

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
		for x := range s.width {
			curr := s.cells[x][y]
			prev := s.prev[x][y]

			if curr == prev && !s.dirty {
				continue
			}

			fmt.Fprintf(*s.out, "\033[%d;%dH", y+1, x+1)
			fmt.Fprintf(*s.out, "%s%c\033[0m", curr.Style.ANSIPrefix(), curr.Rune)
			s.prev[x][y] = curr
		}
	}
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
