package pnm

import (
	"image"
	img "github.com/knieriem/g/image"
	"io"
)

func decodePbmRaw(r io.Reader, width, height, _ int) image.Image {
	im := img.NewBitmap(width, height)
	if _, err := io.ReadFull(r, im.Pix); err != nil {
		panic(err)
	}
	return im
}

var pbmRawFormat = format{'4', false, img.BinaryColorModel, decodePbmRaw}

func init() {
	image.RegisterFormat("pbm", "P4", Decode, DecodeConfig)
}
