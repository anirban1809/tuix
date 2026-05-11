package tuix

import "bytes"

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
	KeyPaste
)

// Key represents a single keyboard event. Either Code or Rune is set.
// When Code == KeyPaste, Paste holds the full pasted text.
type Key struct {
	Code  KeyCode
	Rune  rune
	Paste string
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

// Bracketed paste mode wraps pasted content in these markers so applications
// can distinguish a paste from a fast burst of typed keystrokes. Terminals
// emit them only after we send the enable sequence \x1b[?2004h on startup.
var (
	pasteStart = []byte{0x1B, '[', '2', '0', '0', '~'}
	pasteEnd   = []byte{0x1B, '[', '2', '0', '1', '~'}
)

// KeyScanner is a stateful parser that converts raw stdin reads into Key
// events. State is required because a single paste can span many Read calls
// and the end marker may itself straddle a Read boundary.
type KeyScanner struct {
	inPaste  bool
	pasteBuf []byte
}

// Feed consumes one chunk of stdin bytes and returns any complete Key events
// produced. Unfinished sequences (mid-paste content, partial end marker) are
// retained inside the scanner for the next call.
func (s *KeyScanner) Feed(b []byte) []Key {
	var keys []Key
	for len(b) > 0 {
		if s.inPaste {
			s.pasteBuf = append(s.pasteBuf, b...)
			b = nil
			idx := bytes.Index(s.pasteBuf, pasteEnd)
			if idx < 0 {
				return keys
			}
			pasted := string(s.pasteBuf[:idx])
			rest := append([]byte(nil), s.pasteBuf[idx+len(pasteEnd):]...)
			s.pasteBuf = nil
			s.inPaste = false
			keys = append(keys, Key{Code: KeyPaste, Paste: pasted})
			b = rest
			continue
		}
		if bytes.HasPrefix(b, pasteStart) {
			b = b[len(pasteStart):]
			s.inPaste = true
			continue
		}
		consumed := 1
		if b[0] == 0x1B && len(b) >= 3 && b[1] == '[' {
			consumed = 3
		}
		keys = append(keys, ParseKey(b[:consumed]))
		b = b[consumed:]
	}
	return keys
}
