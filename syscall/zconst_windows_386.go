// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs ,,const.go

package syscall

const (
	MAXDWORD = 0xffffffff

	FILE_FLAG_OVERLAPPED = 0x40000000

	CLRBREAK = 0x9
	CLRDTR   = 0x6
	CLRRTS   = 0x4

	SETBREAK = 0x8
	SETDTR   = 0x5
	SETRTS   = 0x3

	SETXOFF = 0x1
	SETXON  = 0x2

	DTR_CONTROL_DISABLE   = 0x0
	DTR_CONTROL_ENABLE    = 0x1
	DTR_CONTROL_HANDSHAKE = 0x2
	RTS_CONTROL_ENABLE    = 0x1
	RTS_CONTROL_DISABLE   = 0x0
	RTS_CONTROL_HANDSHAKE = 0x2
	RTS_CONTROL_TOGGLE    = 0x3

	ONESTOPBIT  = 0x0
	TWOSTOPBITS = 0x2

	ODDPARITY  = 0x1
	EVENPARITY = 0x2
	NOPARITY   = 0x0

	TRUE  = 0x1
	FALSE = 0x0

	ERROR_SUCCESS = 0x0
)
