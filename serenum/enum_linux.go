package serenum

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Ports gathers information about the serial ports present on a system, by
// examining information from sysfs, especially from /sys/class/tty.
// Devices that have no driver entry are skipped, as well as devices that symlink
// into /sys/devices/platform. If a device is an USB device, additional
// information like VID, PID will be extracted.
// The resulting list containing one item for each serial port is sorted
// in the following order:
//   - PL2303 devices
//   - USB devices
//   - ACM devices
//   - other devices
func Ports() (ports []*PortInfo) {
	f, err := os.Open("/sys/class/tty")
	if err != nil {
		return
	}
	names, err := f.Readdirnames(0)
	if err != nil {
		return
	}
	for _, name := range names {
		if info, ok := readDeviceInfo(name); ok {
			ports = append(ports, info)
		}
	}
	sort.Sort(portList(ports))
	return
}

// Lookup returns information about the named serial port. On failure,
// a PortInfo with just the Device field set to portName will be returned.
func Lookup(portName string) *PortInfo {
	name := strings.TrimPrefix(portName, "/dev/")
	if info, ok := readDeviceInfo(name); ok {
		return info
	}
	return &PortInfo{Device: portName}
}

func readDeviceInfo(name string) (port *PortInfo, ok bool) {
	path := linkTarget("/sys/class/tty", name)

	// skip platform devices -- they appear not to be real devices
	if strings.HasPrefix(path, "/sys/devices/platform/serial8250/") {
		return
	}

	// skip device if it does not have a `driver' symlink
	drvPath := filepath.Join(path, "device", "driver")
	drv, err := os.Lstat(drvPath)
	if err != nil || drv.Mode()&os.ModeSymlink == 0 {
		return
	}

	link, err := os.Readlink(drvPath)
	if err != nil {
		return
	}

	port = new(PortInfo)
	port.Driver = filepath.Base(link)

	port.Device = "/dev/" + name
	if _, err = os.Stat(port.Device); err != nil {
		return
	}

	// try to get some information about the device
	if !readUSBInfo(port, linkTarget(path, "device")) {
		readUSBInfo(port, linkTarget("/sys/bus/usb-serial/devices", name))
	}
	ok = true
	return
}

func linkTarget(path, fileName string) string {
	p := filepath.Join(path, fileName)
	link, err := os.Readlink(p)
	if err == nil {
		if !filepath.IsAbs(link) {
			link = filepath.Join(path, link)
		}
		return link
	}
	return p
}

func readUSBInfo(p *PortInfo, path string) bool {
	for i := 0; i < 3; i++ {
		if p.Desc == "" {
			p.Desc = readTextfile(path, "interface")
		}
		if _, err := os.Stat(filepath.Join(path, "idProduct")); err != nil {
			path = filepath.Dir(path)
			continue
		}
		if p.Desc == "" {
			p.Desc = readTextfile(path, "product")
		}
		p.VendorID += readTextfile(path, "idVendor")
		p.ProductID += readTextfile(path, "idProduct")
		p.Manufacturer = readTextfile(path, "manufacturer")
		p.SerialNumber = readTextfile(path, "serial")
		return true
	}
	return false
}

func readTextfile(path, fileName string) string {
	b, err := ioutil.ReadFile(filepath.Join(path, fileName))
	if err != nil {
		return ""
	}
	f := bytes.SplitN(b, []byte{'\n'}, 2)
	return string(f[0])
}

func (list portList) Less(i, j int) bool {
	p0, p1 := list[i], list[j]
	if isLess, match := matchPL2303(p0, p1); match {
		return isLess
	}
	if isLess, match := matchDesc(p0, p1, "roadband", false); match {
		return isLess
	}
	if isLess, match := matchDevice(p0, p1, "USB"); match {
		return isLess
	}
	if isLess, match := matchDevice(p0, p1, "ACM"); match {
		return isLess
	}
	return less(p0, p1)
}

func matchDevice(p0, p1 *PortInfo, s string) (isLess bool, match bool) {
	c0 := strings.Contains(p0.Device, s)
	c1 := strings.Contains(p1.Device, s)
	if c0 {
		match = true
		if c1 {
			isLess = less(p0, p1)
		} else {
			isLess = true
		}
	} else if c1 {
		match = true
	}
	return
}
