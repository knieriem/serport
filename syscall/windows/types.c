#include <windows.h>

typedef DCB $DCB;

typedef COMMTIMEOUTS $CommTimeouts;

enum {
	$dcbSize = sizeof(DCB),
};



#undef PtrToInt
#define PtrToInt(name) $##name = (char*)(name)-(char*)0,

enum {
	PtrToInt(HKEY_CLASSES_ROOT)
	PtrToInt(HKEY_CURRENT_CONFIG)
	PtrToInt(HKEY_CURRENT_USER)
	PtrToInt(HKEY_LOCAL_MACHINE)
};

typedef REGSAM $REGSAM;
