# Copyright 2011 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

GOOS=windows

include ../../Make.inc 

api=syscall

TARG=github.com/knieriem/g/$(api)

GOFILES=\
	$(api).go\
	z$(api)_$(GOARCH).go\
	ztypes_$(GOARCH).go\
	zconst_$(GOARCH).go\

include $(GOROOT)/src/Make.pkg

z%.go:
	GOARCH=386 sh ./gen.sh

zclean:
	rm -f z*.go


.PHONY: \
	zclean\
