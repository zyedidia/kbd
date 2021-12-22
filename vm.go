package kbd

import (
	"bytes"

	"github.com/micro-editor/tcell/v2"
)

type events []tcell.Event

// when multiple events are concatenated together, a string representation is
// used to unify.
func ev2str(ev tcell.Event) string {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyRune {
			return string(ev.Rune())
		}
	}
	return ""
}

func (evs events) slice(start, end int) string {
	if start >= end {
		return ""
	}

	buf := &bytes.Buffer{}
	slc := evs[start:end]
	for _, e := range slc {
		buf.WriteString(ev2str(e))
	}
	return buf.String()
}

type VM struct {
	// program to be executed
	prog Program
	// separate machines executing each program path
	machines []*machine

	// all events seen so far
	evs events
}

func NewVM(prog Program) *VM {
	return &VM{
		prog:     prog,
		machines: []*machine{newMachine()},
	}
}

func (vm *VM) Reset() {
	vm.machines = []*machine{newMachine()}
	vm.evs = nil
}

// vm is blocked if all submachines are blocked
func (vm *VM) blocked() bool {
	for _, m := range vm.machines {
		if !m.status.blocked {
			return false
		}
	}
	return true
}

func (vm *VM) unblock() {
	for i := range vm.machines {
		vm.machines[i].status.blocked = false
	}
}

type Action struct {
	Cmd  string
	Vars []interface{}
}

// Exec consumes the next event. It returns three values: 'more' indicates that
// there may be more commands in the future if more events are given; 'ok'
// indicates that there is a command to execute now, 'action' is the command to
// execute now if 'ok' is true.
func (vm *VM) Exec(next tcell.Event) (action Action, ok bool, more bool) {
	vm.evs = append(vm.evs, next)

	for {
		for i := 0; i < len(vm.machines); i++ {
			m := vm.machines[i]
			if m.status.blocked {
				continue
			}
			pc, split := m.step(vm.prog, vm.evs)
			if split {
				vm.machines = append(vm.machines, m.cpy(pc))
			}
			if m.status.done {
				// slice tricks delete
				ln := len(vm.machines)
				vm.machines[i] = vm.machines[ln-1]
				vm.machines[ln-1] = nil
				vm.machines = vm.machines[:ln-1]
				i--
				if !m.status.failed {
					action.Cmd = m.cmds[0]
					action.Vars = m.vars
					ok = true
				}
			}
		}
		if len(vm.machines) == 0 || vm.blocked() {
			break
		}
	}
	vm.unblock()
	return action, ok, len(vm.machines) > 0
}
