package main

import (
	"os"

	"golang.org/x/term"
)

type editorConfig struct {
	originTermState *term.State
	screenrows      int
	screencols      int
}

var E editorConfig

/*** init ***/

func (e *editorConfig) getTermSize(fd int) {
	c, r, err := term.GetSize(fd)

	if err != nil {
		panic(err)
	}

	e.screenrows = r
	e.screencols = c
}

func (e *editorConfig) setRawMode(fd int) {
	originState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	e.originTermState = originState
}

func (e *editorConfig) disableRawMode(fd int) {
	term.Restore(fd, e.originTermState)
}


func main() {
	terminal := int(os.Stdin.Fd())
  
  E.getTermSize(terminal)
  E.setRawMode(terminal)

	println("Rows: ", E.screenrows)
	println("Cols: ", E.screencols)

	E.disableRawMode(terminal)
}
