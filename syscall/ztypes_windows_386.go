// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs windows/types.go

//line windows/types.go:1
package syscall

//line windows/types.go:21

//line windows/types.go:20
type DCB struct {
	//line windows/types.go:20
	DCBlength uint32
	//line windows/types.go:20
	BaudRate uint32
	//line windows/types.go:20
	Flags uint32
	//line windows/types.go:20
	WReserved uint16
	//line windows/types.go:20
	XonLim uint16
	//line windows/types.go:20
	XoffLim uint16
	//line windows/types.go:20
	ByteSize uint8
	//line windows/types.go:20
	Parity uint8
	//line windows/types.go:20
	StopBits uint8
	//line windows/types.go:20
	XonChar int8
	//line windows/types.go:20
	XoffChar int8
	//line windows/types.go:20
	ErrorChar int8
	//line windows/types.go:20
	EofChar int8
	//line windows/types.go:20
	EvtChar int8
	//line windows/types.go:20
	WReserved1 uint16
	//line windows/types.go:20
}

//line windows/types.go:23

//line windows/types.go:22
type CommTimeouts struct {
	//line windows/types.go:22
	ReadIntervalTimeout uint32
	//line windows/types.go:22
	ReadTotalTimeoutMultiplier uint32
	//line windows/types.go:22
	ReadTotalTimeoutConstant uint32
	//line windows/types.go:22
	WriteTotalTimeoutMultiplier uint32
	//line windows/types.go:22
	WriteTotalTimeoutConstant uint32
	//line windows/types.go:22
}

//line windows/types.go:25

//line windows/types.go:24
const (
	dcbSize = 0x1c
	//line windows/types.go:28

	//line windows/types.go:27
	HKEY_CLASSES_ROOT   = (1 << 32) - 0x80000000
	HKEY_CURRENT_CONFIG = (1 << 32) - 0x7ffffffb
	HKEY_CURRENT_USER   = (1 << 32) - 0x7fffffff
	HKEY_LOCAL_MACHINE  = (1 << 32) - 0x7ffffffe
)

//line windows/types.go:34

//line windows/types.go:33
type REGSAM uint32
