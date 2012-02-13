// +build ignore

package syscall

/*
#include <termios.h>
#include <unistd.h>
*/
import "C"

type Termios	C.struct_termios
type Int	C.int
