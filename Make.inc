include $(GOROOT)/src/Make.inc

all: install

gofmt:
	gofmt -w -s *.go

doc:
	godoc -html . > doc.html

CLEANFILES += doc.html
