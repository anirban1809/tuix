package tuix

import (
	"os"
	"time"
)

type App struct {
	screen   *Screen
	renderer *ComponentRenderer
}

func NewApp(width, height int) *App {

	screen := NewScreenWriter(width, height, os.Stdout)
	screen.Start()
	screen.SetDimensions(width, height)

	renderer := NewRenderer(screen)

	return &App{
		screen:   screen,
		renderer: renderer,
	}
}

var ticker = make(chan bool, 1)
var CurrentTick bool = false

func (a *App) Run(fn func() Element) {
	a.Render(fn)

	quit := make(chan struct{})

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
			a.Render(fn)
		case tick := <-ticker:
			CurrentTick = tick
			a.Render(fn)
		}
	}
}

func (a *App) Render(fn func() Element) {
	// Pass 1: process key events and mutate state
	StateCursor = 0
	EffectCursor = 0
	fn()

	// Pass 2: render with updated state; key is now consumed
	CurrentKey = Key{}
	StateCursor = 0
	EffectCursor = 0
	next := fn()

	a.renderer.Render(next)
	a.screen.Flush()
	RunEffects()
}
