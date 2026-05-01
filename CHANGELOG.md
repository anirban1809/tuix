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