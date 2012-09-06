// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Derived from image/gif/reader.go, changes
// © 2011 M. Teichgräber

// Package pnm implements a PBM image decoder.
package pnm

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
)

const (
	maxAsciiNum = 0xfffff
)

var formats = []*format{&pbmRawFormat}

type format struct {
	id         byte
	readMaxval bool
	colorModel color.Model
	decode     func(r io.Reader, w, h, maxval int) image.Image
}

// If the io.Reader does not also have (Un)ReadByte, then decode will introduce its own buffering.
type reader interface {
	io.Reader
	io.ByteScanner
}

// decoder is the type used to decode a PNM file.
type decoder struct {
	r reader

	// From header.
	width  int
	height int
	maxval int

	// Computed.
	f *format

	image image.Image
}

// decode reads a PNM image from r and stores the result in d.
func (d *decoder) decode(r io.Reader, configOnly bool) (err error) {

	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok {
				err = errors.New("pnm:" + s)
			} else {
				err = r.(error)
			}
		}
	}()

	// Add buffering if r does not provide ReadByte.
	if rr, ok := r.(reader); ok {
		d.r = rr
	} else {
		d.r = bufio.NewReader(r)
	}

	d.readHeader()
	if configOnly {
		return
	}

	d.image = d.f.decode(d.r, d.width, d.height, d.maxval)

	return
}

func (d *decoder) readHeader() {
	var f *format

	if d.readByte() != 'P' {
		panic("magic number not recognized")
	}
	format := d.readByte()
	for i := range formats {
		f = formats[i]
		if f.id == format {
			goto readDim
		}
	}
	panic("unsupported format")

readDim:
	d.f = f
	d.width = d.btoi(maxAsciiNum)
	d.height = d.btoi(maxAsciiNum)
	if f.readMaxval {
		d.maxval = d.btoi(0xffff)
	}
	switch d.readByte() {
	case '\r', '\n', ' ', '\t':
		return
	}
	panic("header: expected single white space character")
}

// skip PBM whitespace: SPC, TAB, CR, LF, and comments
//
func (d *decoder) skipWhite() {
	var skipLine = false

loop:
	for {
		b := d.readByte()
		switch b {
		case '#':
			skipLine = true
		case '\n':
			skipLine = false
		case ' ', '\r', '\t':
		default:
			if !skipLine {
				break loop
			}
		}
	}
	d.unreadByte()
}

// read an ASCII decimal number
//
func (d *decoder) btoi(max int) (v int) {
	var start = true

	d.skipWhite()
	for {
		b := d.readByte()
		if b >= '0' && b <= '9' {
			v *= 10
			v += int(b - '0')
			if v > max {
				panic("header: number too big")
			}
			start = false
		} else {
			if start {
				panic("header: expected number")
			}
			break
		}
	}
	d.unreadByte()
	return
}

func (d *decoder) readByte() (b uint8) {
	b, err := d.r.ReadByte()
	if err != nil {
		panic(err)
	}
	return
}

func (d *decoder) unreadByte() {
	if err := d.r.UnreadByte(); err != nil {
		panic(err)
	}
}

// Decode reads a PBM image from r and returns the first embedded
// image as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	var d decoder
	if err := d.decode(r, false); err != nil {
		return nil, err
	}
	return d.image, nil
}

// DecodeConfig returns the color model and dimensions of a PBM image
// without decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	var d decoder
	if err := d.decode(r, true); err != nil {
		return image.Config{}, err
	}
	return image.Config{d.f.colorModel, d.width, d.height}, nil
}
