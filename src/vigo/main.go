package main

import (
  "bufio"
  "fmt"
  "os"
)

func enableRawMode() {
  // get terminal state

  // make adjustments

  // set terminal state

}

func main() {
  enableRawMode()

  reader := bufio.NewReader(os.Stdin)

  for {
    c, err := reader.ReadByte()
    if err != nil {
      // if error including EOF, break
      break
    }
    if c == 'q' {
      // if 'q' is entered, break
      break
    }

    fmt.Printf("%c", c)
  }

}

