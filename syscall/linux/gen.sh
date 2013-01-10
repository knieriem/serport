pkg=syscall
OS=$GOOS
ARCH=$GOARCH

mksyscall=$GOROOT/src/pkg/syscall/mksyscall.pl

perl $mksyscall ${pkg}_$OS.go |
	sed 's/^package.*syscall$$/package $*/' |
	sed '/^import/a \
		import "syscall"' |
	sed 's/Syscall/syscall.Syscall/' |
	sed 's/SYS_/syscall.SYS_/' |
	gofmt > z${pkg}_${OS}_$ARCH.go

# note: cgo execution depends on $GOARCH value
go tool cgo -godefs $OS/types.go  |
	gofmt >ztypes_${OS}_$ARCH.go


(
	cat <<EOF
package $pkg
/*
#include <unistd.h>
#include <termios.h>
#include <sys/ioctl.h>
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

go tool cgo -godefs ,,const.go | gofmt > zconst_${OS}_$ARCH.go
rm -f ,,const.go
