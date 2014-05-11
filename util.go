package pixelsort

import (
	"image"
	"image/color"
)

type Region interface {
	// At returns the relative pixel in the region
	At(n int) color.Color

	// Idx returns the first element in bytes
	Idx(n int) int

	// Size returns the maximum index in the region
	Size() int

	// Set sets the value of a pixel in the region
	Set(n int, c color.Color)

	// Bytes returns the slice of bytes the region represents
	Bytes() []uint8
}

type rowRegion struct {
	image *image.RGBA
}

func (r *rowRegion) At(n int) color.Color {
	b := r.image.Bounds()
	return r.image.At(b.Min.X+n, b.Min.Y)
}

func (r *rowRegion) Idx(n int) int {
	b := r.image.Bounds()
	return r.image.PixOffset(b.Min.X+n, b.Min.Y)
}

func (r *rowRegion) Size() int {
	return r.image.Bounds().Dx()
}

func (r *rowRegion) Set(n int, c color.Color) {
	b := r.image.Bounds()
	r.image.Set(b.Min.X+n, b.Min.Y, c)
}

func (r *rowRegion) Bytes() []uint8 {
	return r.image.Pix
}

type RegionEnum interface {
	Next() Region
	Value() Region
	Reset()
}

type rowEnum struct {
	image  *image.RGBA
	curRow int
}

func (re *rowEnum) Next() Region {
	re.curRow++
	return re.Value()
}

func (re *rowEnum) Value() Region {
	b := re.image.Bounds()
	low := b.Min.Add(image.Pt(0, re.curRow))
	high := b.Min.Add(image.Pt(b.Dx(), re.curRow+1))
	sub := image.Rectangle{low, high}
	if !sub.In(b) {
		return nil
	}
	return &rowRegion{
		image: re.image.SubImage(sub).(*image.RGBA),
	}
}

func (re *rowEnum) Reset() {
	re.curRow = 0
}

func RowEnum(im *image.RGBA) RegionEnum {
	return &rowEnum{
		image:  im,
		curRow: 0,
	}
}
