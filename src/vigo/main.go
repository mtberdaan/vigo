package main

import (
  "bufio"
  "fmt"
  "os"
)

func main() {
  reader := bufio.NewReader(os.Stdin)

  for {
    c, err := reader.ReadByte()
    if err != nil {
      // if error including EOF, break
      break
    }
    fmt.Printf("%c", c)
  }

}

