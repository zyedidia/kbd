package kbd

import (
	"strconv"

	"github.com/zyedidia/generic/stack"
)

type machine struct {
	pc int // program counter
	sp int // subject pointer

	cmds []string
	vars []interface{}
	caps *stack.Stack[int]

	status status
}

type status struct {
	blocked bool // blocked waiting for another event
	failed  bool // did not match
	done    bool // finished matching (success if !failed)
}

func newMachine() *machine {
	return &machine{
		pc:   0,
		sp:   0,
		cmds: nil,
		vars: nil,
		caps: stack.New[int](),
	}
}

// copy the machine but start it at a new pc
func (m *machine) cpy(pc int) *machine {
	cmds := make([]string, len(m.cmds))
	vars := make([]interface{}, len(m.vars))
	copy(cmds, m.cmds)
	copy(vars, m.vars)
	return &machine{
		pc:     pc,
		sp:     m.sp,
		cmds:   cmds,
		vars:   vars,
		caps:   m.caps.Copy(),
		status: m.status,
	}
}

func (m *machine) done(success bool) {
	m.status.done = true
	m.status.failed = !success
}

func (m *machine) step(prog Program, evs events) (splitpc int, ok bool) {
	if m.pc < 0 || m.pc >= len(prog) {
		m.done(true)
		return
	}

	insn := prog[m.pc]

	switch t := insn.(type) {
	case iEnd:
		m.done(true)
		return
	case iConsume:
		if m.sp >= len(evs) {
			m.status.blocked = true
			return
		}
		if !t.match.Match(evs[m.sp]) {
			m.done(false)
			return
		}
		m.sp++
		m.pc++
	case iJump:
		m.pc += t.lbl
	case iSplit:
		newpc := m.pc + t.lbl2
		m.pc += t.lbl1
		return newpc, true
	case iCapStart:
		m.caps.Push(m.sp)
		m.pc++
	case iCapEnd:
		// this is confusing so it is heavily commented
		last := m.caps.Pop()
		n, zero := nargs(t.cmd)
		n-- // we only care about the number of args after arg 0
		// the zero arg corresponds to a capture of all the events
		var arg0 interface{}
		if zero && m.sp-last == 1 {
			// one event, encode it directly (TODO)
			arg0 = evs.slice(last, m.sp)
		} else if zero {
			// multiple events, get the concatenated slice as a string (meant
			// to be used only for rune events).
			arg0 = evs.slice(last, m.sp)
		}
		var arg0name string
		if zero {
			// add arg0 to the var list and get its name
			arg0name = m.mkvar(arg0)
		}
		// how many additional arguments do we need?
		var args []string
		if n > 0 {
			// if there are additional arguments fetch them out of m.cmds and
			// remove them from m.cmds
			ln := len(m.cmds)
			// last n values in m.cmds are the args
			args = m.cmds[ln-n : ln]
			// remove the last n vals
			m.cmds = m.cmds[:ln-n]
		}
		// args are arg0 plus any additional args
		args = append([]string{arg0name}, args...)
		// perform expansion
		result := expand(t.cmd, args)
		// add the result directly into m.cmds
		m.cmds = append(m.cmds, result)

		m.pc++
	}
	return
}

func (m *machine) mkvar(val interface{}) string {
	m.vars = append(m.vars, val)
	return "$" + strconv.Itoa(len(m.vars)-1)
}
