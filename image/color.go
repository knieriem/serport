package image

import (
	"image/color"
)


// A BinaryColor represents either black or white.
type BinaryColor struct {
	Black bool
}

func (c BinaryColor) RGBA() (r, g, b, a uint32) {
	a = 0xffff
	if c.Black {
		return
	}
	return a, a, a, a
}

func toBinaryColor(c color.Color) color.Color {
	if _, ok := c.(BinaryColor); ok {
		return c
	}
	// should be some dithering
	r, g, b, _ := c.RGBA()
	return BinaryColor{(299*r+587*g+114*b+500)/1000 < 0x8000}
}

// The ColorModel associated with BinaryColor.
var BinaryColorModel color.Model = color.ModelFunc(toBinaryColor)
