package main

import (
	"github.com/anirban1809/tuix/tuix"
)

func App(props tuix.Props) tuix.Element {
	const sample = `# Markdown Example

## Nested lists (the fix)

- Backend
  - Node.js
    - Express
  - Database
    - PostgreSQL
- Frontend
  - React

## Ordered nested list

1. Install dependencies
   1. Run ` + "`npm install`" + `
   2. Run ` + "`go mod tidy`" + `
2. Start the server

## Indent boundary (3 spaces = list, 4 spaces = paragraph)

   - three spaces: still a bullet
    - four spaces: treated as paragraph text

## Inline styles in nested items

- **Bold** parent
  - *italic* child with ` + "`inline code`" + ` ffdsa
  - ~~strikethrough~~ child

## Task list

- [ ] unfinished item
  - [x] finished sub-item
  - [ ] another sub-item
- [x] done

## Mixed ordered / unordered nesting

1. First step
   - detail A
   - detail B
2. Second step
   - detail C
`

	return tuix.Box(
		tuix.Props{
			Direction: tuix.Column,
			Gap:       1,
			Padding:   [4]int{1, 2, 1, 2},
			Width:     tuix.Grow(1),
		},
		tuix.NewStyle(),
		tuix.Markdown(sample, tuix.NewStyle()),
	)
}

func main() {
	app := tuix.NewApp(80, 25)
	app.Run(App, tuix.Props{})
}
