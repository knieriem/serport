W=../win

OSARCH=\
	linux_arm\
	linux_386\
	linux_amd64\
	windows_386\

ITEMS=\
	const\
	types\
	syscall\

ZFILES=\
	$(foreach f,$(ITEMS),$(foreach a,$(OSARCH),z$(f)_$a.go))\

all:	$(ZFILES)

clean:
	rm -rf _obj	# created by cgo

zclean:
	rm -f z*.go

$(foreach f,$(ITEMS),z$(f)_windows_%.go):	windows/c.h windows/types.go windows/const syscall_windows.go $W/gen.sh
	GOOS=windows GOARCH=$* PKG=syscall ZDIR=$W sh $W/gen.sh

$(foreach f,$(ITEMS),z$(f)_linux_%.go):	linux/types.go linux/const syscall_linux.go linux/gen.sh
	GOOS=linux GOARCH=$* sh ./linux/gen.sh
