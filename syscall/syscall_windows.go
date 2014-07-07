package syscall

import (
	"syscall"
)

const (
	// for CreateEvent
	EvManualReset = true
	EvInitiallyOn = true
)

func CreateEvent(manualReset, initialState bool) (h syscall.Handle, err error) {
	return CreateEventW(nil, b2i(manualReset), b2i(initialState), nil)
}
func b2i(v bool) int {
	if v {
		return 1
	}
	return 0
}

//sys CreateEventW(sa *syscall.SecurityAttributes, manualReset int, initialState int, name *uint16) (hEv syscall.Handle, err error)

//sys SetEvent(h syscall.Handle) (err error)

//sys GetOverlappedResult(h syscall.Handle, ov *syscall.Overlapped, done *uint32, bWait int) (err error)

//sys	EscapeCommFunction(h syscall.Handle, fn uint32) (err error)
//sys SetupComm(h syscall.Handle, inQSize uint32, outQSize uint32) (err error)
//sys SetCommTimeouts(h syscall.Handle, cto *CommTimeouts) (err error)
//sys setCommState(h syscall.Handle, dcb *DCB) (err error) = SetCommState
//sys getCommState(h syscall.Handle, dcb *DCB) (err error) = GetCommState
//sys FlushFileBuffers(h syscall.Handle) (err error)

// Flags for DCB.Flags (simulating a bitfield)
const (
	DCBpBinary, DCBfBinary = iota, 1 << iota
	DCBpParity, DCBfParity
	DCBpOutxCtsFlow, DCBfOutxCtsFlow
	DCBpOutxDsrFlow, DCBfOutxDsrFlow
	DCBpDtrControl, _
	_, _
	DCBpDsrSensitivity, DCBfDsrSensitivity
	DCBpTXContinueOnXoff, DCBfTXContinueOnXoff
	DCBpOutX, DCBfOutX
	DCBpInX, DCBfInX
	DCBpErrorChar, DCBfErrorChar
	DCBpNull, DCBfNull
	DCBpRtsControl, _
	_, _
	DCBpAbortOnError, DCBfAbortOnError

	DCBmDtrControl = 3
	DCBmRtsControl = 3
)

func SetCommState(h syscall.Handle, dcb *DCB) (err error) {
	dcb.DCBlength = dcbSize
	dcb.Flags |= DCBfBinary
	return setCommState(h, dcb)
}
func GetCommState(h syscall.Handle, dcb *DCB) (err error) {
	dcb.DCBlength = dcbSize
	return getCommState(h, dcb)
}

//sys GetConsoleMode(h syscall.Handle, mode *uint32) (err error) [failretval==FALSE]

// registry stuff
type HKEY uintptr

//sys RegOpenKeyEx(h HKEY, name *uint16, options uint32, samDesired REGSAM, result *HKEY) (err error) [failretval!=ERROR_SUCCESS] = advapi32.RegOpenKeyExW
//sys RegEnumValue(h HKEY, index uint32, vName *uint16, vNameLen *uint32, reserved *uint32, typ *uint32, data *byte, sz *uint32) (err error) [failretval!=ERROR_SUCCESS] = advapi32.RegEnumValueW
//sys RegQueryValueEx(h HKEY, vName *uint16, reserved *uint32, typ *uint32, data *byte, sz *uint32) (err error) [failretval!=ERROR_SUCCESS] = advapi32.RegQueryValueExW
//sys RegCloseKey(h HKEY) [failretval!=ERROR_SUCCESS] = advapi32.RegCloseKey
