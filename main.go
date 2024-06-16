package main

import (
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"
)

// ----------------- Window Size  -----------------

type size struct {
	width  int
	height int
}

type winsize struct {
	rows uint16
	cols uint16
	x    uint16
	y    uint16
}

func getSize() (s size, err error) {
	//TODO: check if terminal

	s, err = getTerminalSize(os.Stdout)
	return
}

func getTerminalSize(fp *os.File) (s size, err error) {
	ws := winsize{}

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		fp.Fd(),
		uintptr(unix.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))

	if errno != 0 {
		err = errno
		return
	}

	s = size{
		width:  int(ws.cols),
		height: int(ws.rows),
	}

	fmt.Println("size:", ws)

	return
}

// ----------------- Text Buffer -----------------
type buffer struct {
	original string
	add      string
	pieces   []piece
}

type piece struct {
	start  int
	length int
	source string
}

func (b *buffer) display() string {
	value := reflect.ValueOf(b) // get value of buffer

	for _, piece := range b.pieces {
		source := piece.source // origin or add (select source)
		buffer := value.FieldByName(source).Interface().(string)
		span_of_text := buffer[piece.start : piece.start+piece.length]
		fmt.Println(span_of_text)
	}

	return ""
}

func main() {
	s, err := getSize()
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return
	}

	fmt.Println("terminal size:", s)
}
