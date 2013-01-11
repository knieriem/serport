// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs windows/types.go

package syscall

type DCB struct {
	DCBlength  uint32
	BaudRate   uint32
	Flags      uint32
	WReserved  uint16
	XonLim     uint16
	XoffLim    uint16
	ByteSize   uint8
	Parity     uint8
	StopBits   uint8
	XonChar    int8
	XoffChar   int8
	ErrorChar  int8
	EofChar    int8
	EvtChar    int8
	WReserved1 uint16
}

type CommTimeouts struct {
	ReadIntervalTimeout         uint32
	ReadTotalTimeoutMultiplier  uint32
	ReadTotalTimeoutConstant    uint32
	WriteTotalTimeoutMultiplier uint32
	WriteTotalTimeoutConstant   uint32
}

const (
	dcbSize = 0x1c

	HKEY_CLASSES_ROOT   = 0x80000000
	HKEY_CURRENT_CONFIG = 0x80000005
	HKEY_CURRENT_USER   = 0x80000001
	HKEY_LOCAL_MACHINE  = 0x80000002
)

type REGSAM uint32
