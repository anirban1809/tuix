package tuix

type App struct {
	screen   *Screen
	renderer *ComponentRenderer
}

func NewApp(width, height int) *App {
	StdOutScreen.SetDimensions(width, height)
	return &App{
		screen:   StdOutScreen,
		renderer: Renderer,
	}
}

func (a *App) Run(fn func() Element) {
	// Pass 1: process key events and mutate state
	StateCursor = 0
	fn()

	// Pass 2: render with updated state; key is now consumed
	CurrentKey = Key{}
	StateCursor = 0
	next := fn()

	Renderer.Render(next)
	StdOutScreen.Flush()
}
