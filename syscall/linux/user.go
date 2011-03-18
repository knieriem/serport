package syscall

// #include <unistd.h>
// #include <pwd.h>
import "C"


func GetUserName() string {
	pw := C.getpwuid(C.getuid())
	if pw == nil {
		return "none"
	}
	return C.GoString(pw.pw_name)
}
