// +build ignore

package syscall

/*
#include "windows/c.h"
#include <stdint.h>

#undef OFF
#define OFF(ptr)	((char*)(ptr)-(char*)0)

#undef OffsetOf
#define OffsetOf(ptr) (OFF(ptr)<0? ((uint64_t)((unsigned int)~0) + 1 + OFF(ptr)): OFF(ptr))

#define HKEY_CR_OFF	OffsetOf(HKEY_CLASSES_ROOT)
#define HKEY_CC_OFF	OffsetOf(HKEY_CURRENT_CONFIG)
#define HKEY_CU_OFF	OffsetOf(HKEY_CURRENT_USER)
#define HKEY_LM_OFF	OffsetOf(HKEY_LOCAL_MACHINE)
*/
import "C"

type DCB C.DCB

type CommTimeouts C.COMMTIMEOUTS

const (
	dcbSize = C.sizeof_DCB

	HKEY_CLASSES_ROOT = C.HKEY_CR_OFF
	HKEY_CURRENT_CONFIG = C.HKEY_CC_OFF
	HKEY_CURRENT_USER = C.HKEY_CU_OFF
	HKEY_LOCAL_MACHINE = C.HKEY_LM_OFF
)

type REGSAM C.REGSAM
