package kbd

import (
	"fmt"
	"strings"

	"github.com/zyedidia/kbd/cbind"

	"github.com/micro-editor/tcell/v2"
)

// An Event is a blueprint for an actual tcell event and given a tcell event
// specifies if it matches.
type Event interface {
	Match(ev tcell.Event) bool
	String() string
}

// A KeyEvent matches a specific key event made from the given combination of
// key and modifiers.
type KeyEvent struct {
	ch  rune
	key tcell.Key
	mod tcell.ModMask
}

func (ke *KeyEvent) Match(ev tcell.Event) bool {
	if kev, ok := ev.(*tcell.EventKey); ok {
		if kev.Key() == tcell.KeyRune && ke.key == tcell.KeyRune {
			return kev.Rune() == ke.ch && kev.Modifiers() == ke.mod
		}
		return kev.Key() == ke.key && kev.Modifiers() == ke.mod
	}
	return false
}

func (ke *KeyEvent) String() string {
	s, err := cbind.Encode(ke.mod, ke.key, ke.ch)
	if err != nil {
		return err.Error()
	}
	return s
}

// A MouseEvent matches a combination of mouse button and modifiers. It also
// stores the resulting position of the mouse event after the event has been
// matched.
type MouseEvent struct {
	btn tcell.ButtonMask
	mod tcell.ModMask

	// result
	X, Y int
}

func (me *MouseEvent) Match(ev tcell.Event) bool {
	if mev, ok := ev.(*tcell.EventMouse); ok {
		if mev.Buttons() == me.btn && mev.Modifiers() == me.mod {
			me.X, me.Y = mev.Position()
			return true
		}
	}
	return false
}

func (me *MouseEvent) String() string {
	s, err := cbind.EncodeMouse(me.mod, me.btn)
	if err != nil {
		return err.Error()
	}
	return s
}

// A WildcardRuneEvent matches any rune event in the given range and stores the
// matching event.
type WildcardRuneEvent struct {
	Low, High rune
	// result
	Rune rune
}

func (we *WildcardRuneEvent) Match(ev tcell.Event) bool {
	if kev, ok := ev.(*tcell.EventKey); ok && kev.Key() == tcell.KeyRune {
		r := kev.Rune()
		if r >= we.Low && r <= we.High {
			we.Rune = kev.Rune()
			return true
		}
	}
	return false
}

func (we *WildcardRuneEvent) String() string {
	return fmt.Sprintf("Any [%s-%s]", string(we.Low), string(we.High))
}

// A PasteEvent matches any tcell paste event and stores the pasted text.
type PasteEvent struct {
	// result
	Text string
}

func (pe *PasteEvent) Match(ev tcell.Event) bool {
	if pev, ok := ev.(*tcell.EventPaste); ok {
		pe.Text = pev.Text()
		return true
	}
	return false
}

func (pe *PasteEvent) String() string {
	return "Paste"
}

// A ResizeEvent matches any resize event and stores the new screen dimensions.
type ResizeEvent struct {
	// result
	W, H int
}

func (re *ResizeEvent) Match(ev tcell.Event) bool {
	if rev, ok := ev.(*tcell.EventResize); ok {
		re.W, re.H = rev.Size()
		return true
	}
	return false
}

func (re *ResizeEvent) String() string {
	return "Resize"
}

type RawEvent struct {
	esc string
}

// ToEvent constructs a single event from a string.
func ToEvent(s string) (Event, error) {
	switch strings.ToLower(s) {
	case "paste":
		return &PasteEvent{}, nil
	case "resize":
		return &ResizeEvent{}, nil
	case "any":
		return &WildcardRuneEvent{}, nil
	default:
		mod, btn, err := cbind.DecodeMouse(s)
		if err == nil {
			return &MouseEvent{
				mod: mod,
				btn: btn,
			}, nil
		}
		mod, key, ch, err := cbind.Decode(s)
		if err != nil {
			return nil, err
		}
		return &KeyEvent{
			ch:  ch,
			key: key,
			mod: mod,
		}, nil
	}
}
