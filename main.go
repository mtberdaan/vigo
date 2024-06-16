package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"
)

// ----------------- Window Size  -----------------

type size struct {
	width  int // columns
	height int // rows
}

type winsize struct {
	rows uint16
	cols uint16
	x    uint16 // pixel x
	y    uint16 // pixel y
}

type sizeListener struct {
	change <-chan size
	done   chan struct{}
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

func getTerminalSizeChanges(sc chan size, done chan struct{}) error {
	ch := make(chan os.Signal, 1)

	sig := unix.SIGWINCH

	signal.Notify(ch, sig)
	go func() {
		for {
			select {
			case <-ch:
				var err error
				s, err := getTerminalSize(os.Stdout)
				if err != nil {
					sc <- s
				}
			case <-done:
				signal.Reset(sig)
				close(ch)
				return
			}
		}
	}()
	return nil
}

func (sc *sizeListener) close() (err error) {
	if sc.done != nil {
		close(sc.done)
		sc.done = nil
		sc.change = nil
	}
	return
}

func newSizeListener() (sc *sizeListener, err error) {
	sc = &sizeListener{}

	sizechan := make(chan size, 1)
	sc.change = sizechan
	sc.done = make(chan struct{})

	err = getTerminalSizeChanges(sizechan, sc.done)
	if err != nil {
		close(sizechan)
		close(sc.done)
		sc = &sizeListener{}
		return
	}
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

	sc, err := newSizeListener()
	if err != nil {
		fmt.Println("Error getting terminal size listener:", err)
		return
	}
	defer sc.close()

	for {
		select {
		case s = <-sc.change:
			fmt.Println("terminal size:", s)
		}
	}
}
