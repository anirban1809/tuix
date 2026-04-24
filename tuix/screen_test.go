package tuix_test

import (
	"bytes"
	"strings"
	"testing"

	tuix "github.com/anirban1809/tuix/tuix"
)

func newTestScreen(w, h int) (*tuix.Screen, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	s := tuix.NewScreenWriter(w, h, buf)
	return s, buf
}

// --- Cell storage ---

func TestSetCellStoresCorrectly(t *testing.T) {
	s, _ := newTestScreen(10, 5)
	s.SetCell(3, 2, 'X', tuix.Style{}.Bold(true))

	cell := s.GetCell(3, 2)
	if cell.Rune != 'X' {
		t.Errorf("got rune %q, want 'X'", cell.Rune)
	}
	if !cell.Style.IsBold() {
		t.Error("cell style should be bold")
	}
}

func TestSetCellOutOfBoundsIsNoOp(t *testing.T) {
	s, _ := newTestScreen(10, 5)
	s.SetCell(-1, 0, 'X', tuix.Style{})
	s.SetCell(10, 0, 'X', tuix.Style{})
	s.SetCell(0, -1, 'X', tuix.Style{})
	s.SetCell(0, 5, 'X', tuix.Style{})
}

func TestClearResetsAllCells(t *testing.T) {
	s, _ := newTestScreen(5, 5)
	s.SetCell(2, 2, 'A', tuix.Style{}.Bold(true))
	s.Clear()

	cell := s.GetCell(2, 2)
	if cell.Rune != 0 && cell.Rune != ' ' {
		t.Errorf("after Clear, expected empty cell, got %q", cell.Rune)
	}
	if cell.Style != (tuix.Style{}) {
		t.Error("after Clear, cell style should be zero value")
	}
}

// --- Flush / diffing ---

func TestFlushWritesChangedCellsOnly(t *testing.T) {
	s, buf := newTestScreen(10, 5)

	s.SetCell(0, 0, 'A', tuix.Style{})
	s.SetCell(1, 0, 'B', tuix.Style{})
	s.Flush()
	firstWrite := buf.String()

	buf.Reset()

	s.SetCell(1, 0, 'C', tuix.Style{})
	s.Flush()
	secondWrite := buf.String()

	if len(secondWrite) >= len(firstWrite) {
		t.Errorf("second flush (%d bytes) should be smaller than first (%d bytes)", len(secondWrite), len(firstWrite))
	}
	if !strings.Contains(secondWrite, "C") {
		t.Error("second flush should contain the changed rune 'C'")
	}
	if strings.Contains(secondWrite, "A") {
		t.Error("second flush should not redraw unchanged cell 'A'")
	}
}

func TestFlushAfterNoChangeWritesNothing(t *testing.T) {
	s, buf := newTestScreen(10, 5)

	s.SetCell(0, 0, 'A', tuix.Style{})
	s.Flush()
	buf.Reset()

	s.Flush()
	if buf.Len() != 0 {
		t.Errorf("flush with no changes wrote %d bytes, want 0", buf.Len())
	}
}

func TestFlushPositionsCorrectly(t *testing.T) {
	s, buf := newTestScreen(10, 5)

	s.SetCell(4, 3, 'Z', tuix.Style{})
	s.Flush()

	// ANSI cursor: row and col are 1-indexed → row=4, col=5
	if !strings.Contains(buf.String(), "\033[4;5H") {
		t.Errorf("expected cursor move \\033[4;5H, got: %q", buf.String())
	}
}

func TestFlushClearsRemovedCell(t *testing.T) {
	s, buf := newTestScreen(10, 5)

	s.SetCell(2, 2, 'X', tuix.Style{})
	s.Flush()
	buf.Reset()

	s.SetCell(2, 2, ' ', tuix.Style{})
	s.Flush()

	if !strings.Contains(buf.String(), " ") {
		t.Error("clearing a cell should write a space")
	}
}

// --- Unicode width ---

func TestRuneWidth(t *testing.T) {
	cases := []struct {
		r     rune
		width int
	}{
		{'A', 1},
		{'€', 1},
		{'中', 2},
		{'🐱', 2},
	}

	for _, tt := range cases {
		got := tuix.RuneWidth(tt.r)
		if got != tt.width {
			t.Errorf("RuneWidth(%q) = %d, want %d", tt.r, got, tt.width)
		}
	}
}

func TestWideRuneOccupiesTwoCells(t *testing.T) {
	s, _ := newTestScreen(10, 5)
	s.SetCell(0, 0, '中', tuix.Style{})

	next := s.GetCell(1, 0)
	if !next.Wide {
		t.Error("cell after wide rune should be marked as Wide continuation")
	}
}

func TestWideRuneAtEdgeDoesNotPanic(t *testing.T) {
	s, _ := newTestScreen(5, 5)
	s.SetCell(4, 0, '中', tuix.Style{})
}

// --- Screen dimensions ---

func TestScreenDimensions(t *testing.T) {
	s, _ := newTestScreen(80, 24)
	if s.Width() != 80 || s.Height() != 24 {
		t.Errorf("got %dx%d, want 80x24", s.Width(), s.Height())
	}
}

func TestResizeReallocatesBuffer(t *testing.T) {
	s, _ := newTestScreen(10, 5)
	s.SetCell(5, 3, 'A', tuix.Style{})

	s.Resize(20, 10)

	if s.Width() != 20 || s.Height() != 10 {
		t.Errorf("after resize got %dx%d, want 20x10", s.Width(), s.Height())
	}
	s.SetCell(19, 9, 'B', tuix.Style{})
	s.Flush()
}

func TestResizeMarksDirty(t *testing.T) {
	s, buf := newTestScreen(10, 5)
	s.SetCell(0, 0, 'A', tuix.Style{})
	s.Flush()
	buf.Reset()

	s.Resize(10, 5)
	s.Flush()

	if buf.Len() == 0 {
		t.Error("resize should mark all cells dirty and force full redraw")
	}
}
