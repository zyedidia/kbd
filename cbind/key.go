package cbind

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/micro-editor/tcell/v2"
)

// Modifier labels
const (
	LabelCtrl  = "ctrl"
	LabelAlt   = "alt"
	LabelMeta  = "meta"
	LabelShift = "shift"
)

// ErrInvalidKeyEvent is the error returned when encoding or decoding a key event fails.
var ErrInvalidKeyEvent = errors.New("invalid key event")

// UnifyEnterKeys is a flag that determines whether or not KPEnter (keypad
// enter) key events are interpreted as Enter key events. When enabled, Ctrl+J
// key events are also interpreted as Enter key events.
var UnifyEnterKeys = false

var fullKeyNames = map[string]string{
	"backspace2": "Backspace",
	"pgup":       "PageUp",
	"pgdn":       "PageDown",
	"esc":        "Escape",
}

var ctrlKeys = map[rune]tcell.Key{
	' ':  tcell.KeyCtrlSpace,
	'a':  tcell.KeyCtrlA,
	'b':  tcell.KeyCtrlB,
	'c':  tcell.KeyCtrlC,
	'd':  tcell.KeyCtrlD,
	'e':  tcell.KeyCtrlE,
	'f':  tcell.KeyCtrlF,
	'g':  tcell.KeyCtrlG,
	'h':  tcell.KeyCtrlH,
	'i':  tcell.KeyCtrlI,
	'j':  tcell.KeyCtrlJ,
	'k':  tcell.KeyCtrlK,
	'l':  tcell.KeyCtrlL,
	'm':  tcell.KeyCtrlM,
	'n':  tcell.KeyCtrlN,
	'o':  tcell.KeyCtrlO,
	'p':  tcell.KeyCtrlP,
	'q':  tcell.KeyCtrlQ,
	'r':  tcell.KeyCtrlR,
	's':  tcell.KeyCtrlS,
	't':  tcell.KeyCtrlT,
	'u':  tcell.KeyCtrlU,
	'v':  tcell.KeyCtrlV,
	'w':  tcell.KeyCtrlW,
	'x':  tcell.KeyCtrlX,
	'y':  tcell.KeyCtrlY,
	'z':  tcell.KeyCtrlZ,
	'\\': tcell.KeyCtrlBackslash,
	']':  tcell.KeyCtrlRightSq,
	'^':  tcell.KeyCtrlCarat,
	'_':  tcell.KeyCtrlUnderscore,
}

