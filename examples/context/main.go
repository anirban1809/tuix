// Example: context
// Demonstrates: CreateContext, Provide (with the lazy thunk pattern), and
// UseContext. A single state in App drives the language; deeply nested
// components (Greeting, Footer) read it via UseContext — no prop drilling.
//
// Press Space to cycle through languages.
//
// Run: go run ./examples/context
// See: ../../DOCS.md#usecontext

package main

import "github.com/anirban1809/tuix/tuix"

type Locale struct {
	Code     string
	Greeting string
	Hint     string
}

var locales = []Locale{
	{"en", "Hello, friend.", "space to switch language"},
	{"es", "Hola, amigo.", "espacio para cambiar de idioma"},
	{"fr", "Salut, mon ami.", "espace pour changer de langue"},
	{"ja", "こんにちは、友よ。", "スペースで言語を切り替え"},
}

var LocaleContext = tuix.CreateContext(locales[0])

// Greeting consumes LocaleContext — note it takes no locale prop.
func Greeting() tuix.Element {
	l := tuix.UseContext(LocaleContext)
	return tuix.Text(l.Greeting, tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan))
}

// Footer is two levels deep from Provide, but UseContext still sees the value.
func Footer() tuix.Element {
	l := tuix.UseContext(LocaleContext)
	return tuix.Text("["+l.Code+"] "+l.Hint, tuix.NewStyle().Foreground(tuix.BrightBlack))
}

func App(props tuix.Props) tuix.Element {
	idx, setIdx := tuix.UseState(0)
	if tuix.CurrentKey.Code == tuix.KeySpace {
		setIdx((idx + 1) % len(locales))
	}

	// Provide takes a render thunk. Children evaluated *inside* the thunk
	// see the pushed value; children built outside would miss it.
	return LocaleContext.Provide(locales[idx], func() tuix.Element {
		return tuix.Box(
			tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}, Align: tuix.AlignCenter},
			tuix.NewStyle(),
			Greeting(),
			Footer(),
		)
	})
}

func main() {
	app := tuix.NewApp(60, 7)
	app.Run(App, tuix.Props{})
}
