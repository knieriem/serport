// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs windows/types.go

//line windows/types.go:3
package syscall

//line windows/types.go:23

//line windows/types.go:22
type DCB struct {
	//line windows/types.go:22
	DCBlength uint32
	//line windows/types.go:22
	BaudRate uint32
	//line windows/types.go:22
	Flags uint32
	//line windows/types.go:22
	WReserved uint16
	//line windows/types.go:22
	XonLim uint16
	//line windows/types.go:22
	XoffLim uint16
	//line windows/types.go:22
	ByteSize uint8
	//line windows/types.go:22
	Parity uint8
	//line windows/types.go:22
	StopBits uint8
	//line windows/types.go:22
	XonChar int8
	//line windows/types.go:22
	XoffChar int8
	//line windows/types.go:22
	ErrorChar int8
	//line windows/types.go:22
	EofChar int8
	//line windows/types.go:22
	EvtChar int8
	//line windows/types.go:22
	WReserved1 uint16
	//line windows/types.go:22
}

//line windows/types.go:25

//line windows/types.go:24
type CommTimeouts struct {
	//line windows/types.go:24
	ReadIntervalTimeout uint32
	//line windows/types.go:24
	ReadTotalTimeoutMultiplier uint32
	//line windows/types.go:24
	ReadTotalTimeoutConstant uint32
	//line windows/types.go:24
	WriteTotalTimeoutMultiplier uint32
	//line windows/types.go:24
	WriteTotalTimeoutConstant uint32
	//line windows/types.go:24
}

//line windows/types.go:27

//line windows/types.go:26
const (
	dcbSize = 0x1c
	//line windows/types.go:30

	//line windows/types.go:29
	HKEY_CLASSES_ROOT   = (1 << 32) - 0x80000000
	HKEY_CURRENT_CONFIG = (1 << 32) - 0x7ffffffb
	HKEY_CURRENT_USER   = (1 << 32) - 0x7fffffff
	HKEY_LOCAL_MACHINE  = (1 << 32) - 0x7ffffffe
)

//line windows/types.go:36

//line windows/types.go:35
type REGSAM uint32
