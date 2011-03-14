// godefs -g syscall -c /usr/bin/i586-mingw32msvc-gcc ,,const.c

// MACHINE GENERATED - DO NOT EDIT.

package syscall

// Constants
const (
	MAXDWORD                = (1 << 32) - 0x1
	FILE_FLAG_OVERLAPPED    = 0x40000000
	CLRBREAK                = 0x9
	CLRDTR                  = 0x6
	CLRRTS                  = 0x4
	SETBREAK                = 0x8
	SETDTR                  = 0x5
	SETRTS                  = 0x3
	SETXOFF                 = 0x1
	SETXON                  = 0x2
	DTR_CONTROL_DISABLE     = 0
	DTR_CONTROL_ENABLE      = 0x1
	DTR_CONTROL_HANDSHAKE   = 0x2
	RTS_CONTROL_ENABLE      = 0x1
	RTS_CONTROL_DISABLE     = 0
	RTS_CONTROL_HANDSHAKE   = 0x2
	RTS_CONTROL_TOGGLE      = 0x3
	ONESTOPBIT              = 0
	TWOSTOPBITS             = 0x2
	ODDPARITY               = 0x1
	EVENPARITY              = 0x2
	NOPARITY                = 0
	TRUE                    = 0x1
	FALSE                   = 0
	ERROR_SUCCESS           = 0
	REG_BINARY              = 0x3
	REG_DWORD               = 0x4
	REG_DWORD_LITTLE_ENDIAN = 0x4
	REG_DWORD_BIG_ENDIAN    = 0x5
	REG_QWORD               = 0xb
	REG_SZ                  = 0x1
	REG_MULTI_SZ            = 0x7
	KEY_ALL_ACCESS          = 0xf003f
	KEY_CREATE_LINK         = 0x20
	KEY_CREATE_SUB_KEY      = 0x4
	KEY_ENUMERATE_SUB_KEYS  = 0x8
	KEY_EXECUTE             = 0x20019
	KEY_NOTIFY              = 0x10
	KEY_QUERY_VALUE         = 0x1
	KEY_READ                = 0x20019
	KEY_SET_VALUE           = 0x2
	KEY_WRITE               = 0x20006
)

// Types
