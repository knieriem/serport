## Utility packages for Go.

*	__sercom__

	Access to serial ports on Linux and Windows. A port may
	also be served and dialed to via [*9P*][9P].


*	__go9p__

	Utility functions and types for [go9p][]


*	__util__

	Contains `Tokenize`, an implementation of an
	equally named [function of Plan 9's libc][tokenize],
	which i miss from Go's string package (similar to
	`string.Fields`, but with interpretation of single
	quotes).

*	__registry__

	Access to Windows' registry database (still read-only). 


*	__syscall__

	System functions for Linux and Windows that were
	needed to implement the above packages.

	The make use of the `mksyscall*.sh` scripts from `$GOROOT/src/pkg/syscall`.

[9P]: http://plan9.bell-labs.com/sys/man/5/INDEX.html
[go9p]: http://code.google.com/p/go9p/
[hg-git]: http://hg-git.github.com/
[tokenize]: http://plan9.bell-labs.com/magic/man2html/2/getfields


## Installation

Since `goinstall` cannot cope with GOOS dependent source files yet, the following
commands can be used instead to install and build the packages:

	cd $GOROOT/src/pkg
	mkdir -p github.com/knieriem
	cd github.com/knieriem

Clone repository using Mercurial (utilizing the [hg-git][] extension):

	hg clone git://github.com/knieriem/g

... or using Git:

	git clone https://github.com/knieriem/g.git


Install prerequisites (Go 9P implementation):

	goinstall go9p.googlecode.com/hg/p
	goinstall go9p.googlecode.com/hg/p/clnt
	goinstall go9p.googlecode.com/hg/p/srv

Then,

	cd g
	make

Directory `examples' contains some programs making use of the packages.

I build everything on a 386 machine running Linux, also the windows
packages and binaries. A little `rc` script containing the lines

	#!/usr/local/plan9/bin/rc
	
	GOOS=windows
	GOBIN=$GOROOT/bin/$GOOS
	path=($GOBIN $path)
	prompt=(W-$prompt(1) $prompt(2))
	exec rc

helps switching between windows and host targets.


