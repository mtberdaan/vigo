package main

import (
	"fmt"
	"os"
	"unicode"

	"golang.org/x/term"
	"golang.org/x/sys/unix"
)

type editorConfig struct {
	originTermState unix.Termios 
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
  termios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
  if err != nil {
    panic(err)
  } 

  e.originTermState = *termios

  termios.Lflag &^= unix.ECHO | unix.ICANON

  if err := unix.IoctlSetTermios(fd, unix.TIOCSETA, termios); err != nil {
    panic(err)
  }


}

func (e *editorConfig) disableRawMode(fd int) {
  err := unix.IoctlSetTermios(fd, unix.TIOCSETA, &e.originTermState)
  if err != nil {
    panic(err)
  }

}

func main() {
	fd := int(os.Stdin.Fd())

	E.getTermSize(fd)
	E.setRawMode(fd)

	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil || buf[0] == 'q' {
			break
		}
		if unicode.IsControl(rune(buf[0])) {
			fmt.Println("%d\n", buf[0])
		} else {
			fmt.Println("%d ('%c')\n", buf[0], buf[0])
		}
	}

	E.disableRawMode(fd)
}
