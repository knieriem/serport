// +build ignore

package syscall

/*
#include <termios.h>
#include <unistd.h>

#define AAA 0xF0000001

*/
import "C"

type Termios	C.struct_termios
type Int	C.int
const(
 Hoho = C.AAA
)

