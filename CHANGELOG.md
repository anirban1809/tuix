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