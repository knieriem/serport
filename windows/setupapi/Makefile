GOOS=windows

PKG=setupapi


#include ../Make.zpkg
W=..

ZOSARCH=\
	windows_386\
	windows_amd64\

ZITEMS=\
	const\
	types\
	syscall\

ZFILES=\
	$(foreach f,$(ZITEMS),$(foreach a,$(ZOSARCH),z$(f)_$a.go))\

all:	$(ZFILES)

$(foreach f,$(ZITEMS),z$(f)_windows_%.go):	windows/c.h windows/types.go windows/const setupapi_windows.go $W/gen.sh
	GOOS=windows GOARCH=$* PKG=$(PKG) ZDIR=$W sh $W/gen.sh
