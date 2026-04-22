package tuix

// package level variable, represents the current node being rendered
var currentNode *Node

type HookState struct {
	slots          []any
	cursor         int
	effects        []Effect
	scheduleUpdate func()
}

type Effect struct {
	fn      func() func()
	deps    []any
	cleanup func()
}

var State []any
var StateCursor int = 0

func UseState[T any](initial T) (T, func(T)) {
	State = append(State, initial)

	setter := func(newVal T) {
		State[StateCursor] = newVal
	}

	value := initial
	return value, setter
}

func UseEffect(fn func() func(), deps []any) {

}
