# tuix — Complete Reference

A Go framework for building interactive terminal UIs with a React-style mental
model. Components are plain functions; layout is flexbox; rendering is
cell-diffed.

This document is the full reference. For a one-screen pitch, see
[README.md](README.md). For runnable demos, every section links to a
self-contained program under [`examples/`](examples/) that you can run with
one command.

---

## Table of contents

- [Quick start](#quick-start)
- [Mental model](#mental-model)
- [Easy](#easy)
  - [Text](#text)
  - [Box](#box)
  - [Styling](#styling)
  - [Borders](#borders)
- [Layout](#layout)
  - [Direction, Gap, Padding](#direction-gap-padding)
  - [Sizing: Fit, Fixed, Grow](#sizing-fit-fixed-grow)
  - [Align & Justify](#align--justify)
- [Hooks](#hooks)
  - [UseState](#usestate)
  - [UseEffect](#useeffect)
  - [UseContext](#usecontext)
- [Keyboard input](#keyboard-input)
- [Component library](#component-library)
  - [Display: Badge, Alert, Spinner, ProgressBar, Panel](#display-components)
  - [Interactive: Button, Input, Checkbox, List, SelectPicker](#interactive-components)
  - [Complex: Table, Tabs, Modal](#complex-components)
- [Advanced](#advanced)
  - [Conditional rendering with `If`](#conditional-rendering-with-if)
  - [`WrappedText` vs `MultilineText`](#wrappedtext-vs-multilinetext)
  - [Context API in depth](#context-api-in-depth)
  - [Bracketed paste](#bracketed-paste)
  - [Resize handling](#resize-handling)
  - [The two-pass render](#the-two-pass-render)
- [Recipes](#recipes)
- [API reference index](#api-reference-index)

---

## Quick start

```bash
go get github.com/anirban1809/tuix
```

Requires Go 1.21+.

```go
package main

import "github.com/anirban1809/tuix/tuix"

func App(props tuix.Props) tuix.Element {
    return tuix.Box(
        tuix.Props{Padding: [4]int{1, 2, 1, 2}},
        tuix.NewStyle(),
        tuix.Text("hello, tuix", tuix.NewStyle().Bold(true).Foreground(tuix.Cyan)),
    )
}

func main() {
    app := tuix.NewApp(60, 6)
    app.Run(App, tuix.Props{})
}
```

Press **Ctrl-C** to exit (there is no `Exit()` function).

→ Runnable: [`examples/hello`](examples/hello/main.go) ·
`go run ./examples/hello`

---

## Mental model

Three ideas you need before anything else makes sense:

1. **Components are functions.** A component takes `tuix.Props` and returns a
   `tuix.Element` tree. There is no class, no lifecycle object — just a
   function that gets called on every render.
2. **The call tree is the component tree.** When `App` calls `Header()` and
   `Footer()`, those calls happen *during* `App`'s execution, so any hooks
   they call (and any context they read) are scoped to the current render.
3. **Hooks are positional.** `UseState` and `UseEffect` identify their state
   slot by *call order within a render*, not by name. Never call them inside
   `if`/`for` — the slot index would shift between renders and you'd silently
   read another component's state.

The runtime's job:

```
keyboard / ticker / resize event
        │
        ▼
   re-render the whole tree
        │
        ▼
   measure → layout (2-pass flexbox)
        │
        ▼
   paint into a cell grid
        │
        ▼
   diff against previous frame
        │
        ▼
   write only changed cells to the terminal
```

Source files referenced throughout: [`tuix/runtime.go`](tuix/runtime.go),
[`tuix/renderer.go`](tuix/renderer.go),
[`tuix/layout_engine.go`](tuix/layout_engine.go),
[`tuix/hooks.go`](tuix/hooks.go).

---

## Easy

### Text

```go
tuix.Text("hello", tuix.NewStyle().Bold(true))
```

`Text` renders a single line. Newlines in the string are **not** treated as
line breaks — use [`MultilineText`](#wrappedtext-vs-multilinetext) for that.

### Box

`Box` is the only container. It lays out children using a flexbox-like
algorithm.

```go
tuix.Box(
    tuix.Props{
        Direction: tuix.Column,         // Row or Column (default: Row)
        Gap:       1,                   // empty cells between children
        Padding:   [4]int{1, 2, 1, 2},  // top, right, bottom, left
        Align:     tuix.AlignCenter,    // cross-axis alignment
        Justify:   tuix.JustifyStart,   // main-axis distribution
        Width:     tuix.Grow(1),        // optional; default is Fit()
        Height:    tuix.Fit(),          // optional; default is Fit()
    },
    tuix.NewStyle(),  // background, foreground, border (no padding)
    childA,
    childB,
)
```

→ Runnable: [`examples/layout`](examples/layout/main.go) ·
`go run ./examples/layout`

### Styling

`Style` is immutable and chainable.

```go
s := tuix.NewStyle().
    Bold(true).
    Italic(true).
    Underline(true).
    Foreground(tuix.Hex("#ff6b6b")).
    Background(tuix.ANSI256(236))
```

Colors come in three flavours:

| Constructor                  | Range                           | Example                |
| ---------------------------- | ------------------------------- | ---------------------- |
| `tuix.Red`, `tuix.Cyan`, …   | ANSI 16 (named)                 | `tuix.BrightMagenta`   |
| `tuix.ANSI256(n)`            | 256-color palette (0–255)       | `tuix.ANSI256(214)`    |
| `tuix.Hex("#rrggbb")`        | 24-bit truecolor                | `tuix.Hex("#ffd93d")`  |

Named colors: `Black`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`,
`White`, and their `Bright*` variants — see
[`tuix/style.go`](tuix/style.go).

**Inheritance.** Styles flow from parent to child: a child whose foreground is
`ColorNone` inherits the parent's foreground. The same applies to background.
Bold is *promoted* (a bold parent makes children bold) but italic and
underline are not inherited — they are per-element opt-ins.

→ Runnable: [`examples/styling`](examples/styling/main.go) ·
`go run ./examples/styling`

### Borders

Borders are part of `Style`, not `Props`:

```go
tuix.NewStyle().Border(tuix.Border{
    Top: true, Right: true, Bottom: true, Left: true,
    Chars: tuix.BorderRounded,  // or BorderSharp, BorderDouble, BorderThick
    Color: tuix.Cyan,
})
```

You can toggle individual sides — `Left: true` alone draws a single-side
accent rail.

When a border is active, the layout engine automatically inflates the box's
padding by 1 cell on the bordered side so children don't get clipped. You do
not need to manually account for the border in your own padding.

---

## Layout

### Direction, Gap, Padding

`Direction` is the main axis. `Row` stacks children left-to-right;
`Column` stacks them top-to-bottom.

```go
tuix.Box(
    tuix.Props{Direction: tuix.Row, Gap: 2, Padding: [4]int{0, 1, 0, 1}},
    tuix.NewStyle(),
    childA, childB, childC,
)
```

`Padding` order is `{top, right, bottom, left}` — same as CSS shorthand.

### Sizing: Fit, Fixed, Grow

Each axis (`Width` and `Height`) can use one of three modes:

| Mode             | Behaviour                                                |
| ---------------- | -------------------------------------------------------- |
| `tuix.Fit()`     | Hugs the content (default for `Box`).                    |
| `tuix.Fixed(n)`  | Exactly `n` cells.                                       |
| `tuix.Grow(n)`   | Flex-grow with weight `n`; shares leftover space.        |

`Grow(1)` on the root makes your app fill the terminal width. Multiple
siblings with `Grow` divide leftover space by weight: `Grow(2)` + `Grow(1)`
splits 2:1.

### Align & Justify

`Justify` distributes children along the **main axis**:

| Value                       | Effect                                                  |
| --------------------------- | ------------------------------------------------------- |
| `tuix.JustifyStart`         | Pack at start (default)                                 |
| `tuix.JustifyEnd`           | Pack at end                                             |
| `tuix.JustifyCenter`        | Center as a group                                       |
| `tuix.JustifySpaceBetween`  | First/last hug edges; equal gaps between siblings       |
| `tuix.JustifySpaceAround`   | Equal gaps including half-gaps at the edges             |

`Align` controls the **cross axis** for each child:

| Value                  | Effect                                          |
| ---------------------- | ----------------------------------------------- |
| `tuix.AlignStretch`    | Children stretch to fill cross axis (default)   |
| `tuix.AlignStart`      | Pack at start                                   |
| `tuix.AlignCenter`     | Centered                                        |
| `tuix.AlignEnd`        | Pack at end                                     |

→ Runnable: [`examples/layout`](examples/layout/main.go)

---

## Hooks

Hooks live in [`tuix/hooks.go`](tuix/hooks.go) and follow the same rules as
React: call them at the top of a component, in the same order, every render.

### UseState

```go
value, setValue := tuix.UseState(0)
setValue(value + 1)  // schedules a re-render
```

The initial value is used only on the first call to that slot. The setter is a
closure over the slot index, so it's safe to capture in goroutines and
`UseEffect` callbacks — it always writes to the same slot.

⚠ **Don't read state and write it back unconditionally in a component body.**
The runtime re-renders the component tree **twice per event** (see
[The two-pass render](#the-two-pass-render)); a bare `setValue(value+1)` in
the body increments by 2, not 1. Always gate on a condition (`if key == ...`).

→ Runnable: [`examples/counter`](examples/counter/main.go)

### UseEffect

```go
tuix.UseEffect(func() func() {
    ticker := time.NewTicker(time.Second)
    go func() {
        for range ticker.C {
            setNow(time.Now())
        }
    }()
    return func() { ticker.Stop() }  // cleanup
}, []any{someDep})
```

- The effect runs after the render commits.
- If any element in `deps` differs from the previous render, the previous
  cleanup runs and the effect re-runs.
- Return `nil` if you don't need cleanup.
- An empty `deps` (`[]any{}`) runs the effect exactly once, on mount.

Caveat: state written from a goroutine doesn't trigger a render directly —
the next event (key, internal 500ms tick, or resize) will pick it up.

→ Runnable: [`examples/effect-clock`](examples/effect-clock/main.go)

### UseContext

Share a value across a subtree without prop-drilling.

```go
type Theme struct { Fg, Bg tuix.Color }

var ThemeContext = tuix.CreateContext(Theme{Fg: tuix.White, Bg: tuix.Black})

func Header() tuix.Element {
    t := tuix.UseContext(ThemeContext)
    return tuix.Text("◆ hi", tuix.NewStyle().Foreground(t.Fg))
}

func App(props tuix.Props) tuix.Element {
    return ThemeContext.Provide(Theme{Fg: tuix.BrightCyan, Bg: tuix.Black}, func() tuix.Element {
        return Header()
    })
}
```

**Crucial gotcha:** `Provide` takes a **render thunk** (`func() Element`),
not pre-built children. Children must be created *inside* the thunk so they
run while the value is on the context stack. Children built outside the thunk
have already executed and `UseContext` inside them sees the default value.

See [Context API in depth](#context-api-in-depth) for the why.

→ Runnable: [`examples/context`](examples/context/main.go)

---

## Keyboard input

The current keyboard event lives in the global `tuix.CurrentKey` during each
render. It's zeroed in the second render pass so keys don't get
double-handled — see [The two-pass render](#the-two-pass-render).

```go
type Key struct {
    Code  KeyCode  // a KeyXxx constant, or KeyNone for plain runes
    Rune  rune     // the printable character, or 0
    Paste string   // when Code == KeyPaste, the full pasted text
}
```

| Constant            | Trigger                                |
| ------------------- | -------------------------------------- |
| `tuix.KeyEnter`     | Enter / Return                         |
| `tuix.KeyBackspace` | Backspace (`0x7F` or `0x08`)           |
| `tuix.KeyEscape`    | Esc                                    |
| `tuix.KeyTab`       | Tab                                    |
| `tuix.KeyShiftTab`  | Shift-Tab (CSI `Z`)                    |
| `tuix.KeyUp/Down/Left/Right` | Arrow keys                    |
| `tuix.KeySpace`     | Spacebar (also surfaces as `Rune=' '`) |
| `tuix.KeyCtrlC`     | Quits the app (handled by runtime)     |
| `tuix.KeyPaste`     | Bracketed paste; `Paste` holds the text|

Plain printable input:

```go
if tuix.CurrentKey.Rune != 0 {
    setText(text + string(tuix.CurrentKey.Rune))
}
```

→ Runnable: [`examples/input`](examples/input/main.go)

---

## Component library

All built-ins live in [`tuix/components/`](tuix/components/) and import as
`github.com/anirban1809/tuix/tuix/components`.

### Display components

Source: [`tuix/components/components.go`](tuix/components/components.go).

```go
components.Badge("Active", tuix.Black, tuix.Green)
//          Badge(label, fg, bg)

components.Alert(components.AlertSuccess, "Saved!")
//          Alert(kind, message)
//   Kinds: AlertInfo, AlertSuccess, AlertWarning, AlertError

components.Spinner("loading...")
//          Spinner(label) — animates one frame per render

components.ProgressBar(0.65, 30, tuix.Green)
//          ProgressBar(value 0..1, width, fillColor)

components.Panel("Details", 40, child1, child2)
//          Panel(title, width, children...)
```

### Interactive components

Source: [`tuix/components/interactive.go`](tuix/components/interactive.go).

```go
components.Button("Confirm", focused)
//          Button(label, focused)
//   Render highlighted when focused=true. You handle Enter yourself.

components.Input("name>", "▌", focused, value, setValue)
//          Input(label, cursor, focused, value, onChange)
//   Handles typing / backspace / space / paste internally.

components.Checkbox("Notifications", focused, onChange)
//          Checkbox(label, focused, onChange func(bool))
//   Owns its checked state; calls onChange every render with current value.

components.List(items, focused)
//          List(items, focused)
//   Up/Down moves highlight. Selection state is internal and not surfaced.

components.SelectPicker(options, focused)
//          SelectPicker(options, focused)
//   Left/Right cycles options. Selection state is internal and not surfaced.
```

→ Runnables: [`examples/input`](examples/input/main.go),
[`examples/list`](examples/list/main.go)

### Complex components

Source: [`tuix/components/complex.go`](tuix/components/complex.go).

```go
components.Table(headers, rows, focused, onChange)
//          Table([]string, [][]string, bool, func(int))
//   Up/Down moves row selection; onChange called every render with index.

components.Tabs(tabs, focused, onChange)
//          Tabs([]string, bool, func(int))
//   Left/Right switches active tab.

components.Modal("Confirm?", visible, 36, onClose, child1, child2)
//          Modal(title, visible, width, onClose func(), children...)
//   Esc calls onClose. Place it last in its parent so it paints on top.
```

→ Runnables: [`examples/table`](examples/table/main.go),
[`examples/tabs`](examples/tabs/main.go),
[`examples/modal`](examples/modal/main.go)

---

## Advanced

### Conditional rendering with `If`

```go
tuix.If(loggedIn, dashboard, loginPrompt)
```

`If` returns one of two pre-built elements. Because it's a regular function
call, **both branches are evaluated** before `If` runs. Use it for cheap
elements you've already built; don't try to guard expensive work behind one
branch.

Source: [`tuix/elements.go`](tuix/elements.go).

→ Runnable: [`examples/conditional`](examples/conditional/main.go)

### `WrappedText` vs `MultilineText`

| Constructor                      | Splits on `\n`? | Word-wraps?           |
| -------------------------------- | --------------- | --------------------- |
| `tuix.Text(s, style)`            | no (single line)| no                    |
| `tuix.MultilineText(s, style)`   | yes             | no                    |
| `tuix.WrappedText(s, style)`     | yes             | yes (to parent width) |

`WrappedText` sets `Width: Grow(1)` internally, so it expands to fill its
parent's cross-axis space and breaks lines to fit. It registers a `reflow`
callback with the layout engine, which is why
[`tuix/layout_engine.go`](tuix/layout_engine.go) runs a second measure pass
when a wrapped element is present in the tree.

### Context API in depth

The Context API is in [`tuix/hooks.go`](tuix/hooks.go). Three exports:

- `tuix.CreateContext[T](defaultValue T) *Context[T]` — construct
- `(*Context[T]).Provide(value T, render func() Element) Element` — scope
- `tuix.UseContext[T](*Context[T]) T` — read

Under the hood, each `Context` owns its own `[]T` stack. `Provide` appends a
value, runs the render thunk, and pops via `defer`. `UseContext` returns the
top of the stack, or the context's `defaultValue` if the stack is empty.

**Why a thunk?** Children in tuix are eager Go function arguments. If
`Provide` took children directly — `Provide(value, child1, child2)` — Go
would evaluate `child1`/`child2` *before* `Provide` ran, so any
`UseContext` inside them would see the empty stack. The thunk defers
descendant evaluation until *after* the push.

**Stack identity vs cursor identity.** Unlike `UseState` (slot-by-call-order),
context is keyed by `*Context[T]` pointer. There's no cursor to reset between
renders, and stacks survive across renders because they live on the Context
object, not in a global slab.

→ Runnable: [`examples/context`](examples/context/main.go)

### Bracketed paste

The runtime enables bracketed paste mode on startup (terminal emits
`\x1b[200~`…`\x1b[201~` around clipboard content). A
[`KeyScanner`](tuix/key.go) reassembles paste fragments across multiple
`stdin.Read` calls and delivers them as a single `Key{Code: KeyPaste, Paste:
"…"}` event. The built-in `Input` component sanitises pasted content (strips
CSI escapes, normalises CRLF, drops control chars except `\n`/`\t`) before
inserting it. If you build a custom text field, look at the `sanitizePaste`
helper in [`tuix/components/interactive.go`](tuix/components/interactive.go).

### Resize handling

The runtime listens for `SIGWINCH` and re-queries the terminal size via
`golang.org/x/term`, then re-renders. Before the re-query, the screen is
cleared (`\033[H\033[2J\033[3J`) so leftover glyphs from a smaller-resize
don't linger. This is handled automatically — your code doesn't need to do
anything special.

### The two-pass render

Each event triggers two render passes of your component tree:

1. **Pass 1** — `tuix.CurrentKey` is set; state setters mutate state.
2. **Pass 2** — `tuix.CurrentKey` is zeroed; the tree renders with updated
   state. Only pass 2's tree is painted.

This is why unconditional `setValue(value+1)` in a component body double-
increments per event. Two practical rules:

- **Gate setters on a condition** — usually a key check (`if Code ==
  KeyEnter { ... }`). Pass 1 fires the handler; pass 2's condition is false
  because `CurrentKey` was zeroed.
- **Side effects belong in `UseEffect`**, never in the component body —
  otherwise they fire twice.

Source: [`tuix/runtime.go`](tuix/runtime.go) `App.Render`.

---

## Recipes

Beyond the per-feature examples, here are short patterns for common needs.

### Focus cycling

Track which interactive element is focused with a single `UseState`:

```go
focus, setFocus := tuix.UseState(0)
if tuix.CurrentKey.Code == tuix.KeyTab {
    setFocus((focus + 1) % 3)
}

return tuix.Box(
    tuix.Props{Direction: tuix.Column, Gap: 1},
    tuix.NewStyle(),
    components.Input("name>", "▌", focus == 0, name, setName),
    components.Input("email>", "▌", focus == 1, email, setEmail),
    components.Button("Submit", focus == 2),
)
```

### Polling external data

```go
data, setData := tuix.UseState[*Result](nil)
tuix.UseEffect(func() func() {
    done := make(chan struct{})
    go func() {
        t := time.NewTicker(5 * time.Second)
        defer t.Stop()
        for {
            select {
            case <-done: return
            case <-t.C:
                if r, err := fetch(); err == nil { setData(r) }
            }
        }
    }()
    return func() { close(done) }
}, []any{})
```

### Toast notifications

A short-lived banner that auto-dismisses:

```go
toast, setToast := tuix.UseState("")
tuix.UseEffect(func() func() {
    if toast == "" { return nil }
    timer := time.AfterFunc(3*time.Second, func() { setToast("") })
    return func() { timer.Stop() }
}, []any{toast})
```

---

## API reference index

Quick jump-to-source for everything exported.

### Core types

- [`Props`](tuix/node.go), [`Element`](tuix/node.go), [`LayoutProps`](tuix/node.go)
- [`Box`](tuix/elements.go), [`Text`](tuix/elements.go),
  [`MultilineText`](tuix/elements.go), [`WrappedText`](tuix/elements.go),
  [`If`](tuix/elements.go)

### Layout primitives

- [`Direction`](tuix/layout.go) — `Row`, `Column`
- [`Sizing`](tuix/layout.go) — `Fit()`, `Fixed(n)`, `Grow(n)`
- [`Alignment`](tuix/layout.go) — `AlignStretch`, `AlignStart`,
  `AlignCenter`, `AlignEnd`
- [`Justify`](tuix/layout.go) — `JustifyStart`, `JustifyEnd`,
  `JustifyCenter`, `JustifySpaceBetween`, `JustifySpaceAround`

### Style

- [`Style`](tuix/style.go) — `NewStyle()`, `.Bold()`, `.Italic()`,
  `.Underline()`, `.Foreground()`, `.Background()`, `.Border()`
- [`Color`](tuix/style.go) — `Black`…`BrightWhite`, `Hex(s)`, `ANSI256(n)`
- [`Border`](tuix/style.go) — sides + `Chars` + `Color`
- Presets: `BorderSharp`, `BorderRounded`, `BorderDouble`, `BorderThick`

### Hooks

- [`UseState[T](initial T) (T, func(T))`](tuix/hooks.go)
- [`UseEffect(fn func() func(), deps []any)`](tuix/hooks.go)
- [`CreateContext[T](defaultValue T) *Context[T]`](tuix/hooks.go)
- [`UseContext[T](c *Context[T]) T`](tuix/hooks.go)
- [`(*Context[T]).Provide(v T, render func() Element) Element`](tuix/hooks.go)

### Keyboard

- [`Key`](tuix/key.go), [`KeyCode`](tuix/key.go), [`CurrentKey`](tuix/key.go)
- [`ParseKey`](tuix/key.go), [`KeyScanner`](tuix/key.go)

### Runtime

- [`NewApp(width, height int) *App`](tuix/runtime.go)
- [`(*App).Run(fn func(Props) Element, props Props)`](tuix/runtime.go)

### Components

- Display: [`Badge`](tuix/components/components.go),
  [`Alert`](tuix/components/components.go),
  [`Spinner`](tuix/components/components.go),
  [`ProgressBar`](tuix/components/components.go),
  [`Panel`](tuix/components/components.go)
- Interactive: [`Button`](tuix/components/interactive.go),
  [`Input`](tuix/components/interactive.go),
  [`Checkbox`](tuix/components/interactive.go),
  [`List`](tuix/components/interactive.go),
  [`SelectPicker`](tuix/components/interactive.go)
- Complex: [`Table`](tuix/components/complex.go),
  [`Tabs`](tuix/components/complex.go),
  [`Modal`](tuix/components/complex.go)

### Examples directory

| Example                                              | Demonstrates                       |
| ---------------------------------------------------- | ---------------------------------- |
| [`hello`](examples/hello/main.go)                    | minimal program                    |
| [`counter`](examples/counter/main.go)                | `UseState` + keyboard              |
| [`styling`](examples/styling/main.go)                | colors, borders, text styles       |
| [`layout`](examples/layout/main.go)                  | flexbox: direction/sizing/justify  |
| [`input`](examples/input/main.go)                    | `Input` component + paste          |
| [`list`](examples/list/main.go)                      | navigable `List`                   |
| [`table`](examples/table/main.go)                    | `Table` with `onChange`            |
| [`tabs`](examples/tabs/main.go)                      | `Tabs` switching content panels    |
| [`modal`](examples/modal/main.go)                    | `Modal` open/close                 |
| [`effect-clock`](examples/effect-clock/main.go)      | `UseEffect` + goroutine cleanup    |
| [`context`](examples/context/main.go)                | `Context` + `Provide` thunk        |
| [`conditional`](examples/conditional/main.go)        | `If` helper                        |

Run any of them:

```bash
go run ./examples/<name>
```