// KeyNamesUniform is a normalized key name map, which is the same as KeyNames
// but all keys are lowercase and '-' is replaced with '+'. This is an
// optimization for cbind.
var keyNamesUniform = map[string]tcell.Key{
	"enter":      tcell.KeyEnter,
	"backspace":  tcell.KeyBackspace,
	"tab":        tcell.KeyTab,
	"backtab":    tcell.KeyBacktab,
	"esc":        tcell.KeyEsc,
	"backspace2": tcell.KeyBackspace2,
	"delete":     tcell.KeyDelete,
	"insert":     tcell.KeyInsert,
	"up":         tcell.KeyUp,
	"down":       tcell.KeyDown,
	"left":       tcell.KeyLeft,
	"right":      tcell.KeyRight,
	"home":       tcell.KeyHome,
	"end":        tcell.KeyEnd,
	"upleft":     tcell.KeyUpLeft,
	"upright":    tcell.KeyUpRight,
	"downleft":   tcell.KeyDownLeft,
	"downright":  tcell.KeyDownRight,
	"center":     tcell.KeyCenter,
	"pgdn":       tcell.KeyPgDn,
	"pgup":       tcell.KeyPgUp,
	"clear":      tcell.KeyClear,
	"exit":       tcell.KeyExit,
	"cancel":     tcell.KeyCancel,
	"pause":      tcell.KeyPause,
	"print":      tcell.KeyPrint,
	"f1":         tcell.KeyF1,
	"f2":         tcell.KeyF2,
	"f3":         tcell.KeyF3,
	"f4":         tcell.KeyF4,
	"f5":         tcell.KeyF5,
	"f6":         tcell.KeyF6,
	"f7":         tcell.KeyF7,
	"f8":         tcell.KeyF8,
	"f9":         tcell.KeyF9,
	"f10":        tcell.KeyF10,
	"f11":        tcell.KeyF11,
	"f12":        tcell.KeyF12,
	"f13":        tcell.KeyF13,
	"f14":        tcell.KeyF14,
	"f15":        tcell.KeyF15,
	"f16":        tcell.KeyF16,
	"f17":        tcell.KeyF17,
	"f18":        tcell.KeyF18,
	"f19":        tcell.KeyF19,
	"f20":        tcell.KeyF20,
	"f21":        tcell.KeyF21,
	"f22":        tcell.KeyF22,
	"f23":        tcell.KeyF23,
	"f24":        tcell.KeyF24,
	"f25":        tcell.KeyF25,
	"f26":        tcell.KeyF26,
	"f27":        tcell.KeyF27,
	"f28":        tcell.KeyF28,
	"f29":        tcell.KeyF29,
	"f30":        tcell.KeyF30,
	"f31":        tcell.KeyF31,
	"f32":        tcell.KeyF32,
	"f33":        tcell.KeyF33,
	"f34":        tcell.KeyF34,
	"f35":        tcell.KeyF35,
	"f36":        tcell.KeyF36,
	"f37":        tcell.KeyF37,
	"f38":        tcell.KeyF38,
	"f39":        tcell.KeyF39,
	"f40":        tcell.KeyF40,
	"f41":        tcell.KeyF41,
	"f42":        tcell.KeyF42,
	"f43":        tcell.KeyF43,
	"f44":        tcell.KeyF44,
	"f45":        tcell.KeyF45,
	"f46":        tcell.KeyF46,
	"f47":        tcell.KeyF47,
	"f48":        tcell.KeyF48,
	"f49":        tcell.KeyF49,
	"f50":        tcell.KeyF50,
	"f51":        tcell.KeyF51,
	"f52":        tcell.KeyF52,
	"f53":        tcell.KeyF53,
	"f54":        tcell.KeyF54,
	"f55":        tcell.KeyF55,
	"f56":        tcell.KeyF56,
	"f57":        tcell.KeyF57,
	"f58":        tcell.KeyF58,
	"f59":        tcell.KeyF59,
	"f60":        tcell.KeyF60,
	"f61":        tcell.KeyF61,
	"f62":        tcell.KeyF62,
	"f63":        tcell.KeyF63,
	"f64":        tcell.KeyF64,
	"ctrl+a":     tcell.KeyCtrlA,
	"ctrl+b":     tcell.KeyCtrlB,
	"ctrl+c":     tcell.KeyCtrlC,
	"ctrl+d":     tcell.KeyCtrlD,
	"ctrl+e":     tcell.KeyCtrlE,
	"ctrl+f":     tcell.KeyCtrlF,
	"ctrl+g":     tcell.KeyCtrlG,
	"ctrl+j":     tcell.KeyCtrlJ,
	"ctrl+k":     tcell.KeyCtrlK,
	"ctrl+l":     tcell.KeyCtrlL,
	"ctrl+n":     tcell.KeyCtrlN,
	"ctrl+o":     tcell.KeyCtrlO,
	"ctrl+p":     tcell.KeyCtrlP,
	"ctrl+q":     tcell.KeyCtrlQ,
	"ctrl+r":     tcell.KeyCtrlR,
	"ctrl+s":     tcell.KeyCtrlS,
	"ctrl+t":     tcell.KeyCtrlT,
	"ctrl+u":     tcell.KeyCtrlU,
	"ctrl+v":     tcell.KeyCtrlV,
	"ctrl+w":     tcell.KeyCtrlW,
	"ctrl+x":     tcell.KeyCtrlX,
	"ctrl+y":     tcell.KeyCtrlY,
	"ctrl+z":     tcell.KeyCtrlZ,
	"ctrl+space": tcell.KeyCtrlSpace,
	"ctrl+_":     tcell.KeyCtrlUnderscore,
	"ctrl+]":     tcell.KeyCtrlRightSq,
	"ctrl+\\":    tcell.KeyCtrlBackslash,
	"ctrl+^":     tcell.KeyCtrlCarat,
}

