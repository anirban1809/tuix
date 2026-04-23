package tuix

type Effect struct {
	fn      func() func()
	deps    []any
	cleanup func()
}

var State []any
var StateCursor int = -1

func UseState[T any](initial T) (T, func(T)) {
	idx := StateCursor
	StateCursor++

	if idx >= len(State) {
		State = append(State, initial)
	}

	current := State[idx].(T)

	setter := func(next T) {
		State[idx] = next
	}

	return current, setter
}

func UseEffect(fn func() func(), deps []any) {

}
