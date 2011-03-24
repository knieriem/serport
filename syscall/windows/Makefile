# Copyright 2011 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

GOOS=windows

include ../../Make.inc 

PKG=syscall

TARG=github.com/knieriem/g/$(PKG)

GOFILES=\
	$(PKG).go\
	z$(PKG)_$(GOARCH).go\
	ztypes_$(GOARCH).go\
	zconst_$(GOARCH).go\

ZDIR=../../win
include $(ZDIR)/Make.zpkg
