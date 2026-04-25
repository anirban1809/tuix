# tuix

A Go framework for building interactive terminal UIs with React-style components and hooks.

Tuix brings the declarative, composable model of React — functional components, `UseState`, `UseEffect`, flexbox layout — to the terminal. You describe *what* your UI looks like; tuix handles measuring, laying out, rendering, and diffing only the cells that changed.

```
┌─────────────────────────────────────────────────┐
│  Contacts                                        │
│  Search: john_                                   │
│ ┌──────────┬──────────┬──────────┐               │
│ │  All     │  Active  │  Away    │               │
│ ├──────────┴──────────┴──────────┤               │
│ │ Name          Status   Email   │               │
│ │ John Doe    ● Active   j@…     │               │
│ │ Jane Smith  ○ Away     s@…     │               │
│ └────────────────────────────────┘               │
└─────────────────────────────────────────────────┘
```

---

## Features

- **Functional components** — plain Go functions that return an `Element` tree
- **Hooks** — `UseState` and `UseEffect` with dependency tracking
- **Flexbox layout engine** — two-pass (measure → layout) with `Row`/`Column` direction, `Gap`, `Padding`, `Align`, `Justify`, and `Fixed`/`Grow`/`Fit` sizing
- **Rich styling** — ANSI16, ANSI256, and RGB/Hex colors; bold, italic, underline
- **Built-in component library** — Table, Tabs, Modal, Input, Button, Checkbox, List, SelectPicker, Spinner, ProgressBar, Alert, Badge
- **Efficient rendering** — cell-level diffing; only changed cells are written to the terminal
- **Full Unicode support** — proper character-width handling via `go-runewidth`

---

## Installation

```bash
go get github.com/anirban/tuix
```

Requires Go 1.21+.

---

## Quick Start

```go
package main

import tuix "github.com/anirban/tuix"

func App(props tuix.Props) tuix.Element {
    count, setCount := tuix.UseState(0)

    if tuix.CurrentKey.Code == tuix.KeyEnter {
        setCount(count + 1)
    }
    if tuix.CurrentKey.Code == tuix.KeyEscape {
        tuix.Exit()
    }

    label := tuix.NewStyle().Bold(true).Foreground(tuix.Cyan)

    return tuix.Box(
        tuix.Props{Direction: tuix.Column, Gap: 1},
        tuix.NewStyle(),
        tuix.Text("Press Enter to count, Esc to quit", tuix.NewStyle()),
        tuix.Text("Count: "+fmt.Sprintf("%d", count), label),
    )
}

func main() {
    app := tuix.NewApp(80, 24)
    app.Run(App, tuix.Props{})
}
```

Run it:

```bash
go run .
```

---

## Core Concepts

### Components

A component is a Go function that accepts `tuix.Props` and returns a `tuix.Element`. Components can call hooks and compose child elements freely.

```go
func Greeting(props tuix.Props) tuix.Element {
    name := props.Values["name"].(string)
    return tuix.Text("Hello, "+name+"!", tuix.NewStyle().Bold(true))
}

// Use it inside another component:
Greeting(tuix.Props{Values: map[string]any{"name": "world"}})
```

The top-level component is passed to `app.Run`. Re-renders are triggered by state changes or keyboard events.

---

### Layout with Box

`Box` is the primary layout container. It arranges children along a main axis and supports flexbox-style sizing.

```go
tuix.Box(
    tuix.Props{
        Direction: tuix.Row,          // or tuix.Column
        Gap:       2,                 // space between children
        Padding:   [4]int{1, 2, 1, 2}, // top, right, bottom, left
        Align:     tuix.AlignCenter,  // cross-axis alignment
        Justify:   tuix.JustifySpaceBetween, // main-axis distribution
        Width:     tuix.Grow(1),      // fill available width
        Height:    tuix.Fixed(10),    // exactly 10 rows tall
    },
    tuix.NewStyle(),
    child1,
    child2,
)
```

**Sizing modes:**

| Mode | Description |
|------|-------------|
| `tuix.Fixed(n)` | Exactly `n` characters wide/tall |
| `tuix.Grow(n)` | Flex-grow with weight `n`; shares remaining space proportionally |
| `tuix.Fit()` | Sizes to content (default) |

**Alignment (cross-axis):** `AlignStart`, `AlignCenter`, `AlignEnd`, `AlignStretch`

**Justification (main-axis):** `JustifyStart`, `JustifyEnd`, `JustifyCenter`, `JustifySpaceBetween`, `JustifySpaceAround`

---

### Styling

