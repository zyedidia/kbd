package syntax

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/zyedidia/gpeg/charset"
	"github.com/zyedidia/gpeg/memo"
	"github.com/zyedidia/gpeg/pattern"
	"github.com/zyedidia/gpeg/vm"
	"github.com/zyedidia/kbd"
)

var parser vm.Code

func init() {
	prog := pattern.MustCompile(pattern.Grammar("Pattern", grammar))
	parser = vm.Encode(prog)
}

func compile(name string, root *memo.Capture, s string) kbd.Pattern {
	var p kbd.Pattern
	fmt.Println(root.Id(), idExpression)
	switch root.Id() {
	case idPattern:
		p = compile(name, root.Child(0), s)
	// case idGrammar:
	// 	nonterms := make(map[string]kbd.Pattern)
	// 	it := root.ChildIterator(0)
	// 	for c := it(); c != nil; c = it() {
	// 		k, v := compileDef(name, c, s)
	// 		nonterms[name+k] = v
	// 	}
	// 	p = kbd.Grammar(name+"token", nonterms)
	case idExpression:
		alternations := make([]kbd.Pattern, 0, root.NumChildren())
		it := root.ChildIterator(0)
		for c := it(); c != nil; c = it() {
			alternations = append(alternations, compile(name, c, s))
		}
		p = kbd.Alt(alternations...)
	case idSequence:
		concats := make([]kbd.Pattern, 0, root.NumChildren())
		it := root.ChildIterator(0)
		for c := it(); c != nil; c = it() {
			concats = append(concats, compile(name, c, s))
		}
		p = kbd.Seq(concats...)
	case idSuffix:
		if root.NumChildren() == 2 {
			c := root.Child(1)
			switch c.Id() {
			case idQUESTION:
				p = kbd.Opt(compile(name, root.Child(0), s))
			case idSTAR:
				p = kbd.Star(compile(name, root.Child(0), s))
			case idPLUS:
				p = kbd.Plus(compile(name, root.Child(0), s))
			}
		} else {
			p = compile(name, root.Child(0), s)
		}
	case idPrimary:
		switch root.Child(0).Id() {
		case idBRACEO:
			cpatt := compile(name, root.Child(1), s)
			group := literal(root.Child(2), s)
			p = kbd.Cap(cpatt, group)
		case idIdentifier, idLiteral, idClass:
			p = compile(name, root.Child(0), s)
		case idOPEN:
			p = compile(name, root.Child(1), s)
		case idDOT:
			p = kbd.AnyRune()
		}
	case idLiteral:
		p = kbd.MustLit(literal(root, s))
	case idClass:
		var set charset.Set
		if root.NumChildren() <= 0 {
			break
		}
		complement := false
		if root.Child(0).Id() == idCARAT {
			complement = true
		}
		it := root.ChildIterator(0)
		i := 0
		for c := it(); c != nil; c = it() {
			if i == 0 && complement {
				i++
				continue
			}
			set = set.Add(compileSet(c, s))
		}
		if complement {
			set = set.Complement()
		}
		p = kbd.Set(set)
		// case idIdentifier:
		// 	p = pattern.NonTerm(name + parseId(root, s))
	}
	return p
}

var special = map[byte]byte{
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'\'': '\'',
	'"':  '"',
	'[':  '[',
	']':  ']',
	'\\': '\\',
	'-':  '-',
}

func parseChar(char string) byte {
	switch char[0] {
	case '\\':
		for k, v := range special {
			if char[1] == k {
				return v
			}
		}

		i, _ := strconv.ParseInt(string(char[1:]), 8, 8)
		return byte(i)
	default:
		return char[0]
	}
}

func parseId(root *memo.Capture, s string) string {
	ident := &bytes.Buffer{}
	it := root.ChildIterator(0)
	for c := it(); c != nil; c = it() {
		ident.WriteString(s[c.Start():c.End()])
	}
	return ident.String()
}

func literal(root *memo.Capture, s string) string {
	lit := &bytes.Buffer{}
	it := root.ChildIterator(0)
	for c := it(); c != nil; c = it() {
		lit.WriteByte(parseChar(s[c.Start():c.End()]))
	}
	return lit.String()
}

func compileDef(name string, root *memo.Capture, s string) (string, kbd.Pattern) {
	id := root.Child(0)
	exp := root.Child(1)
	return parseId(id, s), compile(name, exp, s)
}

func compileSet(root *memo.Capture, s string) charset.Set {
	switch root.NumChildren() {
	case 1:
		c := root.Child(0)
		return charset.New([]byte{parseChar(s[c.Start():c.End()])})
	case 2:
		c1, c2 := root.Child(0), root.Child(1)
		return charset.Range(parseChar(s[c1.Start():c1.End()]), parseChar(s[c2.Start():c2.End()]))
	}
	return charset.Set{}
}

func Compile(name, s string) (kbd.Pattern, error) {
	match, n, ast, errs := parser.Exec(strings.NewReader(s), memo.NoneTable{})
	if len(errs) != 0 {
		return nil, errs[0]
	}
	if !match {
		return nil, fmt.Errorf("Invalid PEG: failed at %d", n)
	}

	return compile(name, ast.Child(0), s), nil
}

func MustCompile(name, s string) kbd.Pattern {
	p, err := Compile(name, s)
	if err != nil {
		panic(err)
	}
	return p
}
