package main

import (
	"fmt"
	"log"
	"os"

	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/kbd"
)

func main() {
	f, err := os.Create("log.txt")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	prog := vim()
	log.Println(prog.Compile())

	vm := kbd.NewVM(prog.Compile())

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	for {
		ev := s.PollEvent()
		action, ok, more := vm.Exec(ev)
		log.Println(action.Cmd, ok, more)
		for i, v := range action.Vars {
			log.Printf("\t$%d: %v\n", i, v)
		}
		if !more {
			vm.Reset()
		}
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		default:
		}
	}

}
