package tuix

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	screen   *Screen
	renderer *ComponentRenderer
}

func NewApp(width, height int) *App {

	screen := NewScreenWriter(width, height, os.Stdout)
	screen.Start()

	// Prefer the real terminal dimensions over the constructor args so
	// layout fills the actual viewport. The args remain a fallback for
	// environments where term.GetSize fails (e.g. piped output).
	if screen.termCols > 0 && screen.termRows > 0 {
		screen.SetDimensions(screen.termCols, screen.termRows)
	} else {
		screen.SetDimensions(width, height)
	}

	renderer := NewRenderer(screen)

	return &App{
		screen:   screen,
		renderer: renderer,
	}
}

var ticker = make(chan bool, 1)
var CurrentTick bool = false

func (a *App) Run(fn func(props Props) Element, props Props) {
	a.Render(fn, props)

	quit := make(chan struct{})

	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)

	go func() {
		tick := false
		for {
			time.Sleep(time.Millisecond * 500)
			tick = !tick
			ticker <- tick
		}
	}()

	go func() {
		buf := make([]byte, 32)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				close(quit)
				return
			}
			key := ParseKey(buf[:n])
			if key.Code == KeyEscape || key.Code == KeyCtrlC {
				close(quit)
				a.screen.Stop()
				os.Exit(0)
				return
			}
			Keys <- key
		}
	}()

	for {
		select {
		case <-quit:
			return
		case key := <-Keys:
			CurrentKey = key
			a.Render(fn, props)
		case tick := <-ticker:
			CurrentTick = tick
			a.Render(fn, props)
		case <-resize:
			a.screen.HandleResize()
			a.Render(fn, props)
		}
	}
}

func (a *App) Render(fn func(props Props) Element, props Props) {
	// Pass 1: process key events and mutate state
	StateCursor = 0
	EffectCursor = 0
	fn(props)

	// Pass 2: render with updated state; key is now consumed
	CurrentKey = Key{}
	StateCursor = 0
	EffectCursor = 0
	next := fn(props)

	a.renderer.Render(next)
	a.screen.Flush()

	pendingRender = false
	RunEffects()

	if pendingRender {
		CurrentKey = Key{}
		StateCursor = 0
		EffectCursor = 0
		next := fn(props)

		a.renderer.Render(next)
		a.screen.Flush()
	}
}
