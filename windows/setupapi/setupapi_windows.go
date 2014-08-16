// Some utility functions from SetupAPI.
// (see "Public Device Installation Functions",
//	http://msdn.microsoft.com/en-us/library/ff549791.aspx)
//
package setupapi

// http://support.microsoft.com/kb/259695/en-us

var GuidSerialPorts = &Guid{
	Data1: 0x4D36E978,
	Data2: 0xE325,
	Data3: 0x11CE,
	Data4: [8]uint8{
		0xBF, 0xC1, 0x08, 0x00, 0x2B, 0xE1, 0x03, 0x18,
	},
}

//sys SetupDiGetClassDevs(class *Guid, enum *uint16, parent syscall.Handle, flags uint32) (devInfoSet syscall.Handle, err error) [failretval==syscall.InvalidHandle] = setupapi.SetupDiGetClassDevsW

//sys SetupDiGetDeviceRegistryProperty(devInfoSet syscall.Handle, diData *SpDevinfoData, prop uint32, regDataType *uint32, buf []byte, size *uint32) (err error) [failretval == 0] = setupapi.SetupDiGetDeviceRegistryPropertyW

//sys SetupDiEnumDeviceInfo(devInfoSet syscall.Handle, index uint32, diData *SpDevinfoData) (err error) [failretval == 0] = setupapi.SetupDiEnumDeviceInfo

//sys SetupDiCreateDeviceInfo(devInfoSet syscall.Handle, devName *uint16, g *Guid, devDesc *uint16, hwnd uintptr, cflags uint32, dataOut *SpDevinfoData) (err error) = setupapi.SetupDiCreateDeviceInfoW

//sys SetupDiCreateDeviceInfoList(g *Guid, hwnd uintptr) (devInfoSet syscall.Handle, err error) [failretval==syscall.InvalidHandle] = setupapi.SetupDiCreateDeviceInfoList

//sys SetupDiSetDeviceRegistryProperty(devInfoSet syscall.Handle, data *SpDevinfoData, prop uint32, buf *byte, sz uint32) (err error) = setupapi.SetupDiSetDeviceRegistryPropertyW

//sys SetupDiCallClassInstaller(installFn uintptr, devInfoSet syscall.Handle, data *SpDevinfoData) (err error) = setupapi.SetupDiCallClassInstaller

//sys SetupDiDestroyDeviceInfoList(devInfoSet syscall.Handle) (err error) = setupapi.SetupDiDestroyDeviceInfoList

//sys	SetupDiGetINFClass(infPath *uint16, guid *Guid, className []uint16, reqSz *uint32) (err error) = setupapi.SetupDiGetINFClassW

//sys SetupDiOpenDevRegKey(devInfoSet syscall.Handle, diData *SpDevinfoData, scope uint32, hwProfile uint32, keyType uint32, desiredAccess uint32) (h syscall.Handle, err error) [failretval==syscall.InvalidHandle] = setupapi.SetupDiOpenDevRegKey

//sys SetupDiGetDeviceInstanceId(devInfoSet syscall.Handle, diData *SpDevinfoData, id []uint16, reqSz *uint32) (err error) = setupapi.SetupDiGetDeviceInstanceIdW

func NewDevinfoData() *SpDevinfoData {
	d := new(SpDevinfoData)
	d.CbSize = SpDevinfoDataSz
	return d
}
