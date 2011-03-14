#include <stdio.h>
#include <sys/types.h>
#include <pwd.h>
#include "user.h"
#include "_cgo_export.h"

// from plan9port, slightly adjusted
// http://swtch.com/usr/local/plan9/src/lib9/getuser.c

char*
GetUserName(int uid)
{
	struct passwd *pw;

	pw = getpwuid(uid);
	if (pw==NULL) {
		return NULL;
	}
	return  pw->pw_name;
}
