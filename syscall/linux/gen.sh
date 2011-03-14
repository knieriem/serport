pkg=syscall

mksyscall=$GOROOT/src/pkg/syscall/mksyscall.sh

ARCH=$GOARCH

case $ARCH in
386)
	f=-m32
	;;
amd64)
	f=-m64
	;;
esac


$mksyscall $pkg.go |
	sed 's/^package.*syscall$$/package $*/' |
	sed '/^import/a \
		import "syscall"' |
	sed 's/Syscall/syscall.Syscall/' |
	sed 's/SYS_/syscall.SYS_/' |
	gofmt > z${pkg}_$ARCH.go

godefs -g $pkg -f$f types.c  |
	gofmt >ztypes_$ARCH.go


(
	cat <<EOF
#include <unistd.h>
#include <termios.h>
#include <sys/ioctl.h>
enum {
EOF
	<const awk '
		/^[^\/]/ { print "$" $1 "= " $1 "," ; next}
		{ print }
	'
	echo '};'
) > ,,const.c

godefs -g $pkg -f$f ,,const.c | gofmt > zconst_$ARCH.go
rm -f ,,const.c
