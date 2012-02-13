// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs ,,const.go

//line ,,const.go:1
package syscall

//line ,,const.go:9

//line ,,const.go:8
const (
	MAXDWORD = (1 << 32) - 0x1
	//line ,,const.go:13

	//line ,,const.go:12
	FILE_FLAG_OVERLAPPED = 0x40000000
	//line ,,const.go:16

	//line ,,const.go:15
	CLRBREAK = 0x9
	CLRDTR   = 0x6
	CLRRTS   = 0x4
	//line ,,const.go:20

	//line ,,const.go:19
	SETBREAK = 0x8
	SETDTR   = 0x5
	SETRTS   = 0x3
	//line ,,const.go:24

	//line ,,const.go:23
	SETXOFF = 0x1
	SETXON  = 0x2
	//line ,,const.go:28

	//line ,,const.go:27
	DTR_CONTROL_DISABLE   = 0x0
	DTR_CONTROL_ENABLE    = 0x1
	DTR_CONTROL_HANDSHAKE = 0x2
	RTS_CONTROL_ENABLE    = 0x1
	RTS_CONTROL_DISABLE   = 0x0
	RTS_CONTROL_HANDSHAKE = 0x2
	RTS_CONTROL_TOGGLE    = 0x3
	//line ,,const.go:36

	//line ,,const.go:35
	ONESTOPBIT  = 0x0
	TWOSTOPBITS = 0x2
	//line ,,const.go:39

	//line ,,const.go:38
	ODDPARITY  = 0x1
	EVENPARITY = 0x2
	NOPARITY   = 0x0
	//line ,,const.go:43

	//line ,,const.go:42
	TRUE  = 0x1
	FALSE = 0x0
	//line ,,const.go:46

	//line ,,const.go:45
	ERROR_SUCCESS = 0x0
	//line ,,const.go:50

	//line ,,const.go:49
	REG_BINARY              = 0x3
	REG_DWORD               = 0x4
	REG_DWORD_LITTLE_ENDIAN = 0x4
	REG_DWORD_BIG_ENDIAN    = 0x5
	//line ,,const.go:55

	//line ,,const.go:54
	REG_QWORD              = 0xb
	REG_SZ                 = 0x1
	REG_MULTI_SZ           = 0x7
	KEY_ALL_ACCESS         = 0xf003f
	KEY_CREATE_LINK        = 0x20
	KEY_CREATE_SUB_KEY     = 0x4
	KEY_ENUMERATE_SUB_KEYS = 0x8
	KEY_EXECUTE            = 0x20019
	KEY_NOTIFY             = 0x10
	KEY_QUERY_VALUE        = 0x1
	KEY_READ               = 0x20019
	KEY_SET_VALUE          = 0x2
	//line ,,const.go:69

	//line ,,const.go:68
	KEY_WRITE = 0x20006
)
