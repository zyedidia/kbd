package kbd

import (
	"bytes"
	"fmt"
)

// A Program is a sequence of key parsing instructions.
type Program []insn

func (p Program) String() string {
	s := &bytes.Buffer{}
	for _, insn := range p {
		s.WriteString(insn.String())
		s.WriteByte('\n')
	}
	return s.String()
}

type insn interface {
	String() string
}

type iCapStart struct{}

func (i iCapStart) String() string {
	return "cap start"
}

type iCapEnd struct {
	cmd string
}

func (i iCapEnd) String() string {
	return fmt.Sprintf("cap end '%v'", i.cmd)
}

type iConsume struct {
	match Event
}

func (i iConsume) String() string {
	return fmt.Sprintf("consume %v", i.match)
}

type iEnd struct{}

func (i iEnd) String() string {
	return "end"
}

type iJump struct {
	lbl int
}

func (i iJump) String() string {
	return fmt.Sprintf("jump %v", i.lbl)
}

type iSplit struct {
	lbl1, lbl2 int
}

func (i iSplit) String() string {
	return fmt.Sprintf("split %v, %v", i.lbl1, i.lbl2)
}
