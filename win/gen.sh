# depends on environment variables:
#	PKG, ZDIR, GOARCH

pkg=$PKG
OS=$GOOS

mksyscall=$GOROOT/src/pkg/syscall/mksyscall_windows.pl

case $GOARCH in
386)
	gccarch=i586
	arch=-l32
	;;
amd64)
	gccarch=amd64
	arch=
	;;
*)
	echo GOARCH $GOARCH not supported
	exit 1
	;;
esac

GCC=/usr/bin/$gccarch-mingw32msvc-gcc

SFX=_${OS}_$GOARCH.go

perl $mksyscall $arch ${pkg}_$OS.go |
	sed 's/^package.*syscall$/package '$PKG'/' |
	sed '/^import/a \
		import "syscall"' |
	sed '/import *"DISABLEDunsafe"/d' |
	sed 's/Syscall/syscall.Syscall/' |
	sed 's/NewLazyDLL/syscall.&/' |
	sed 's/EINVAL/syscall.EINVAL/' |
	gofmt > z$pkg$SFX

if test -f $OS/types.go; then
	# note: cgo execution depends on $GOARCH value
	GCC=$GCC cgo -godefs $OS/types.go  |
		sed '/Pad_cgo_0/c\
		Flags	uint32' |
		awk -f $ZDIR/fixtype.awk |
		gofmt >ztypes$SFX
fi

if test -f $OS/const; then :
else
	exit 0
fi

(
	cat <<EOF
package $pkg

/*
#include "$OS/c.h"
*/
import "C"

const (
EOF
	<$OS/const awk '
		/^[^\/]/ { print "\t" $1 "= C." $1 ; next}
		{ print }
	'
	echo ')'
) > ,,const.go

GCC=$GCC cgo -godefs ,,const.go |
	awk -f $ZDIR/fixtype.awk |
	gofmt > zconst$SFX
rm -f ,,const.go
