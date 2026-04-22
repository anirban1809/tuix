# Hooks System

This is a React-style hooks system for a TUI framework. The core trick is the **cursor pattern**: hooks don't have names or IDs — they're identified purely by *the order they're called* during a render. This is why React's "rules of hooks" (no hooks in conditionals/loops) exist.

## The Mental Model

**`HookState` lives on each component `Node`.** When a component renders, it reads and writes to its own `slots []any` using `cursor` as an index — each hook call (`useState`, `useEffect`) advances the cursor by one slot.

```
First render:          slots = []           cursor = 0
  useState(0)    →     slots = [0]          cursor = 1  (slot created)
  useState("")   →     slots = [0, ""]      cursor = 2  (slot created)

Re-render:             slots = [42, "hi"]   cursor = 0  (reset)
  useState(0)    →     returns slots[0]=42  cursor = 1  (slot read)
  useState("")   →     returns slots[1]="hi" cursor = 2 (slot read)
```

## The Three Pieces to Implement

### 1. `useState[T](initial T) (T, func(T))`
- On first call (slot doesn't exist): write `initial` to `slots[cursor]`
- On re-render: read existing value from `slots[cursor]`
- Returns the value + a setter that triggers a re-render

### 2. `useEffect(fn func() func(), deps []any)`
- Store `fn` and `deps` in `slots[cursor]`
- On re-render: compare old deps vs new deps — only re-run if changed
- The `effects []Effect` field on `HookState` is a queue that gets flushed *after* the render completes
- `fn` can return a cleanup function, called before the next effect or on unmount

### 3. Wiring into the Reconciler
- Each component `Node` needs to own a `*HookState`
- Before calling `Render(props)`, reset `cursor` to 0
- After render, flush the `effects` queue
- When a setter fires, mark the node dirty and trigger a re-render loop

## The Flow

```
Input event / setState called
        ↓
  Reset cursor to 0
        ↓
  Call component's Render fn
    → useState reads slots[0], slots[1]...
    → useEffect queues into effects[]
        ↓
  Reconciler diffs the new element tree
        ↓
  Flush effects[] (run new/changed ones)
        ↓
  Screen.Flush() → paint to terminal
```

## The `Effect` Type

```go
type Effect struct {
    fn      func() func()   // the effect body + cleanup returner
    deps    []any           // dependency values
    cleanup func()          // result of previous fn() call
}
```

## Key Design Note

`HookState` needs to **survive re-renders** (must live on the `Node`, not be recreated each render), but be **discarded on unmount** (handled by the `deletions` list returned from `reconcile.go`'s `Reconcile` function).
