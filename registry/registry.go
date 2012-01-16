package registry

import (
	"bytes"
	"encoding/binary"
	"fmt"

	win "github.com/knieriem/g/syscall"
	"runtime"
	"syscall"
)

type Key struct {
	win.HKEY
}

var (
	KeyClassesRoot   = &Key{win.HKEY_CLASSES_ROOT}
	KeyLocalMachine  = &Key{win.HKEY_LOCAL_MACHINE}
	KeyCurrentUser   = &Key{win.HKEY_CURRENT_USER}
	KeyCurrentConfig = &Key{win.HKEY_CURRENT_CONFIG}
)

func (k *Key) Subkey(subkey ...string) (result *Key, err error) {
	var key win.HKEY

	s := ""
	for i, v := range subkey {
		if i > 0 {
			s += "\\" + v
		} else {
			s += v
		}
	}
	if e := win.RegOpenKeyEx(k.HKEY, syscall.StringToUTF16Ptr(s), 0, win.KEY_READ, &key); e != 0 {
		err = syscall.Errno(e)
	} else {
		result = &Key{key}
		runtime.SetFinalizer(result, (*Key).Close)
	}
	return
}

func (k *Key) Close() {
	win.RegCloseKey(k.HKEY)
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

func (k *Key) Values() (vmap map[string]Value) {
	var ulen, typ, sz uint32

	ubuf := make([]uint16, 256)
	vmap = make(map[string]Value, 8)

	for i := 0; ; i++ {
		ulen = uint32(len(ubuf))
		e := win.RegEnumValue(k.HKEY, uint32(i), &ubuf[0], &ulen, nil, &typ, nil, &sz)
		if e != 0 {
			break
		}
		s := syscall.UTF16ToString(ubuf[:ulen])
		utf16 := make([]uint16, ulen+1)
		copy(utf16, ubuf[:len(utf16)])
		v := value{utf16, int(typ), int(sz), k}
		var value Value
		switch typ {
		case win.REG_SZ:
			value = &String{v}
		case win.REG_DWORD_LITTLE_ENDIAN, win.REG_DWORD_BIG_ENDIAN:
			value = &Uint32{v}
		default:
			fallthrough
		case win.REG_BINARY:
			value = &Binary{v}
		}
		vmap[s] = value
	}
	return vmap
}

type Binary struct {
	value
}

func (v *value) Data() (data []byte) {
	data = make([]byte, v.sz)
	sz := uint32(v.sz)
	if e := win.RegQueryValueEx(v.HKEY, &v.name[0], nil, nil, &data[0], &sz); e == 0 {
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
	return fmt.Sprint(v.Data())
}

type Uint32 struct {
	value
}

func (v *Uint32) Uint32() (u uint32) {
	data := v.Data()
	if len(data) < 4 {
		return
	}
	binary.Read(bytes.NewBuffer(data), byteOrder(v.typ == win.REG_DWORD_LITTLE_ENDIAN), &u)
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
	isLittleEndian = win.REG_DWORD == win.REG_DWORD_LITTLE_ENDIAN
)

func byteOrder(little bool) binary.ByteOrder {
	if little {
		return binary.LittleEndian
	}
	return binary.BigEndian
}
