include $(GOROOT)/src/Make.inc 

ifeq ($(GOOS),windows)
REGISTRY=registry
endif

DIRS=\
	syscall\
	go9p\
	$(REGISTRY)\
	strings\
	sercom\
	\
	examples\

include Make.dirs