// Decode decodes a string as a key or combination of keys.
func Decode(s string) (mod tcell.ModMask, key tcell.Key, ch rune, err error) {
	if len(s) == 0 {
		return 0, 0, 0, fmt.Errorf("%s: %w", s, ErrInvalidKeyEvent)
	}

	// Special case for plus rune decoding
	if s[len(s)-1:] == "+" {
		key = tcell.KeyRune
		ch = '+'

		if len(s) == 1 {
			return mod, key, ch, nil
		} else if len(s) == 2 {
			return 0, 0, 0, fmt.Errorf("%s: %w", s, ErrInvalidKeyEvent)
		} else {
			s = s[:len(s)-2]
		}
	}

	split := strings.Split(s, "+")
DECODEPIECE:
	for _, piece := range split {
		// Decode modifiers
		pieceLower := strings.ToLower(piece)
		switch pieceLower {
		case LabelCtrl:
			mod |= tcell.ModCtrl
			continue
		case LabelAlt:
			mod |= tcell.ModAlt
			continue
		case LabelMeta:
			mod |= tcell.ModMeta
			continue
		case LabelShift:
			mod |= tcell.ModShift
			continue
		}

		// Decode key
		for shortKey, fullKey := range fullKeyNames {
			if pieceLower == strings.ToLower(fullKey) {
				pieceLower = shortKey
				break
			}
		}
		switch pieceLower {
		case "backspace":
			key = tcell.KeyBackspace2
			continue
		case "space", "spacebar":
			key = tcell.KeyRune
			ch = ' '
			continue
		}
		if k, ok := keyNamesUniform[pieceLower]; ok {
			key = k
			if key < 0x80 {
				ch = rune(k)
			}
			continue DECODEPIECE
		}

		// Decode rune
		if len(piece) > 1 {
			return 0, 0, 0, fmt.Errorf("%s: %w", s, ErrInvalidKeyEvent)
		}

		key = tcell.KeyRune
		ch = rune(piece[0])
	}

	if mod&tcell.ModCtrl != 0 {
		k, ok := ctrlKeys[unicode.ToLower(ch)]
		if ok {
			key = k
			if UnifyEnterKeys && key == ctrlKeys['j'] {
				key = tcell.KeyEnter
			} else if key < 0x80 {
				ch = rune(key)
			}
		}
	}

	return mod, key, ch, nil
}

// Encode encodes a key or combination of keys a string.
func Encode(mod tcell.ModMask, key tcell.Key, ch rune) (string, error) {
	var b strings.Builder
	var wrote bool

	if mod&tcell.ModCtrl != 0 {
		if key == tcell.KeyBackspace || key == tcell.KeyTab || key == tcell.KeyEnter {
			mod ^= tcell.ModCtrl
		} else {
			for _, ctrlKey := range ctrlKeys {
				if key == ctrlKey {
					mod ^= tcell.ModCtrl
					break
				}
			}
		}
	}

	if key != tcell.KeyRune {
		if UnifyEnterKeys && key == ctrlKeys['j'] {
			key = tcell.KeyEnter
		} else if key < 0x80 {
			ch = rune(key)
		}
	}

	// Encode modifiers
	if mod&tcell.ModCtrl != 0 {
		b.WriteString(upperFirst(LabelCtrl))
		wrote = true
	}
	if mod&tcell.ModAlt != 0 {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(upperFirst(LabelAlt))
		wrote = true
	}
	if mod&tcell.ModMeta != 0 {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(upperFirst(LabelMeta))
		wrote = true
	}
	if mod&tcell.ModShift != 0 {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(upperFirst(LabelShift))
		wrote = true
	}

	if key == tcell.KeyRune && ch == ' ' {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString("Space")
	} else if key != tcell.KeyRune {
		// Encode key
		keyName := tcell.KeyNames[key]
		if keyName == "" {
			return "", ErrInvalidKeyEvent
		}
		keyName = strings.ReplaceAll(keyName, "-", "+")
		fullKeyName := fullKeyNames[strings.ToLower(keyName)]
		if fullKeyName != "" {
			keyName = fullKeyName
		}

		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(keyName)
	} else {
		// Encode rune
		if wrote {
			b.WriteRune('+')
		}
		b.WriteRune(ch)
	}

	return b.String(), nil
}

func upperFirst(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
