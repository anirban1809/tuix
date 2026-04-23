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

var ElementsTable map[string]bool

func (a *App) Run(fn func() Element) {
	StateCursor = 0
	next := fn()

	Renderer.Render(next)
	StdOutScreen.Flush()
}
