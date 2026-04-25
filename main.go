package main

import (
	"fmt"
	"strings"

	"github.com/anirban1809/tuix/tuix"
	"github.com/anirban1809/tuix/tuix/components"
)

var tableHeaders = []string{"Name", "Role", "Status"}

var contacts = [][]string{
	{"Alice Chen", "Engineer", "● Active"},
	{"Bob Smith", "Designer", "○ Away"},
	{"Carol Wu", "Manager", "● Active"},
	{"Dave Kim", "Engineer", "○ Inactive"},
	{"Eva Park", "Designer", "● Active"},
}

func App(props tuix.Props) tuix.Element {
	focus, setFocus := tuix.UseState(0) // 0=search, 1=tabs, 2=table
	search, setSearch := tuix.UseState("")
	modalOpen, setModalOpen := tuix.UseState(false)

	// Tabs must render before filtering: Pass 2 calls onChange synchronously,
	// setting activeTab from Tabs' already-updated internal state.
	activeTab := 0
	tabs := components.Tabs(
		[]string{"All", "Active", "Away / Inactive"},
		focus == 1,
		func(i int) { activeTab = i },
	)

	// Filter: tab category first, then search string.
	q := strings.ToLower(search)
	filtered := make([][]string, 0, len(contacts))
	for _, row := range contacts {
		switch activeTab {
		case 1:
			if !strings.Contains(row[2], "Active") {
				continue
			}
		case 2:
			if strings.Contains(row[2], "Active") {
				continue
			}
		}
		for _, cell := range row {
			if q == "" || strings.Contains(strings.ToLower(cell), q) {
				filtered = append(filtered, row)
				break
			}
		}
	}

	// Tab cycles focus forward; ShiftTab cycles backward.
	if !modalOpen {
		if tuix.CurrentKey.Code == tuix.KeyTab {
			setFocus((focus + 1) % 3)
		}
		if tuix.CurrentKey.Code == tuix.KeyShiftTab {
			setFocus((focus + 2) % 3)
		}
	}

	// Search input keys (handled in App so we own the state for filtering).
	if focus == 0 && !modalOpen {
		switch tuix.CurrentKey.Code {
		case tuix.KeyBackspace:
			if r := []rune(search); len(r) > 0 {
				setSearch(string(r[:len(r)-1]))
			}
		case tuix.KeySpace:
			setSearch(search + " ")
		default:
			if tuix.CurrentKey.Rune != 0 {
				setSearch(search + string(tuix.CurrentKey.Rune))
			}
		}
	}

	// Enter on the table opens the modal.
	if focus == 2 && !modalOpen && tuix.CurrentKey.Code == tuix.KeyEnter && len(filtered) > 0 {
		setModalOpen(true)
	}

	// ── Search bar ─────────────────────────────────────────────────────────────
	cursor := ""
	if focus == 0 {
		cursor = "█"
	}
	labelStyle := tuix.NewStyle().Foreground(tuix.BrightBlack)
	inputStyle := tuix.NewStyle().Foreground(tuix.White)
	if focus == 0 {
		inputStyle = tuix.NewStyle().Foreground(tuix.Cyan).Bold(true)
	}
	searchBar := tuix.Box(
		tuix.Props{Direction: tuix.Row},
		tuix.NewStyle(),
		tuix.Text("Search: ", labelStyle),
		tuix.Text(search+cursor, inputStyle),
	)

	countLine := tuix.Text(
		fmt.Sprintf("%d of %d contacts", len(filtered), len(contacts)),
		labelStyle,
	)

	// ── Table ──────────────────────────────────────────────────────────────────
	selectedIdx := 0
	table := components.Table(tableHeaders, filtered, focus == 2, func(i int) {
		selectedIdx = i
	})

	// ── Modal ──────────────────────────────────────────────────────────────────
	var modalBody tuix.Element
	if len(filtered) > 0 && selectedIdx < len(filtered) {
		row := filtered[selectedIdx]
		modalBody = tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 1},
			tuix.NewStyle(),
			tuix.Text("Name:   "+row[0], tuix.NewStyle().Foreground(tuix.White).Bold(true)),
			tuix.Text("Role:   "+row[1], tuix.NewStyle().Foreground(tuix.White)),
			tuix.Text("Status: "+row[2], tuix.NewStyle().Foreground(tuix.White)),
		)
	} else {
		modalBody = tuix.Text("No contact selected.", labelStyle)
	}

	modal := components.Modal("Contact Details", modalOpen, 36, func() {
		setModalOpen(false)
	}, modalBody)

	// ── Hint bar ───────────────────────────────────────────────────────────────
	hint := tuix.Text(
		"Tab/ShiftTab: switch focus   ←/→: change tab   Enter: view details   Esc: quit",
		labelStyle,
	)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1},
		tuix.NewStyle(),
		searchBar,
		tabs,
		countLine,
		table,
		hint,
		modal,
	)
}

func main() {
	app := tuix.NewApp(200, 25)
	app.Run(App, tuix.Props{})
}
