# v0.0.16

  1. DOCS.md
  New file: the full feature-by-feature reference for tuix. Replaces the long "Core Concepts" / "Component Library" / "Full Example" sections that used to live in `README.md` and had drifted significantly from the actual code. Structured as Quick Start → Mental Model → Easy → Layout → Hooks → Keyboard → Components → Advanced → Recipes → API Reference Index. Every exported symbol is cross-linked to its source file (e.g. `tuix/hooks.go`, `tuix/components/complex.go`), every feature section ends with a link to a runnable example, and the in-doc TOC uses GitHub-flavored anchor links so jumping around is one click. Two non-obvious facts that were undocumented before now have dedicated subsections: the **two-pass render** (each event renders the tree twice — once with `CurrentKey` set, once with it zeroed — which is why unconditional `setValue(value+1)` in a component body increments by 2 per event) and the **Provide thunk rationale** (Go evaluates function arguments eagerly, so a hypothetical `Provide(value, children...)` form would evaluate children before pushing the context value; the thunk defers descendant evaluation until *after* the push).

  2. examples/ (new directory, 12 self-contained programs)
  Twelve runnable example programs, each in its own subdirectory with a single `main.go` under 100 lines. They progress easy → advanced and double as compile-time smoke tests for every signature documented in DOCS.md (if a future refactor renames `Provide` or reorders `Input`'s args, `go build ./examples/...` breaks loudly). Each file opens with a header comment naming the feature, the run command, and a deep-link into the relevant DOCS.md section.
    - `examples/hello/main.go` — minimal Box + Text program; entry point for "did installation work?"
    - `examples/counter/main.go` — `UseState` + key handling (Enter/Space/Backspace)
    - `examples/styling/main.go` — every color mode (ANSI16 swatches, ANSI256 indexes, hex truecolor) and every border preset (`BorderSharp`/`Rounded`/`Double`/`Thick`) plus bold/italic/underline composition
    - `examples/layout/main.go` — flexbox dashboard (header with `JustifySpaceBetween`, fixed sidebar + `Grow(1)` main, right-aligned footer)
    - `examples/input/main.go` — `components.Input` with paste support, live "you typed:" preview pane
    - `examples/list/main.go` — `components.List` playlist navigator (also documents the API gap that List has no `onChange`)
    - `examples/table/main.go` — `components.Table` leaderboard with `onChange` callback driving a caption below
    - `examples/tabs/main.go` — `components.Tabs` switching between three content panels via a `switch active` block
    - `examples/modal/main.go` — `components.Modal` open/close, Enter to open + Esc to close, exercising the runtime's modal overlay semantics
    - `examples/effect-clock/main.go` — `UseEffect` with empty deps spawning a 1s goroutine ticker, cleanup channel for the goroutine, demonstrating that setter closures captured at effect-creation time remain valid because they close over a stable state slot
    - `examples/context/main.go` — `CreateContext`/`Provide`/`UseContext` with a multi-locale demo (en/es/fr/ja); two consumer components (`Greeting`, `Footer`) read the same context without prop-drilling, and Space cycles locales
    - `examples/conditional/main.go` — `tuix.If` swapping between a logged-in and logged-out view; doc comment in the file calls out the eager-evaluation caveat

  3. README.md
  Rewrote as a slim entry point: hero + features + install + quick start + examples table + architecture diagram + contributing. Deleted ~360 lines of stale feature reference (the old "Core Concepts", "Component Library", and "Full Example" sections) that DOCS.md now owns. Fixed three drift bugs that had been live in the README: (a) wrong import path `github.com/anirban/tuix` → `github.com/anirban1809/tuix`, (b) reference to a nonexistent `tuix.Exit()` function in the quick-start (Ctrl-C is the actual exit path, handled by the runtime), and (c) "Full Example" section described `main.go` as a contact-management app with search/tabs/table/modal, which has been false since v0.0.13 when `main.go` was rewritten as the paste demo (and is now the theme/context demo as of v0.0.15). Also replaced the `git clone github.com/anirban/tuix` line in the contributing section with the correct path. Every feature bullet in the "Features" list now hyperlinks into the relevant DOCS.md anchor.

# v0.0.15

  1. tuix/hooks.go
  Added a React-style Context API for sharing values across a component subtree without prop-drilling. Three new exports: the generic type `Context[T any]` (holds a `defaultValue T` and an internal `stack []T`), the constructor `CreateContext[T](defaultValue T) *Context[T]`, and the reader `UseContext[T](c *Context[T]) T` which returns the top of the context's stack or its `defaultValue` when the stack is empty. The provider is a method `(*Context[T]).Provide(value T, render func() Element) Element` that appends `value` to the stack, runs `render` (during which any descendant calling `UseContext` on the same Context observes `value`), and pops via `defer` so a panic in `render` still unwinds the stack cleanly. Each Context owns its own independent stack, keyed by the Context pointer's identity rather than by a positional cursor like `UseState` — so there is no cursor reset to coordinate with `Render`'s two-pass model. Important shape note: `Provide` takes a **render thunk**, not pre-built children. Because Go evaluates function arguments eagerly, a hypothetical `Provide(value, child1, child2)` form would execute the children before the value was pushed onto the stack, and `UseContext` inside them would see the default value instead of the provided one. The thunk defers descendant evaluation until *after* the push, which is the only place during a synchronous render where the new value is visible.

# v0.0.14

  1. tuix/elements.go
  Added a new `If(condition bool, choice1, choice2 Element) Element` helper for inline conditional composition of element trees, returning `choice1` when `condition` is true and `choice2` otherwise. The helper is a plain function call rather than a control structure, so both `choice1` and `choice2` are evaluated by the caller before `If` runs — it's intended for picking between already-constructed elements, not for guarding expensive work behind a branch. Doc comment follows the godoc convention (leading with the identifier) used by the other constructors in this file (`MultilineText`, `WrappedText`) so it surfaces cleanly in generated documentation.

# v0.0.13

  1. tuix/key.go
  Added bracketed paste support to the key parser. New `KeyPaste` constant on `KeyCode`, and a new `Paste string` field on `Key` that carries the full pasted text when `Code == KeyPaste`. Introduced `KeyScanner`, a stateful parser that converts raw stdin reads into Key events — state is required because a single paste can span many `os.Stdin.Read` calls and the end marker (`\x1b[201~`) can itself straddle a read boundary. `Feed(b []byte) []Key` accumulates paste bytes into an internal `pasteBuf` while `inPaste` is true, scans for the end marker via `bytes.Index`, and emits a single `KeyPaste` event with the full content once the marker is found. Non-paste bytes still flow through the existing `ParseKey` for one- or three-byte sequences (CSI arrows, control codes).

  2. tuix/screen.go
  `Start` now emits `\033[?2004h` (enable bracketed paste mode) after raw-mode setup and cursor hide. `Stop` emits `\033[?2004l` (disable) before restoring the terminal. Without these, the terminal emulator never wraps pasted content in `\x1b[200~ … \x1b[201~` and pastes arrive as a burst of indistinguishable keystrokes.

  3. tuix/runtime.go
  Replaced the per-Read `ParseKey(buf[:n])` call with a long-lived `KeyScanner` whose `Feed` is called for every chunk. The stdin buffer was bumped from 32 to 1024 bytes — pastes can be many kilobytes, and a 32-byte buffer would slice them into ~32 chunks per kB. The escape/Ctrl-C exit check now runs inside the per-key loop emitted by the scanner so a single read containing multiple events still triggers the right shutdown.

  4. tuix/components/interactive.go
  Added paste handling to `Input`. New module-level `ansiSequence` regexp (`\x1b\[[0-9;?]*[a-zA-Z]`) strips CSI escape sequences from clipboard content so colored terminal output pasted into a field doesn't render as literal `[42m` garbage. New `lineEndings` `strings.NewReplacer` normalizes `\r\n` → `\n` and lone `\r` → `\n` (macOS clipboards deliver newlines as `\r`, Windows as `\r\n`; the renderer only understands `\n`). New `sanitizePaste(s string) string` combines those with a `strings.Map` that drops control characters below `0x20` and `0x7F` except `\n` and `\t`. The `Input` body now handles `tuix.KeyPaste` by appending `sanitizePaste(tuix.CurrentKey.Paste)` to the field value. The rendered structure switched from `MultilineText` to `WrappedText` so wrapped pastes break on the field's available width, and the outer Box gained `Width: tuix.Grow(1)` plus `Align: tuix.AlignStart` so the label sits inline at the top of the row while continuation lines indent past the label width.

  5. tuix/layout_engine.go
  Extended the layout engine so reflow callbacks fire correctly when wrapped text lives inside a Row parent, and so reflowed heights propagate up through ancestor boxes. Three coordinated changes:
    - In `layout()`, added a post-grow-distribution block that fires `child.reflow(childRects[i].Width)` for every reflow-capable child when `n.Direction == Row`. The Column path already did this inline because cross-axis width is resolved per-child; in Row direction, Grow children's widths aren't known until after the grow loop completes. `SizingFit`-height children also get `childRects[i].Height` updated to the reflowed value.
    - In `measure()`, added a clause that preserves an already-reflowed `intrinsicHeight` (`if n.reflow != nil && n.intrinsicHeight > height { height = n.intrinsicHeight }`). Without this, a second measure pass would clobber the reflowed leaf height back to 0 (since reflow leaves have no children to sum from under `SizingFit`).
    - `ComputeLayout` now runs measure→layout twice when the tree contains any reflow node (`hasReflow(root)`). The first pass establishes widths and fires reflow; the second pass re-measures with the now-correct leaf heights so ancestor boxes (padded containers, borders) grow to fit. Bounded to two passes since width allocations don't depend on heights in this engine, so reflow stabilizes after one round.

  6. tuix/renderer.go
  After `ComputeLayout` returns, `Render` now reads `layoutRoot.intrinsicHeight` (which reflects the second measure pass) instead of the pre-reflow `contentH` computed earlier. If the post-reflow height exceeds the screen's current height, the screen is resized and `ComputeLayout` runs a third time against the enlarged available rect. `EnsureRoom` is called with the post-reflow height too, so scrollback bookkeeping accounts for any growth caused by wrapped paste content.

  7. main.go
  Rewrote the demo as a bracketed-paste verification UI. The previous spacebar-appends-a-line demo is replaced with a focused `components.Input` field (rounded yellow border) plus a sibling stats block showing `paste events: N` and a `WrappedText` view of `last paste (raw):` — the side-by-side lets you compare the sanitized field value against the unfiltered clipboard content. Outer column uses `Width: tuix.Grow(1)` so the input has the full terminal width to wrap into, and `Height` stays at default `Fit` so the app sizes to its content.

# v0.0.12

  1. tuix/runtime.go
  Reworked the shutdown path in `Run` so all exit signals funnel through a single `sync.Once`-guarded `requestQuit` closure. Previously the input goroutine closed `quit` directly on stdin errors and called `close(quit)` + `a.screen.Stop()` + `os.Exit(0)` inline on Escape/Ctrl+C. The inline `os.Exit` bypassed any deferred cleanup, and concurrent shutdown paths could double-close `quit` and panic. Now the input goroutine only signals `requestQuit()`, and the main loop's `case <-quit` branch calls `a.screen.Stop()` before returning, so terminal cleanup runs deterministically on every exit path.

  2. tuix/screen.go
  `HandleResize` now emits `\033[H\033[2J\033[3J` (cursor home + clear visible + clear scrollback) before re-querying terminal dimensions via `term.GetSize`. Without this, leftover glyphs from the pre-resize viewport remained on screen until something repainted over them, producing visible artifacts when the terminal was made smaller mid-run. Clearing first guarantees the next paint draws onto a blank canvas sized to the new dimensions.

# v0.0.11

  1. tuix/node.go
  Extended `Props` with four new layout-control fields. `Align Alignment` and `Justify Justify` expose the cross-axis and main-axis alignment knobs the layout engine already supported (alignment was previously only reachable via the internal `LayoutNode` builder methods). `Width Sizing` and `Height Sizing` let callers override the previously-hardcoded `Fit()` defaults on `Box`. The zero value (`Sizing{} == Fixed(0)`) is treated as "unset" and falls back to `Fit()` inside `Box`, so existing callers are unaffected.

  2. tuix/elements.go
  `Box` now wires `props.Align` / `props.Justify` into the constructed `LayoutProps`, and respects `props.Width` / `props.Height` (with `Fit()` fallback when the field is unset). Previously these were silently ignored — the layout engine and renderer already supported them, but there was no path from the user-facing `Props` API to the internal `LayoutNode`.

  3. tuix/renderer.go
  The `Render` pass now respects the root node's own sizing intent when sizing the `available` rect. For each axis, if the root's sizing is `Fit`, `available` is trimmed down to the intrinsic content size; for `Grow` or `Fixed` roots, the full screen dimension is used. Previously the root always adopted the full screen rect regardless of its sizing — making `Fit` roots silently stretch to fill the terminal once the screen tracked real terminal dimensions.

  4. tuix/runtime.go
  `NewApp` now prefers the real terminal dimensions over the hardcoded constructor args. After `screen.Start()` populates `termCols` / `termRows` via `term.GetSize`, those values are passed to `SetDimensions`. The constructor args remain a fallback for environments where `term.GetSize` fails (e.g. when stdout is piped).

  5. tuix/screen.go
  `HandleResize` now also calls `SetDimensions(cols, rows)` so that SIGWINCH events propagate the new terminal size to the cell grid and to `s.width` / `s.height`. Previously only `termCols` / `termRows` (used by `Flush` / `EnsureRoom` for cursor bookkeeping) were updated, so the layout never reflowed when the terminal was resized mid-run.

  6. main.go
  Reworked the demo to exercise the new `Justify` / `Align` API and the renderer's sizing fix. `header` now has three labels (`◆ tuix demo`, `spacebar appends a line`, `v0.1`) with `Justify: JustifySpaceBetween` + `Align: AlignCenter`, demonstrating main-axis distribution across the full row width. Added a new `footer` Box that wraps `hint` in a `Justify: JustifyEnd` row to right-align it within the column. The outer column box now sets `Width: tuix.Grow(1)` so it fills the terminal width (giving `JustifySpaceBetween` slack to distribute), while leaving `Height` at the default `Fit` so the app doesn't stretch vertically.

# v0.0.10

  1. tuix/style.go
  Added border support on Style. New `Border` struct with per-side toggles (`Top`, `Right`, `Bottom`, `Left`), a `Chars BorderChars` glyph set, and a `Color`. New `BorderChars` struct holding the eight runes (four edges + four corners). Shipped four presets: `BorderSharp`, `BorderRounded`, `BorderDouble`, `BorderThick`. Added a `Border(b Border) Style` builder that defaults `Chars` to `BorderSharp` when unset, and a `Border.Any()` helper that reports whether any side is active. Borders are intentionally not propagated through `mergeStyles` — each element opts in explicitly.

  2. tuix/renderer.go
  - In `buildLayoutTree`, padding is now inflated by 1 cell on each side that has an active border before constructing the `LayoutNode`. This keeps `Box()` itself unaware of borders while ensuring child layout shrinks correctly to leave room for the frame.
  - In `paint` for `ElementBox`, after the background fill the new `paintBorder(screen, rect, effective, element.Style.border)` call draws the frame. Border foreground overrides the inherited foreground only when `border.Color.Type != ColorNone`, so an unset color falls through to the box's own foreground.
  - Added `paintBorder` which walks the four edges between corner cells (top, bottom, left, right) and then dispatches the four corner cells through `cornerGlyph`. Skips entirely when no side is active or the rect is empty; guards against degenerate height/width by checking `y1 != y0` / `x1 != x0` before drawing the bottom/right edges.
  - Added `cornerGlyph(cornerChar, hChar, vChar, hasH, hasV) rune` which returns the full corner glyph when both adjacent edges are active, the horizontal/vertical edge rune when only one adjacent edge is active (so partial borders render as clean line continuations), and `0` when neither is active.

  3. main.go
  Reworked the demo to exercise all three border modes: `header` now has a full rounded cyan border on all four sides; `wrapped` is a Box with only `Top` + `Bottom` set using `BorderSharp` in bright yellow (showing partial borders rendering as horizontal rules with no stray corners); `hint` is a Box with only `Left: true` using `BorderThick` in bright black (showing a single-side accent rail). Added inner padding on the partial-border boxes so text doesn't butt against the frame.

# v0.0.9

  1. tuix/node.go
  Replaced WrapWidth int with Wrap bool on Element. The element now signals that it wants to wrap, not how wide — the layout engine decides the width.

  2. tuix/elements.go
  WrappedText now takes (text, style) only — no maxWidth parameter. It sets Wrap: true on the element.

  3. tuix/layout.go
  Added a new optional field on LayoutNode: reflow func(crossSize int) int. This is the escape valve for content whose main-axis size depends on its cross-axis size. Currently invoked only when the parent's Direction is Column.

  4. tuix/layout_engine.go
  In the per-child loop inside layout(), immediately after setCrossSize(...), added a reflow call: if the parent is Column-direction and the child has a reflow callback, child.intrinsicHeight is recomputed from childRect.Width. This lets the subsequent Fit-on-height branch pick up the corrected line count so the parent's main-axis stacking is accurate.

  5. tuix/renderer.go
  - Replaced the dual-purpose multilineLines(element) with two narrow helpers: rawLines(text) splits on \n only; wrappedLines(text, maxWidth) splits on \n then wraps each segment.
  - In buildLayoutTree for ElementMultilineText, branched on element.Wrap. When Wrap is true: WidthSizing = Grow(1), HeightSizing = Fit(), and a reflow closure returns len(wrappedLines(text, width)) (falling back to 1 if width <= 0). Otherwise: same Fixed(width) × Fixed(height) derived from raw \n splits as before.
  - In paint for ElementMultilineText, branched on element.Wrap. When Wrap is true: calls wrappedLines(text, rect.Width) — the rect width is the layout engine's authoritative answer. Otherwise: calls rawLines(text).
  - Bug fix in wrapText: changed the boundary guard from `if i%maxWidth == 0 { ...; continue }` to `if i > 0 && i%maxWidth == 0 { ... }` (no continue). The original dropped one byte at every wrap boundary, including i=0, so each wrapped line lost its leading character.

  6. main.go
  Dropped the 30 argument from the WrappedText call. The wrapped paragraph now adapts to whatever width the surrounding column box settles on.

# v0.0.8

  1. tuix/node.go                                                                                                                                                                                                                                                         
  Added a WrapWidth int field to the Element struct so any element can carry an optional column-width cap that downstream code can act
   on.                                                                                                                                
                                                                                                                                      
  2. tuix/elements.go                                       
  Added a new WrappedText(text, style, maxWidth) constructor that builds an ElementMultilineText with WrapWidth set. The original     
  MultilineText constructor was left unchanged — when WrapWidth == 0, the renderer treats text exactly as before.
                                                                                                                                      
  3. tuix/renderer.go                                       
  - Added import "strings".                                                                                                           
  - Introduced multilineLines(element Element) []string as the single source of truth for "what lines does this element render?" It
  always splits on \n first, then — only if WrapWidth > 0 — passes each segment through wrapText.                                     
  - Added wrapText(text, maxWidth) []string (the function you implemented).
  - Replaced the inline \n-only logic in both buildLayoutTree (the measure pass) and paint (the draw pass) with calls to              
  multilineLines. Width is now derived by summing RuneWidth over each returned line and taking the max; height is len(lines).         
                                                                                                                                      
  4. main.go                                                                                                                                                                              
  Added a WrappedText demo: a 96-character sentence rendered with maxWidth: 30 and the same blue body style as longBlock, placed      
  between longBlock and hint in the column layout.