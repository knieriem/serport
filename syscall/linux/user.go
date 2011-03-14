package syscall

// #include "user.h"
import "C"

import "os"

func GetUserName() string {
	s := C.GetUserName(C.int(os.Getuid()))
	if s == nil {
		return "none"
	}
	return C.GoString(s)
}
