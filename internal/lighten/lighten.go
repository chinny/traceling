// Package lighten converts rendered page images into light-grey tracing images.
package lighten

import (
	"image"
	"image/color"
)

// ToTraceable converts src to 8-bit grayscale and remaps its tonal range so
// that black becomes the given grey level while white stays white. Values in
// between scale linearly: out = grey + (255-grey)*in/255.
func ToTraceable(src image.Image, grey uint8) *image.Gray {
	lut := buildLUT(grey)
	b := src.Bounds()
	dst := image.NewGray(image.Rect(0, 0, b.Dx(), b.Dy()))

	if rgba, ok := src.(*image.RGBA); ok {
		for y := 0; y < b.Dy(); y++ {
			srow := rgba.Pix[y*rgba.Stride : y*rgba.Stride+b.Dx()*4]
			drow := dst.Pix[y*dst.Stride : y*dst.Stride+b.Dx()]
			for x := 0; x < b.Dx(); x++ {
				drow[x] = lut[luminance(srow[x*4], srow[x*4+1], srow[x*4+2])]
			}
		}
		return dst
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			g := color.GrayModel.Convert(src.At(x, y)).(color.Gray)
			dst.SetGray(x-b.Min.X, y-b.Min.Y, color.Gray{Y: lut[g.Y]})
		}
	}
	return dst
}

func buildLUT(grey uint8) [256]uint8 {
	var lut [256]uint8
	for v := 0; v < 256; v++ {
		lut[v] = grey + uint8((255-int(grey))*v/255)
	}
	return lut
}

// luminance is the ITU-R BT.601 weighted average used by color.GrayModel.
func luminance(r, g, b uint8) uint8 {
	return uint8((19595*uint32(r) + 38470*uint32(g) + 7471*uint32(b) + 1<<15) >> 16)
}
