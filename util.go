package pixelsort

import (
	"image"
	"image/color"
)

type Region interface {
	// At returns the relative pixel in the region
	At(n int) color.Color

	// Size returns the maximum index in the region
	Size() int

	// Set sets the value of a pixel in the region
	Set(n int, c color.Color)
}

type rowRegion struct {
	image *image.RGBA
	row   int
}

func (r rowRegion) At(n int) color.Color {
	b := r.image.Bounds()
	return r.image.At(b.Min.X+n, b.Min.Y+r.row)
}

func (r rowRegion) Size() int {
	return r.image.Bounds().Dx()
}

func (r rowRegion) Set(n int, c color.Color) {
	b := r.image.Bounds()
	r.image.Set(b.Min.X+n, b.Min.Y+r.row, c)
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
	if re.curRow+b.Min.Y >= b.Max.Y-1 {
		return nil
	}
	return &rowRegion{
		image: re.image,
		row:   re.curRow,
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
