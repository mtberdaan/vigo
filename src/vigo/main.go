package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"
)

// ----------------- Application Flow -----------------
// 1. Term input (vigo filename)
// 2. Read file -> buffer
// 3. Display buffer -> screen

type vigo struct {
	buffers []*buffer
}

type buffer struct {
	original string
	add      string
	pieces   []piece

	// file path if read from file
	path string

	// unique bugger name
	name string
}

type piece struct {
	start  int
	length int
	source string
}

func new_vigo(filenames []string) *vigo {
	v := new(vigo)
	v.buffers = make([]*buffer, 0, 20)
	for _, filename := range filenames {
		v.new_buffer_from_file(filename)
	}
	// create new empty buffer if no file
	return v
}

func (v *vigo) new_buffer_from_file(filename string) (*buffer, error) {
	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	// what if file is not there?

	f, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := new_buffer(f)
	if err != nil {
		return nil, err
	}
	buf.path = fullpath
	buf.name = filename

	v.buffers = append(v.buffers, buf)

	return buf, nil
}

func new_buffer(r io.Reader) (*buffer, error) {
	buf := new(buffer)

	bo, err := io.ReadAll(r)
	if err != nil && err != io.EOF {
		return nil, err
	}

	buf.original = string(bo)
	// read file into buffer
	return buf, nil
}

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

// ----------------- View Functions -----------------

func draw() {
	// clear screen
	// draw status bar
	// draw text buffer
}

// ----------------- Status Bar  -----------------

// ----------------- Key Controls -----------------

// ----------------- Text Buffer -----------------

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

// ----------------- Debug Functions -----------------
func showTermSize() {
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

// ----------------- Main Loop -----------------

func main() {
	vigo := new_vigo(os.Args[1:])

}
