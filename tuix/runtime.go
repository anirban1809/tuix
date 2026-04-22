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

func (a *App) Run(fn func(Props) []Element) {
	next := Component(fn, Props{})

	Renderer.Render(next)
	StdOutScreen.Flush()

}
