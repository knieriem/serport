// +build windows

package registry

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	win "github.com/knieriem/g/syscall"
	"runtime"
	"syscall"
)

type Key struct {
	syscall.Handle
}

var (
	KeyClassesRoot   = &Key{syscall.HKEY_CLASSES_ROOT}
	KeyLocalMachine  = &Key{syscall.HKEY_LOCAL_MACHINE}
	KeyCurrentUser   = &Key{syscall.HKEY_CURRENT_USER}
	KeyCurrentConfig = &Key{syscall.HKEY_CURRENT_CONFIG}
)

func (k *Key) Subkey(path ...string) (result *Key, err error) {
	var key syscall.Handle

	s := ""
	for i, v := range path {
		if i > 0 {
			s += "\\" + v
		} else {
			s += v
		}
	}
	u, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		return
	}
	err = syscall.RegOpenKeyEx(k.Handle, u, 0, syscall.KEY_READ, &key)
	if err == nil {
		result = &Key{key}
		runtime.SetFinalizer(result, (*Key).Close)
	}
	return
}

type KeyBaseInfo struct {
	Name      string
	LastWrite time.Time
}

type KeyInfo struct {
	KeyBaseInfo

	NumSubKeys   int
	MaxSubKeyLen int

	MaxClassLen int

	NumValues       int
	MaxValueNameLen int
	MaxValueLen     int
}

func (k *Key) Subkeys() (list []KeyBaseInfo, err error) {
	var n, maxNameLen, ulen uint32
	var ft syscall.Filetime

	err = syscall.RegQueryInfoKey(k.Handle, nil, nil, nil, &n, &maxNameLen, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return
	}

	maxNameLen++
	ubuf := make([]uint16, maxNameLen)

	list = make([]KeyBaseInfo, n)
	for i := uint32(0); i < n; i++ {
		ulen = maxNameLen
		err = syscall.RegEnumKeyEx(k.Handle, i, &ubuf[0], &ulen, nil, nil, nil, &ft)
		if err != nil {
			break
		}
		list[i] = KeyBaseInfo{
			Name:      syscall.UTF16ToString(ubuf[:ulen]),
			LastWrite: time.Unix(0, ft.Nanoseconds()),
		}
	}

	return
}

func (k *Key) LoopSubKeys(f func(*Key, *KeyBaseInfo) error) (err error) {
	sub, err := k.Subkeys()
	if err != nil {
		return
	}
	for i := range sub {
		subk, err1 := k.Subkey(sub[i].Name)
		if err1 != nil {
			err = err1
			return
		}
		err = f(subk, &sub[i])
		subk.Close()
		if err != nil {
			break
		}
	}
	return
}

func (k *Key) Close() {
	syscall.RegCloseKey(k.Handle)
	runtime.SetFinalizer(k, nil)
}

type Value interface {
	Uint32() uint32
	String() string
	Data() []byte
}

type value struct {
	name []uint16
	typ  int
	sz   int
	*Key
}

func (k *Key) Values() (vmap map[string]Value, err error) {
	var n, maxNameLen, ulen, typ, sz uint32

	err = syscall.RegQueryInfoKey(k.Handle, nil, nil, nil, nil, nil, nil, &n, &maxNameLen, nil, nil, nil)
	if err != nil {
		return
	}

	maxNameLen++
	ubuf := make([]uint16, maxNameLen)
	vmap = make(map[string]Value, 8)

	for i := uint32(0); i < n; i++ {
		ulen = maxNameLen
		err = win.RegEnumValue(k.Handle, i, &ubuf[0], &ulen, nil, &typ, nil, &sz)
		if err != nil {
			break
		}
		s := syscall.UTF16ToString(ubuf[:ulen])
		vmap[s] = newValue(ubuf[:ulen+1], typ, sz, k)
	}
	return
}

func (k *Key) Value(name string) (v Value, err error) {
	var typ, sz uint32

	uname, err := syscall.UTF16FromString(name)
	if err != nil {
		return
	}
	err = syscall.RegQueryValueEx(k.Handle, &uname[0], nil, &typ, nil, &sz)
	if err == nil {
		v = newValue(uname, typ, sz, k)
	}
	return
}

func newValue(uname []uint16, typ, sz uint32, k *Key) Value {
	utf16 := make([]uint16, len(uname))
	copy(utf16, uname)
	v := value{utf16, int(typ), int(sz), k}
	var value Value
	switch typ {
	case syscall.REG_SZ:
		value = &String{v}
	case syscall.REG_DWORD_LITTLE_ENDIAN, syscall.REG_DWORD_BIG_ENDIAN:
		value = &Uint32{v}
	default:
		fallthrough
	case syscall.REG_BINARY:
		value = &Binary{v}
	}
	return value
}

type Binary struct {
	value
}

func (v *value) Data() (data []byte) {
	data = make([]byte, v.sz)
	sz := uint32(v.sz)
	if err := syscall.RegQueryValueEx(v.Handle, &v.name[0], nil, nil, &data[0], &sz); err == nil {
		data = data[:sz]
	} else {
		data = []byte{}
	}
	return
}
func (v *value) Uint32() uint32 {
	return 0
}
func (v *value) String() string {
	return string(v.Data())
}

type Uint32 struct {
	value
}

func (v *Uint32) Uint32() (u uint32) {
	data := v.Data()
	if len(data) < 4 {
		return
	}
	binary.Read(bytes.NewBuffer(data), byteOrder(v.typ == syscall.REG_DWORD_LITTLE_ENDIAN), &u)
	return u
}
func (v *Uint32) String() string {
	return fmt.Sprint(v.Uint32())
}

type Uint64 struct {
	value
}

type String struct {
	value
}

func (v *String) String() (s string) {
	data := v.Data()
	if len(data) < 2 {
		return
	}
	ubuf := make([]uint16, len(data)/2)
	binary.Read(bytes.NewBuffer(data), byteOrder(isLittleEndian), ubuf)
	return syscall.UTF16ToString(ubuf)
}

const (
	isLittleEndian = syscall.REG_DWORD == syscall.REG_DWORD_LITTLE_ENDIAN
)

func byteOrder(little bool) binary.ByteOrder {
	if little {
		return binary.LittleEndian
	}
	return binary.BigEndian
}
