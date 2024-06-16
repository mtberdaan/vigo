package main

import ("reflect")

// ----------------- Window Size  -----------------

type Window struct {
  Width int
  Height int
}


// ----------------- Text Buffer -----------------
type Buffer struct {
  original string
  add string
  pieces []Piece
}

type Piece struct {
  start int
  length int
  source string
}

func (b *Buffer) Display() string {
  value := reflect.ValueOf(b) // get value of buffer
  
  for piece := range b.pieces {
    source string := piece.source // origin or add (select source)
    buffer string := value.FieldByName(source) 
    span_of_text = buffer[piece.start:piece.start + piece.length]
    fmt.print(span_of_text) 
}

func main() {
  
}
