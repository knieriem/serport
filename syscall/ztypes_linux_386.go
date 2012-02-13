// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs linux/types.go

//line linux/types.go:3
package syscall

//line linux/types.go:12

//line linux/types.go:11
type Termios struct {
	//line linux/types.go:11
	Iflag uint32
	//line linux/types.go:11
	Oflag uint32
	//line linux/types.go:11
	Cflag uint32
	//line linux/types.go:11
	Lflag uint32
	//line linux/types.go:11
	Line uint8
	//line linux/types.go:11
	Cc [32]uint8
	//line linux/types.go:11
	Pad_cgo_0 [3]byte
	//line linux/types.go:11
	Ispeed uint32
	//line linux/types.go:11
	Ospeed uint32
	//line linux/types.go:11
}
type Int int32
