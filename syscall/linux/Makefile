GOOS=linux

include ../../Make.inc 

TARG=github.com/knieriem/g/syscall

GOFILES=\
	syscall.go\
	zsyscall_$(GOARCH).go\
	ztypes_$(GOARCH).go\
	zconst_$(GOARCH).go\

CGOFILES=\
	user.go\

include $(GOROOT)/src/Make.pkg

z%.go:
	GOARCH=386 sh ./gen.sh
	GOARCH=amd64 sh ./gen.sh

zclean:
	rm -f z*.go


.PHONY: \
	zclean\

