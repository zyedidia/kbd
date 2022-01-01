package kbd

import (
	"math"

	"github.com/zyedidia/gpeg/charset"
)

type Pattern interface {
	Compile() Program
}

type AltNode struct {
	s1 Pattern
	s2 Pattern
}

func Alt(s ...Pattern) Pattern {
	acc := s[len(s)-1]
	for i := len(s) - 2; i >= 0; i-- {
		acc = &AltNode{
			s1: s[i],
			s2: acc,
		}
	}

	return acc
}

func (n *AltNode) Compile() Program {
	p1 := n.s1.Compile()
	p2 := n.s2.Compile()
	var prog Program
	prog = append(prog, iSplit{1, len(p1) + 2})
	prog = append(prog, p1...)
	prog = append(prog, iJump{len(p2) + 1})
	prog = append(prog, p2...)
	return prog
}

type SeqNode struct {
	s1 Pattern
	s2 Pattern
}

func Seq(s ...Pattern) Pattern {
	acc := s[0]
	for _, p := range s[1:] {
		acc = &SeqNode{
			s1: acc,
			s2: p,
		}
	}

	return acc
}

func (n *SeqNode) Compile() Program {
	p1 := n.s1.Compile()
	p2 := n.s2.Compile()
	var prog Program
	prog = append(prog, p1...)
	prog = append(prog, p2...)
	return prog
}

type StarNode struct {
	s Pattern
}

func Star(s Pattern) *StarNode {
	return &StarNode{
		s: s,
	}
}

func (n *StarNode) Compile() Program {
	p := n.s.Compile()
	var prog Program
	prog = append(prog, iSplit{1, len(p) + 2})
	prog = append(prog, p...)
	prog = append(prog, iJump{-len(p) - 1})
	return prog
}

type PlusNode struct {
	s Pattern
}

func Plus(s Pattern) *PlusNode {
	return &PlusNode{
		s: s,
	}
}

func (n *PlusNode) Compile() Program {
	p := n.s.Compile()
	var prog Program
	prog = append(prog, p...)
	prog = append(prog, iSplit{-len(p), 1})
	return prog
}

type OptNode struct {
	s Pattern
}

func Opt(s Pattern) *OptNode {
	return &OptNode{
		s: s,
	}
}

func (n *OptNode) Compile() Program {
	p := n.s.Compile()
	var prog Program
	prog = append(prog, iSplit{1, len(p) + 1})
	prog = append(prog, p...)
	return prog
}

type LitNode struct {
	ev Event
}

func Lit(ev Event) *LitNode {
	return &LitNode{
		ev: ev,
	}
}

func MustLit(s string) *LitNode {
	ev, err := ToEvent(s)
	if err != nil {
		panic(err)
	}
	return &LitNode{
		ev: ev,
	}
}

func AnyRune() *LitNode {
	return &LitNode{
		ev: &WildcardRuneEvent{
			Low:  0,
			High: math.MaxInt32,
		},
	}
}

func RangeRune(low, high rune) *LitNode {
	return &LitNode{
		ev: &WildcardRuneEvent{
			Low:  low,
			High: high,
		},
	}
}

func Set(s charset.Set) *LitNode {
	return &LitNode{
		ev: &WildcardRuneSetEvent{
			Set: s,
		},
	}
}

func (n *LitNode) Compile() Program {
	return Program{
		iConsume{
			match: n.ev,
		},
	}
}

type CapNode struct {
	s   Pattern
	cmd string
}

func Cap(n Pattern, cmd string) *CapNode {
	return &CapNode{
		s:   n,
		cmd: cmd,
	}
}

func (c *CapNode) Compile() Program {
	p := c.s.Compile()
	var prog Program
	prog = append(prog, iCapStart{})
	prog = append(prog, p...)
	prog = append(prog, iCapEnd{c.cmd})
	return prog
}

type EndNode struct{}

func End() *EndNode {
	return &EndNode{}
}

func (n *EndNode) Compile() Program {
	return Program{
		iEnd{},
	}
}

type NonTermNode struct {
	name string
}

func NonTerm(name string) Pattern {
	return &NonTermNode{
		name: name,
	}
}

func (n *NonTermNode) Compile() Program {
	var prog Program
	prog = append(prog, iOpenCall{n.name})
	return prog
}

func Grammar(fns map[string]Program, root string) Program {
	var prog Program

	fnlocs := make(map[string]int)

	i := 0
	for name, fn := range fns {
		fnlocs[name] = i
		prog = append(prog, fn...)
		prog = append(prog, iRet{})
		i += 1 + len(fn)
	}

	for j, insn := range prog {
		switch t := insn.(type) {
		case iOpenCall:
			prog[j] = iCall{fnlocs[t.name]}
		}
	}

	return prog
}
