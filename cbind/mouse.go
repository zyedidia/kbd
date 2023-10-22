package cbind

import (
	"errors"
	"fmt"
	"strings"

	"github.com/micro-editor/tcell/v2"
)

var nameBtn = map[string]tcell.ButtonMask{
	"mouseleft":       tcell.ButtonPrimary,
	"mouseright":      tcell.ButtonSecondary,
	"mousemiddle":     tcell.ButtonMiddle,
	"mousethumbnext":  tcell.Button4,
	"mousethumbprev":  tcell.Button5,
	"mousebutton6":    tcell.Button6,
	"mousebutton7":    tcell.Button7,
	"mousebutton8":    tcell.Button8,
	"mousewheelup":    tcell.WheelUp,
	"MouseWheelDown":  tcell.WheelDown,
	"mousewheelleft":  tcell.WheelLeft,
	"mousewheelright": tcell.WheelRight,
	"mousenone":       tcell.ButtonNone,
}

var btnName = map[tcell.ButtonMask]string{
	tcell.ButtonPrimary:   "MouseLeft",
	tcell.ButtonSecondary: "MouseRight",
	tcell.ButtonMiddle:    "MouseMiddle",
	tcell.Button4:         "MouseThumbNext",
	tcell.Button5:         "MouseThumbPrev",
	tcell.Button6:         "MouseButton6",
	tcell.Button7:         "MouseButton7",
	tcell.Button8:         "MouseButton8",
	tcell.WheelUp:         "MouseWheelUp",
	tcell.WheelDown:       "MouseWheelDown",
	tcell.WheelLeft:       "MouseWheelLeft",
	tcell.WheelRight:      "MouseWheelRight",
	tcell.ButtonNone:      "MouseNone",
}

func DecodeMouse(s string) (mod tcell.ModMask, btn tcell.ButtonMask, err error) {
	if len(s) == 0 {
		return 0, 0, fmt.Errorf("%s: %w", s, ErrInvalidKeyEvent)
	}

	split := strings.Split(s, "+")
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

		if btn, ok := nameBtn[pieceLower]; ok {
			return mod, btn, nil
		}
	}
	return mod, btn, errors.New("not enough parts")
}

func EncodeMouse(mod tcell.ModMask, btn tcell.ButtonMask) (string, error) {
	var b strings.Builder
	var wrote bool

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

	if name, ok := btnName[btn]; ok {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(name)
	} else {
		return "", ErrInvalidKeyEvent
	}
	return b.String(), nil
}
