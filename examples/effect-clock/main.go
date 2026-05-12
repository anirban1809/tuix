// Example: effect-clock
// Demonstrates: UseEffect with a goroutine + cleanup, driving a live clock.
//
// The effect runs once on mount (empty deps), spawns a 1-second ticker, and
// pushes the current time into state. Cleanup closes a `done` channel so the
// goroutine exits if the effect re-runs. Setter closures are stable across
// renders (they close over a fixed state slot) so calling them from a
// goroutine is safe.
//
// Note: state updates from outside the runtime event loop don't trigger a
// render until the next event arrives. The runtime emits an internal tick
// every 500ms, so the clock visibly updates within that window.
//
// Run: go run ./examples/effect-clock
// See: ../../DOCS.md#useeffect

package main

import (
	"time"

	"github.com/anirban1809/tuix/tuix"
)

func App(props tuix.Props) tuix.Element {
	now, setNow := tuix.UseState(time.Now())

	tuix.UseEffect(func() func() {
		done := make(chan struct{})
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					setNow(t)
				}
			}
		}()
		return func() { close(done) }
	}, []any{})

	bigTime := tuix.NewStyle().Bold(true).Foreground(tuix.BrightCyan)
	dim := tuix.NewStyle().Foreground(tuix.BrightBlack)

	return tuix.Box(
		tuix.Props{Direction: tuix.Column, Gap: 1, Padding: [4]int{1, 2, 1, 2}, Align: tuix.AlignCenter},
		tuix.NewStyle(),
		tuix.Text(now.Format("15:04:05"), bigTime),
		tuix.Text(now.Format("Mon, 02 Jan 2006"), dim),
		tuix.Text("ctrl-c to quit", dim),
	)
}

func main() {
	app := tuix.NewApp(60, 7)
	app.Run(App, tuix.Props{})
}
