# depends on environment variables:
#	PKG, ZDIR, GOARCH

pkg=$PKG

CC=/usr/bin/i586-mingw32msvc-gcc
mksyscall=$GOROOT/src/pkg/syscall/mksyscall_windows.pl

case $GOARCH in
386)
	f=-m32
	;;
amd64)
	f=-m64
	;;
*)
	echo GOARCH $GOARCH not supported
	exit 1
	;;
esac

SFX=_$GOARCH.go

perl $mksyscall $pkg.go |
	sed 's/^package.*syscall$/package '$PKG'/' |
	sed '/^import/a \
		import "syscall"' |
	sed '/import *"DISABLEDunsafe"/d' |
	sed 's/Syscall/syscall.Syscall/' |
	sed 's/NewLazyDLL/syscall.&/' |
	sed 's/EINVAL/syscall.EINVAL/' |
	gofmt > z$pkg$SFX


godefs -g $pkg -f$f -c $CC types.c  |
	sed '/Pad_godefs_0/c\
	Flags	uint32' |
	awk -f $ZDIR/fixtype.awk |
	gofmt >ztypes$SFX


(
	echo '#include "c.h"'
	echo 'enum {'
	sed '/^[^/]/ s/.*/	$& = &,/' < const
	echo '};'
) > ,,const.c

godefs -g $pkg -c $CC ,,const.c |
	awk -f $ZDIR/fixtype.awk |
	gofmt > zconst$SFX
rm -f ,,const.c
