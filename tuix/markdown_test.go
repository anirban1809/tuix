package tuix

import "testing"

func TestMarkdownFencedCodeWithLanguageCloses(t *testing.T) {
	lines := renderMarkdownLines("```go\nfmt.Println(\"hi\")\n```\nafter", 80, Style{})

	got := markdownLinesText(lines)
	want := []string{"  fmt.Println(\"hi\")", "", "after"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d: %#v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d: expected %q, got %q", i, want[i], got[i])
		}
	}
}

func TestMarkdownStrikethroughKeepsText(t *testing.T) {
	lines := renderMarkdownLines("keep ~~drop~~ done", 80, Style{})

	got := markdownLinesText(lines)
	want := []string{"keep drop done"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d: %#v", len(want), len(got), got)
	}
	if got[0] != want[0] {
		t.Fatalf("expected %q, got %q", want[0], got[0])
	}
}

func TestMarkdownIndentedListMarkersStayParagraph(t *testing.T) {
	lines := renderMarkdownLines("Command output:\n    - not a list\n    1. not ordered", 80, Style{})

	got := markdownLinesText(lines)
	want := []string{"Command output: - not a list 1. not ordered"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d: %#v", len(want), len(got), got)
	}
	if got[0] != want[0] {
		t.Fatalf("expected %q, got %q", want[0], got[0])
	}
}

func TestMarkdownListAllowsUpToThreeSpacesIndent(t *testing.T) {
	lines := renderMarkdownLines("   - item", 80, Style{})

	got := markdownLinesText(lines)
	want := []string{"• item"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d: %#v", len(want), len(got), got)
	}
	if got[0] != want[0] {
		t.Fatalf("expected %q, got %q", want[0], got[0])
	}
}

func markdownLinesText(lines []markdownLine) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		runes := make([]rune, 0, len(line))
		for _, cell := range line {
			runes = append(runes, cell.r)
		}
		out[i] = string(runes)
	}
	return out
}
