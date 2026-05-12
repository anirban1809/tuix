package tuix

type Effect struct {
	fn      func() func()
	deps    []any
	cleanup func()
	dirty   bool
}

var State []any
var StateCursor int = 0

var Effects []Effect
var EffectCursor int = 0

var pendingRender bool

func UseState[T any](initial T) (T, func(T)) {
	idx := StateCursor
	StateCursor++

	if idx >= len(State) {
		State = append(State, initial)
	}

	current := State[idx].(T)

	setter := func(next T) {
		State[idx] = next
		pendingRender = true
	}

	return current, setter
}

func UseEffect(fn func() func(), deps []any) {
	idx := EffectCursor
	EffectCursor++

	newEffect := Effect{fn: fn, deps: deps, dirty: true}

	if idx >= len(Effects) {
		Effects = append(Effects, newEffect)
	}

	for i, dep := range newEffect.deps {
		if Effects[idx].deps[i] != dep {
			Effects[idx].fn = newEffect.fn
			Effects[idx].dirty = true
			break
		}
	}
	Effects[idx].deps = newEffect.deps
}

// RunEffects runs all effects marked dirty since the last render.
// Called after screen flush so effects fire after paint, matching React semantics.
func RunEffects() {
	for i := range Effects {
		if !Effects[i].dirty {
			continue
		}
		if Effects[i].cleanup != nil {
			Effects[i].cleanup()
		}
		Effects[i].cleanup = Effects[i].fn()
		Effects[i].dirty = false
	}
}

// Context carries a value down the component tree without prop-drilling.
// Each Context owns an independent stack of values; the innermost active
// Provide call wins. The zero value of Context is not usable — construct
// one with CreateContext so the defaultValue is set.
type Context[T any] struct {
	defaultValue T
	stack        []T
}

// CreateContext returns a new Context whose UseContext readers see
// defaultValue when no enclosing Provide is active.
func CreateContext[T any](defaultValue T) *Context[T] {
	return &Context[T]{defaultValue: defaultValue}
}

// Provide pushes value onto the context's stack, runs render (during which
// any descendant calling UseContext on this context observes value), then
// pops the value back off. The pop runs via defer so a panic in render
// still unwinds the stack cleanly.
func (c *Context[T]) Provide(value T, render func() Element) Element {
	c.stack = append(c.stack, value)
	defer func() { c.stack = c.stack[:len(c.stack)-1] }()
	return render()
}

// UseContext returns the value of the innermost active Provide for c, or
// c's defaultValue if no Provide is currently on the stack.
func UseContext[T any](c *Context[T]) T {
	if len(c.stack) == 0 {
		return c.defaultValue
	}
	return c.stack[len(c.stack)-1]
}
