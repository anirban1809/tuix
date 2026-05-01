# v0.0.7

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