Styles are built with a fluent API and passed as the second argument to most element constructors.

```go
style := tuix.NewStyle().
    Bold(true).
    Italic(true).
    Underline(true).
    Foreground(tuix.Green).
    Background(tuix.Hex("#1E1E2E"))
```

**Color types:**

```go
tuix.Red          // ANSI16 named color
tuix.ANSI256(214) // ANSI 256-color palette
tuix.Hex("#FF6B6B") // RGB truecolor
```

Available ANSI16 colors: `Black`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`, `White`, and their `Bright` variants.

---

### Hooks

Hooks must be called unconditionally at the top level of a component function (same rules as React).

#### UseState

```go
value, setValue := tuix.UseState(initialValue)

// Read:
fmt.Println(value)

// Write (triggers re-render):
setValue(value + 1)
```

#### UseEffect

Runs a side-effect after render. Return a cleanup function (or `nil`) to run on dependency change or unmount.

```go
tuix.UseEffect(func() func() {
    // effect runs when deps change
    ticker := time.NewTicker(time.Second)
    go func() {
        for range ticker.C {
            setTick(t => t + 1)
        }
    }()
    return func() {
        ticker.Stop() // cleanup
    }
}, []any{someDepValue})
```

---

### Keyboard Events

The current key press is available globally as `tuix.CurrentKey` during each render cycle.

```go
type Key struct {
    Code KeyCode // for special keys
    Rune rune    // for printable characters
}
```

**Special key codes:**

| Constant | Key |
|----------|-----|
| `tuix.KeyEnter` | Enter |
| `tuix.KeyBackspace` | Backspace |
| `tuix.KeyEscape` | Escape |
| `tuix.KeyTab` | Tab |
| `tuix.KeyShiftTab` | Shift+Tab |
| `tuix.KeyUp / Down / Left / Right` | Arrow keys |
| `tuix.KeySpace` | Space |
| `tuix.KeyCtrlC` | Ctrl+C |

**Handling printable characters:**

```go
if tuix.CurrentKey.Rune != 0 {
    setText(text + string(tuix.CurrentKey.Rune))
}
if tuix.CurrentKey.Code == tuix.KeyBackspace && len(text) > 0 {
    setText(text[:len(text)-1])
}
```

---

## Component Library

All components live in `components/` and are imported alongside the core package.

### Display

#### Text

```go
tuix.Text("Hello world", tuix.NewStyle().Bold(true))
```

#### Badge

```go
components.Badge("Active", tuix.NewStyle().
    Background(tuix.Green).Foreground(tuix.Black))
```

#### Alert

```go
components.Alert("Saved successfully", components.AlertSuccess)
// Variants: AlertInfo, AlertSuccess, AlertWarning, AlertError
```

#### Spinner

Automatically advances its animation frame on each render (no state needed).

```go
components.Spinner(tuix.NewStyle().Foreground(tuix.Cyan))
```

#### ProgressBar

```go
components.ProgressBar(0.65, 30, tuix.NewStyle().Foreground(tuix.Green))
// args: fraction (0.0–1.0), width, style
```

#### Panel

A bordered container with an optional title.

```go
components.Panel("Details", tuix.NewStyle(), childElement)
```

---

### Interactive

Interactive components are focus-aware. Pass `focused bool` via `props.Values` and handle keyboard events in your component to update focus state.

#### Button

```go
components.Button("Confirm", focused, tuix.NewStyle())
```

Renders highlighted when `focused` is `true`. Trigger with `tuix.KeyEnter`.

#### Input

```go
text, setText := tuix.UseState("")

// In your key handler:
if tuix.CurrentKey.Rune != 0 {
    setText(text + string(tuix.CurrentKey.Rune))
}
if tuix.CurrentKey.Code == tuix.KeyBackspace && len(text) > 0 {
    setText(text[:len(text)-1])
}

components.Input(text, focused, tuix.NewStyle())
```

#### Checkbox

```go
checked, setChecked := tuix.UseState(false)

if focused && tuix.CurrentKey.Code == tuix.KeySpace {
    setChecked(!checked)
}

components.Checkbox("Enable notifications", checked, focused, tuix.NewStyle())
```

#### List

```go
items := []string{"Apple", "Banana", "Cherry"}
cursor, setCursor := tuix.UseState(0)

if focused {
    if tuix.CurrentKey.Code == tuix.KeyDown {
        setCursor(min(cursor+1, len(items)-1))
    }
    if tuix.CurrentKey.Code == tuix.KeyUp {
        setCursor(max(cursor-1, 0))
    }
}

