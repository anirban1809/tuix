package components

import (
	"regexp"
	"strings"

	"github.com/anirban1809/tuix/tuix"
)

// ansiSequence matches CSI escape sequences (colors, cursor moves, etc).
// Pasted content from a colored terminal can carry these; we drop them so
// they don't render as literal "[42m" garbage in the input.
var ansiSequence = regexp.MustCompile(`\x1b\[[0-9;?]*[a-zA-Z]`)

// lineEndings normalizes clipboard line breaks to '\n'. macOS pastes use
// '\r', Windows uses '\r\n', and our renderer only respects '\n'.
var lineEndings = strings.NewReplacer("\r\n", "\n", "\r", "\n")

// sanitizePaste filters clipboard text for safe display in a text input:
// normalizes line endings, strips ANSI escape sequences, drops control
// characters except newline and tab (which the multiline renderer handles),
// and preserves all printable unicode.
func sanitizePaste(s string) string {
	s = lineEndings.Replace(s)
	s = ansiSequence.ReplaceAllString(s, "")
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return r
		}
		if r < 0x20 || r == 0x7F {
			return -1
		}
		return r
	}, s)
}

// Button renders a pressable label. Highlighted when focused.
func Button(label string, focused bool) tuix.Element {
	var style tuix.Style
	if focused {
		style = tuix.NewStyle().
			Foreground(tuix.Black).
			Background(tuix.Cyan).
			Bold(true)
	} else {
		style = tuix.NewStyle().Foreground(tuix.White)
	}
	return tuix.Text("[ "+label+" ]", style)
}

// Input renders a labeled text field. Shows a custom cursor when focused.
func Input(
	label string,
	focused bool,
	value string,
	onChange func(value string),
) tuix.Element {
	runes := []rune(value)

	// Cursor position is the insertion point: 0..len(runes)
	pos, setPos := tuix.UseState(len(runes))
	clamped := pos
	if clamped < 0 {
		clamped = 0
	}
	if clamped > len(runes) {
		clamped = len(runes)
	}
	if clamped != pos {
		setPos(clamped)
		pos = clamped
	}

	if focused {
		switch tuix.CurrentKey.Code {
		case tuix.KeyLeft:
			if pos > 0 {
				setPos(pos - 1)
			}
		case tuix.KeyRight:
			if pos < len(runes) {
				setPos(pos + 1)
			}
		case tuix.KeyBackspace:
			if pos > 0 {
				newRunes := append([]rune{}, runes[:pos-1]...)
				newRunes = append(newRunes, runes[pos:]...)
				onChange(string(newRunes))
				setPos(pos - 1)
			}
		case tuix.KeySpace:
			newRunes := append([]rune{}, runes[:pos]...)
			newRunes = append(newRunes, ' ')
			newRunes = append(newRunes, runes[pos:]...)
			onChange(string(newRunes))
			setPos(pos + 1)
		case tuix.KeyPaste:
			insert := []rune(sanitizePaste(tuix.CurrentKey.Paste))
			if len(insert) > 0 {
				newRunes := append([]rune{}, runes[:pos]...)
				newRunes = append(newRunes, insert...)
				newRunes = append(newRunes, runes[pos:]...)
				onChange(string(newRunes))
				setPos(pos + len(insert))
			}
		default:
			if tuix.CurrentKey.Rune != 0 {
				newRunes := append([]rune{}, runes[:pos]...)
				newRunes = append(newRunes, tuix.CurrentKey.Rune)
				newRunes = append(newRunes, runes[pos:]...)
				onChange(string(newRunes))
				setPos(pos + 1)
			}
		}
	}

	var fieldStyle tuix.Style
	if focused {
		fieldStyle = tuix.NewStyle().Foreground(tuix.White)
	} else {
		fieldStyle = tuix.NewStyle().Foreground(tuix.BrightBlack)
	}

	var display string
	if focused {
		if pos < len(runes) {
			display = string(runes[:pos]) + "█" + string(runes[pos+1:])
		} else {
			display = string(runes) + "█"
		}
	} else {
		display = value
	}

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Row,
			Width:     tuix.Grow(1),
			Align:     tuix.AlignStart,
		},
		tuix.NewStyle(),
		tuix.Text(label+" ", tuix.NewStyle().Foreground(tuix.White)),
		tuix.WrappedText(display, fieldStyle),
	)
}

// Checkbox renders a boolean toggle. Space or Enter toggles when focused.
func Checkbox(label string, focused bool, onChange func(bool)) tuix.Element {
	checked, setChecked := tuix.UseState(false)

	if focused {
		if tuix.CurrentKey.Code == tuix.KeySpace ||
			tuix.CurrentKey.Code == tuix.KeyEnter {
			setChecked(!checked)
		}
	}

	if onChange != nil {
		onChange(checked)
	}

	box := "[ ]"
	if checked {
		box = "[x]"
	}
	var style tuix.Style
	if focused {
		style = tuix.NewStyle().Foreground(tuix.Cyan).Bold(true)
	} else {
		style = tuix.NewStyle().Foreground(tuix.White)
	}
	return tuix.Text(box+" "+label, style)
}

// List renders a vertical item list with a cursor on the selected item.
// Up/Down arrows move the selection when focused.
func List(items []string, focused bool) tuix.Element {
	selected, setSelected := tuix.UseState(0)

	if focused {
		if tuix.CurrentKey.Code == tuix.KeyDown && selected < len(items)-1 {
			setSelected(selected + 1)
		}
		if tuix.CurrentKey.Code == tuix.KeyUp && selected > 0 {
			setSelected(selected - 1)
		}
	}

	children := make([]tuix.Element, len(items))
	for i, item := range items {
		prefix := "  "
		var style tuix.Style
		if i == selected {
			prefix = "> "
			if focused {
				style = tuix.NewStyle().
					Background(tuix.Blue).
					Foreground(tuix.Cyan).
					Bold(true)
			} else {
				style = tuix.NewStyle().Foreground(tuix.White).Bold(true)
			}
		} else {
			style = tuix.NewStyle().Foreground(tuix.BrightBlack)
		}
		children[i] = tuix.Text(prefix+item, style)
	}
	return tuix.Box(
		tuix.Props{Direction: tuix.Column},
		tuix.NewStyle(),
		children...)
}

// SelectPicker renders a single-line option cycler with < > arrows.
// Left/Right arrows cycle options when focused.
func SelectPicker(options []string, focused bool) tuix.Element {
	selected, setSelected := tuix.UseState(0)

	if focused {
		if tuix.CurrentKey.Code == tuix.KeyLeft && selected > 0 {
			setSelected(selected - 1)
		} else if tuix.CurrentKey.Code == tuix.KeyRight && selected < len(options)-1 {
			setSelected(selected + 1)
		}
	}

	label := options[selected]
	const optWidth = 12
	for len([]rune(label)) < optWidth {
		label += " "
	}
	var style tuix.Style
	if focused {
		style = tuix.NewStyle().Foreground(tuix.Cyan).Bold(true)
	} else {
		style = tuix.NewStyle().Foreground(tuix.White)
	}
	return tuix.Text("< "+label+" >", style)
}
