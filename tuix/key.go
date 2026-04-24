package tuix

// KeyCode identifies special (non-printable) keys.
type KeyCode int

var Keys = make(chan Key, 4)

// CurrentKey holds the key being processed in the current render pass.
var CurrentKey Key

const (
	KeyNone KeyCode = iota
	KeyEnter
	KeyBackspace
	KeyEscape
	KeyTab
	KeyShiftTab
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeySpace
	KeyCtrlC
)

// Key represents a single keyboard event. Either Code or Rune is set.
type Key struct {
	Code KeyCode
	Rune rune
}

// ParseKey converts raw terminal bytes into a Key.
// Arrow keys arrive as 3-byte escape sequences; all other specials are 1 byte.
func ParseKey(b []byte) Key {
	if len(b) == 0 {
		return Key{}
	}
	if len(b) >= 3 && b[0] == 0x1B && b[1] == '[' {
		switch b[2] {
		case 'A':
			return Key{Code: KeyUp}
		case 'B':
			return Key{Code: KeyDown}
		case 'C':
			return Key{Code: KeyRight}
		case 'D':
			return Key{Code: KeyLeft}
		case 'Z':
			return Key{Code: KeyShiftTab}
		}
	}
	switch b[0] {
	case 0x1B:
		return Key{Code: KeyEscape}
	case 0x0D, 0x0A:
		return Key{Code: KeyEnter}
	case 0x7F, 0x08:
		return Key{Code: KeyBackspace}
	case 0x09:
		return Key{Code: KeyTab}
	case 0x20:
		return Key{Code: KeySpace}
	case 0x03:
		return Key{Code: KeyCtrlC}
	}
	return Key{Rune: rune(b[0])}
}
