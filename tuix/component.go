package tuix

// DefaultProps returns a Props with sensible zero-value defaults.
// Callers can override individual fields without specifying everything.
func DefaultProps() Props {
	return Props{
		WidthSizing:  Grow(1),
		HeightSizing: Grow(1),
	}
}
