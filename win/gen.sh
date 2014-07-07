# depends on environment variables:
#	PKG, ZDIR, GOARCH

set -e

pkg=$PKG
OS=$GOOS
GOROOT=`go env GOROOT`

case $GOARCH in
386)
	gccarch=i686
	;;
amd64)
	gccarch=x86_64
	;;
*)
	echo GOARCH $GOARCH not supported
	exit 1
	;;
esac

GCC=/usr/bin/$gccarch-w64-mingw32-gcc

SFX=_${OS}_$GOARCH.go

src=${pkg}_$OS.go
mv $src _$src
sed '/^package/s,syscall,none,' <_$src >$src
GOARCH= GOOS= go build $mksyscall $GOROOT/src/pkg/syscall/mksyscall_windows.go
./mksyscall_windows $src |
	sed '/import.*unsafe/a\
import "syscall"' |
	sed 's,EINVAL,syscall.EINVAL,g;s,Syscall,syscall.Syscall,;s,NewLazyDLL,syscall.NewLazyDLL,;s/^package.*syscall/package '$pkg'/' |
	gofmt > z$pkg$SFX
rm -f mksyscall_windows
rm -f $src
mv _$src $src

if test -f $OS/types.go; then
	# note: cgo execution depends on $GOARCH value
	GCC=$GCC go tool cgo -godefs $OS/types.go  |
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

GCC=$GCC go tool cgo -godefs ,,const.go |
	awk -f $ZDIR/fixtype.awk |
	gofmt > zconst$SFX
rm -f ,,const.go
