package main

import (
  "os"

  "golang.org/x/term"
)

type editorConfig struct {
  screenrows int
  screencols int
}

var E editorConfig

func getTermSize(fd int) editorConfig {
  c, r, err := term.GetSize(fd)
  
  if err != nil {
    panic(err)
  }

  return editorConfig{screenrows: r, screencols: c} 
}

func main() {
  E := getTermSize(int(os.Stdin.Fd()))
  println("Rows: ", E.screenrows)
  println("Cols: ", E.screencols)
}
