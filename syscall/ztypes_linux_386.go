// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs linux/types.go

//line linux/types.go:2
package syscall

//line linux/types.go:14

//line linux/types.go:13
type Termios struct {
	//line linux/types.go:13
	Iflag uint32
	//line linux/types.go:13
	Oflag uint32
	//line linux/types.go:13
	Cflag uint32
	//line linux/types.go:13
	Lflag uint32
	//line linux/types.go:13
	Line uint8
	//line linux/types.go:13
	Cc [32]uint8
	//line linux/types.go:13
	Pad_cgo_0 [3]byte
	//line linux/types.go:13
	Ispeed uint32
	//line linux/types.go:13
	Ospeed uint32
	//line linux/types.go:13
}
type Int int32

//line linux/types.go:16

//line linux/types.go:15
const (
	Hoho = 0xf0000001
)
