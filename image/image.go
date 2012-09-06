// Package image implements a bitmap image type, and a binary color type
package image

import (
	i "image"
	"image/color"
)

// A Bitmap is an in-memory image of BinaryColor values.
type Bitmap struct {
	// Pix holds the image's pixels. The pixel at (x, y) is Pix[y*Stride + x/8] >> (7 - x&7) & 1.
	Pix    []uint8
	Stride int
	// Rect is the image's bounds.
	Rect i.Rectangle
}

const (
	wShift = 3
	wS     = 1 << wShift
	wSmod  = wS - 1
)

var bits = []uint8{1 << 7, 1 << 6, 1 << 5, 1 << 4, 1 << 3, 1 << 2, 1 << 1, 1 << 0}

func (p *Bitmap) bitAddr(x, y int) (addr *uint8, bit uint8) {
	r := &p.Rect
	addr = &p.Pix[p.Stride*(y-r.Min.Y)+(x>>wShift)-(r.Min.X>>wShift)]
	bit = bits[x&wSmod]
	return

}

func (p *Bitmap) ColorModel() color.Model { return BinaryColorModel }

func (p *Bitmap) Bounds() i.Rectangle { return p.Rect }

func (p *Bitmap) At(x, y int) color.Color {
	if !(i.Point{x, y}.In(p.Rect)) {
		return BinaryColor{}
	}
	addr, bit := p.bitAddr(x, y)
	return BinaryColor{(*addr)&bit != 0}
}

func (p *Bitmap) Set(x, y int, c color.Color) {
	if !(i.Point{x, y}.In(p.Rect)) {
		return
	}
	addr, bit := p.bitAddr(x, y)
	if toBinaryColor(c).(BinaryColor).Black {
		*addr |= bit
	} else {
		*addr &^= bit
	}
}

func (p *Bitmap) SetBinary(x, y int, c BinaryColor) {
	if !(i.Point{x, y}.In(p.Rect)) {
		return
	}

	addr, bit := p.bitAddr(x, y)
	if c.Black {
		*addr |= bit
	} else {
		*addr &^= bit
	}
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *Bitmap) SubImage(r i.Rectangle) i.Image {
	return &Bitmap{
		Pix:    p.Pix,
		Stride: p.Stride,
		Rect:   p.Rect.Intersect(r),
	}
}

// Opaque scans the entire image and returns whether or not it is fully opaque.
func (p *Bitmap) Opaque() bool {
	return true
}

// NewBitmap returns a new Bitmap with the given width and height.
func NewBitmap(w, h int) *Bitmap {
	bytesPerRow := (w-1)/8 + 1
	pix := make([]uint8, h*bytesPerRow)
	return &Bitmap{pix, bytesPerRow, i.Rectangle{i.ZP, i.Point{w, h}}}
}
