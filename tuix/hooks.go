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
		}
	}
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
