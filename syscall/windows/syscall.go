package syscall

import (
	"syscall"
	"log"
	"unsafe"
)


type Handle uint32

func (h Handle) Close() {
	syscall.CloseHandle(int32(h))
}
func (h Handle) Byteptr() *byte {
	return (*byte)(unsafe.Pointer(uintptr(h)))
}

const (
	// for CreateEvent
	EvManualReset = true
	EvInitiallyOn = true
)

func CreateEvent(manualReset, initialState bool) (h Handle, errno int) {
	return CreateEventW(nil, b2i(manualReset), b2i(initialState), nil)
}
func b2i(v bool) int {
	if v {
		return 1
	}
	return 0
}

//sys CreateEventW(sa *syscall.SecurityAttributes, manualReset int, initialState int, name *uint16) (hEv Handle, errno int)


//sys GetOverlappedResult(h uint32, ov *syscall.Overlapped, done *uint32, bWait int) (errno int)


//sys	EscapeCommFunction(h uint32, fn uint32) (errno int)
//sys SetupComm(h uint32, inQSize uint32, outQSize uint32) (errno int)
//sys SetCommTimeouts(h uint32, cto *CommTimeouts) (errno int)
//sys setCommState(h uint32, dcb *DCB) (errno  int) = SetCommState
//sys getCommState(h uint32, dcb *DCB) (errno int) = GetCommState
//sys FlushFileBuffers(h uint32) (errno int)

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

func SetCommState(h uint32, dcb *DCB) (errno int) {
	dcb.DCBlength = dcbSize
	dcb.Flags |= DCBfBinary
	return setCommState(h, dcb)
}
func GetCommState(h uint32, dcb *DCB) (errno int) {
	dcb.DCBlength = dcbSize
	return getCommState(h, dcb)
}


// registry stuff
type HKEY uintptr

//sys RegOpenKeyEx(h HKEY, name *uint16, options uint32, samDesired REGSAM, result *HKEY) (errno int) [failretval!=ERROR_SUCCESS] = advapi32.RegOpenKeyExW
//sys RegEnumValue(h HKEY, index uint32, vName *uint16, vNameLen *uint32, reserved *uint32, typ *uint32, data *byte, sz *uint32) (errno int) [failretval!=ERROR_SUCCESS] = advapi32.RegEnumValueW
//sys RegQueryValueEx(h HKEY, vName *uint16, reserved *uint32, typ *uint32, data *byte, sz *uint32) (errno int) [failretval!=ERROR_SUCCESS] = advapi32.RegQueryValueExW
//sys RegCloseKey(h HKEY) [failretval!=ERROR_SUCCESS] = advapi32.RegCloseKey


//sys getUserName(buf *uint16, sz *uint32) (errno int) [failretval==0] = advapi32.GetUserNameW

func GetUserName() string {
	var (
		buf = make([]uint16, 128)
		sz  = uint32(len(buf))
	)
	if e := getUserName(&buf[0], &sz); e != 0 {
		return "none"
	}
	return syscall.UTF16ToString(buf[:sz])
}


func loadDll(fname string) uint32 {
	h, e := syscall.LoadLibrary(fname)
	if e != 0 {
		log.Fatalf("LoadLibrary(%s) failed with err=%d.\n", fname, e)
	}
	return h
}
func getSysProcAddr(m uint32, pname string) uintptr {
	p, e := syscall.GetProcAddress(m, pname)
	if e != 0 {
		log.Fatalf("GetProcAddress(%s) failed with err=%d.\n", pname, e)
	}
	return uintptr(p)
}
