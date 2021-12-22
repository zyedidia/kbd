package cbind

import (
	"testing"

	"github.com/micro-editor/tcell/v2"
)

type mouseTestCase struct {
	mod     tcell.ModMask
	btn     tcell.ButtonMask
	encoded string
}

var mouseTestCases = []mouseTestCase{
	{mod: tcell.ModNone, btn: tcell.ButtonSecondary, encoded: "MouseRight"},
	{mod: tcell.ModNone, btn: tcell.WheelUp, encoded: "MouseWheelUp"},
	{mod: tcell.ModCtrl, btn: tcell.ButtonPrimary, encoded: "Ctrl+MouseLeft"},
	{mod: tcell.ModCtrl | tcell.ModAlt, btn: tcell.ButtonPrimary, encoded: "Ctrl+Alt+MouseLeft"},
}

func TestMouseEncode(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		encoded, err := Encode(c.mod, c.key, c.ch)
		if err != nil {
			t.Errorf("failed to encode key %d %d %d: %s", c.mod, c.key, c.ch, err)
		}
		if encoded != c.encoded {
			t.Errorf("failed to encode key %d %d %d: got %s, want %s", c.mod, c.key, c.ch, encoded, c.encoded)
		}
	}
}

func TestMouseDecode(t *testing.T) {
	t.Parallel()

	for _, c := range testCases {
		mod, key, ch, err := Decode(c.encoded)
		if err != nil {
			t.Errorf("failed to decode key %s: %s", c.encoded, err)
		}
		if mod != c.mod {
			t.Errorf("failed to decode key %s: invalid modifiers: got %d, want %d", c.encoded, mod, c.mod)
		}
		if key != c.key {
			t.Errorf("failed to decode key %s: invalid key: got %d, want %d", c.encoded, key, c.key)
		}
		if ch != c.ch {
			t.Errorf("failed to decode key %s: invalid rune: got %d, want %d", c.encoded, ch, c.ch)
		}
	}
}
