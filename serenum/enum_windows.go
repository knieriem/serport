package serenum

import (
	"regexp"
	"sort"
	"strings"

	"github.com/knieriem/g/windows/setupapi"
)

// Ports gathers information about the serial ports present on a system,
// making use of information provided by the SetupDi* functions.
// The resulting list is sorted in the following order:
//   - PL2303 devices
//   - FTDIBUS devices
//   - USB devices
//   - ACPI devices
//   - other devices
//
// Line printer devices (LPT...) will be skipped.
func Ports() (ports []*PortInfo) {
	walkDevices(func(name string, di *setupapi.DevinfoSet) (match bool) {
		ports = append(ports, makeInfo(name, di))
		return
	})
	sort.Sort(portList(ports))
	return
}

// Lookup returns information about the named serial port. On failure,
// a PortInfo with just the Device field set to portName will be returned.
func Lookup(portName string) (port *PortInfo) {
	want := strings.ToUpper(portName)
	match := walkDevices(func(name string, di *setupapi.DevinfoSet) (match bool) {
		have := strings.ToUpper(name)
		if want != have {
			return
		}
		match = true
		port = makeInfo(name, di)
		return
	})
	if !match {
		port = &PortInfo{Device: portName}
	}
	return
}

type walkFunc func(portName string, di *setupapi.DevinfoSet) (match bool)

func walkDevices(f walkFunc) (match bool) {
	h, err := setupapi.SetupDiGetClassDevs(setupapi.GuidSerialPorts, nil, 0, setupapi.DIGCF_PRESENT)
	if err != nil {
		return
	}
	defer setupapi.SetupDiDestroyDeviceInfoList(h)

	di := setupapi.NewDevinfoSet(h)
	for i := 0; di.Enum(i); i++ {
		portName := ""
		if k, err := di.OpenDevRegKey(); err == nil {
			v, err := k.Value("PortName")
			if err == nil {
				portName = v.String()
			}
		}
		if portName == "" || strings.HasPrefix(portName, "LPT") {
			continue
		}
		match = f(portName, di)
		if match {
			break
		}
	}
	return
}

func makeInfo(portName string, di *setupapi.DevinfoSet) (port *PortInfo) {
	port = new(PortInfo)
	port.Device = portName
	port.Desc = di.DeviceRegistryProperty(setupapi.SPDRP_DEVICEDESC).String()
	port.Enumerator = di.DeviceRegistryProperty(setupapi.SPDRP_ENUMERATOR_NAME).String()

	// Try to extract vendor and product IDs.
	f := di.DeviceRegistryProperty(setupapi.SPDRP_HARDWAREID).Strings()
	if len(f) > 0 {
		s := f[0]
		s = strings.ToUpper(s)
		s = strings.Replace(s, "\\", "&", -1)
		for _, f := range strings.Split(s, "&") {
			if strings.HasPrefix(f, "VID_") {
				port.VendorID = f[4:]
			}
			if strings.HasPrefix(f, "PID_") {
				port.ProductID = f[4:]
			}
		}
	}
	if id, err := di.DeviceInstanceID(); err == nil {
		if m := ftdiSerialNumRE.FindStringSubmatch(id); len(m) == 2 {
			port.SerialNumber = m[1]
		}
	}
	return
}

var ftdiSerialNumRE = regexp.MustCompile(`FTDIBUS\\VID_\w+\+PID_\w+\+(\w+)A\\00.*`)

func (list portList) Less(i, j int) bool {
	p0, p1 := list[i], list[j]
	if isLess, match := matchPL2303(p0, p1); match {
		return isLess
	}
	if isLess, match := matchEnumerator(p0, p1, "FTDIBUS"); match {
		return isLess
	}
	if isLess, match := matchDesc(p0, p1, "roadband", false); match {
		return isLess
	}
	if isLess, match := matchEnumerator(p0, p1, "USB"); match {
		return isLess
	}
	if isLess, match := matchEnumerator(p0, p1, "ACPI"); match {
		return isLess
	}
	return less(p0, p1)
}

func matchEnumerator(p0, p1 *PortInfo, e string) (isLess bool, match bool) {
	if p0.Enumerator == e {
		match = true
		if p1.Enumerator == e {
			isLess = less(p0, p1)
		} else {
			isLess = true
		}
	}
	return
}
