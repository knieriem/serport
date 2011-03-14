// godefs -g syscall -f-m32 -c /usr/bin/i586-mingw32msvc-gcc types.c

// MACHINE GENERATED - DO NOT EDIT.

package syscall

// Constants
const (
	dcbSize             = 0x1c
	HKEY_CLASSES_ROOT   = (1 << 32) - 0x80000000
	HKEY_CURRENT_CONFIG = (1 << 32) - 0x7ffffffb
	HKEY_CURRENT_USER   = (1 << 32) - 0x7fffffff
	HKEY_LOCAL_MACHINE  = (1 << 32) - 0x7ffffffe
)

// Types

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

type REGSAM uint32
