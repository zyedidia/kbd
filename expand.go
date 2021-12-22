package kbd

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// these functions perform dollar expansion (adapted from the Go regexp implementation).

// Returns the number of arguments this template needs for expansion. Also
// returns a bool indicating whether arg 0 is used.
func nargs(template string) (n int, zero bool) {
	for len(template) > 0 {
		i := strings.Index(template, "$")
		if i < 0 {
			break
		}
		template = template[i:]
		if len(template) > 1 && template[1] == '$' {
			// Treat $$ as $.
			template = template[2:]
			continue
		}
		num, rest, ok := extract(template)
		if !ok {
			// Malformed; treat $ as raw text.
			template = template[1:]
			continue
		}
		template = rest
		if num >= 0 {
			if num == 0 {
				zero = true
			}
			n = max(n, num+1)
		}
	}
	return n, zero
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Expands all occurrences of $x, where x is a number, with the corresponding
// entry from args. Use $$ to escape a dollar sign.
func expand(template string, args []string) string {
	buf := &bytes.Buffer{}
	for len(template) > 0 {
		i := strings.Index(template, "$")
		if i < 0 {
			break
		}
		buf.WriteString(template[:i])
		template = template[i:]
		if len(template) > 1 && template[1] == '$' {
			// Treat $$ as $.
			buf.WriteByte('$')
			template = template[2:]
			continue
		}
		num, rest, ok := extract(template)
		if !ok {
			// Malformed; treat $ as raw text.
			buf.WriteByte('$')
			template = template[1:]
			continue
		}
		template = rest
		if num >= 0 {
			if num < len(args) {
				buf.WriteString(args[num])
			}
		}
	}
	buf.WriteString(template)
	return buf.String()
}

// looks for $x numbers and extracts the number.
func extract(str string) (num int, rest string, ok bool) {
	if len(str) < 2 || str[0] != '$' {
		return
	}
	str = str[1:]
	i := 0
	for i < len(str) {
		rune, size := utf8.DecodeRuneInString(str[i:])
		if !unicode.IsLetter(rune) && !unicode.IsDigit(rune) && rune != '_' {
			break
		}
		i += size
	}
	if i == 0 {
		// empty name is not okay
		return
	}
	name := str[:i]

	// Parse number.
	num = 0
	for i := 0; i < len(name); i++ {
		if name[i] < '0' || '9' < name[i] || num >= 1e8 {
			num = -1
			break
		}
		num = num*10 + int(name[i]) - '0'
	}
	// Disallow leading zeros.
	if name[0] == '0' && len(name) > 1 {
		num = -1
	}

	rest = str[i:]
	ok = num != -1
	return
}