components.List(items, cursor, focused, tuix.NewStyle())
```

#### SelectPicker

Cycles through options with Left/Right arrows.

```go
options := []string{"Small", "Medium", "Large"}
selected, setSelected := tuix.UseState(0)

if focused {
    if tuix.CurrentKey.Code == tuix.KeyRight {
        setSelected((selected + 1) % len(options))
    }
    if tuix.CurrentKey.Code == tuix.KeyLeft {
        setSelected((selected - 1 + len(options)) % len(options))
    }
}

components.SelectPicker(options, selected, focused, tuix.NewStyle())
```

---

### Complex

#### Table

```go
headers := []string{"Name", "Status", "Email"}
rows := [][]string{
    {"Alice", "Active", "alice@example.com"},
    {"Bob",   "Away",   "bob@example.com"},
}
cursor, setCursor := tuix.UseState(0)

if focused {
    if tuix.CurrentKey.Code == tuix.KeyDown {
        setCursor(min(cursor+1, len(rows)-1))
    }
    if tuix.CurrentKey.Code == tuix.KeyUp {
        setCursor(max(cursor-1, 0))
    }
}

components.Table(headers, rows, cursor, focused, tuix.NewStyle())
```

#### Tabs

```go
tabs := []string{"All", "Active", "Away"}
active, setActive := tuix.UseState(0)

if focused {
    if tuix.CurrentKey.Code == tuix.KeyRight {
        setActive((active + 1) % len(tabs))
    }
    if tuix.CurrentKey.Code == tuix.KeyLeft {
        setActive((active - 1 + len(tabs)) % len(tabs))
    }
}

components.Tabs(tabs, active, focused, tuix.NewStyle())
```

#### Modal

Renders as an overlay. Place it last in the child list so it paints on top.

```go
open, setOpen := tuix.UseState(false)

if tuix.CurrentKey.Code == tuix.KeyEscape {
    setOpen(false)
}

if open {
    return components.Modal(
        "Contact Details",
        tuix.NewStyle(),
        detailsContent,
    )
}
```

---

## Full Example

The file `main.go` in the repository is a complete contact management app showing:

- Search input with live filtering
- Tabbed navigation (All / Active / Away)
- Scrollable table with row selection
- Modal detail overlay
- Focus cycling between panes with Tab / Shift+Tab

Use it as a reference for structuring a real application.

---

## Architecture

```
keyboard / ticker
       │
       ▼
   runtime.go          ← event loop, re-render scheduling
       │
       ▼
 component tree         ← functional components + hooks
       │
       ▼
  layout engine         ← 2-pass flexbox (measure → layout)
       │
       ▼
   renderer.go          ← element tree → screen cells
       │
       ▼
    screen.go           ← cell diffing → ANSI output → terminal
```

**Hooks cursor pattern:** State is identified by call order within a render, not by name. This is why hooks must never be called conditionally — the nth `UseState` call always corresponds to the same state slot.

**Two-pass layout:**
1. *Measure* — bottom-up: each node reports its intrinsic size
2. *Layout* — top-down: parent distributes space and assigns concrete `Rect` to each child

---

## Contributing

Contributions are welcome. Please follow these guidelines to keep the codebase consistent.

### Getting Started

```bash
git clone https://github.com/anirban/tuix
cd tuix
go mod download
go test ./...
```

### Workflow

1. **Open an issue first** for non-trivial changes to align on the approach before writing code.
2. **Branch off `main`:** `git checkout -b feat/my-feature`
3. **Keep commits focused** — one logical change per commit with a clear message.
4. **Add tests** for new layout or rendering behaviour in `*_test.go` files.
5. **Run tests and vet before opening a PR:**
   ```bash
   go test ./...
   go vet ./...
   ```
6. **Open a pull request** against `main` with a description of what changed and why.

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Keep component functions pure where possible; side effects belong in `UseEffect`
- New built-in components go in `components/` — simple/display in `components.go`, interactive in `interactive.go`, compound in `complex.go`
- Avoid adding dependencies; the stdlib + the two existing deps cover most needs

### Adding a Component

1. Write the component function in the appropriate file under `components/`
2. It should accept a `tuix.Props` parameter (use `props.Values` for component-specific options)
3. Demonstrate it in `main.go` or a separate example file
4. Document its props and keyboard contract in this README under the relevant section

### Reporting Bugs

Open a GitHub issue with:
- Go version (`go version`)
- Terminal emulator and OS
- Minimal reproduction case
- What you expected vs. what happened

---

## License

MIT — see [LICENSE](LICENSE).
