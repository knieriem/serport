// +build ignore

package syscall

/*
#include "windows/c.h"
*/
import "C"

type DCB C.DCB

type CommTimeouts C.COMMTIMEOUTS

const (
	dcbSize = C.sizeof_DCB
)
