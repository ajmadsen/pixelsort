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

func (r *rowRegion) At(n int) color.Color {
	b := r.image.Bounds()
	return r.image.At(b.Min.X+n, b.Min.Y+r.row)
}

func (r *rowRegion) Size() int {
	return r.image.Bounds().Dx()
}

func (r *rowRegion) Set(n int, c color.Color) {
	b := r.image.Bounds()
	r.image.Set(b.Min.X+n, b.Min.Y+r.row, c)
}

type RegionEnum interface {
	Next() Region
}

type rowEnum struct {
	image *image.RGBA
	row   int
}

func (re *rowEnum) Next() Region {
	b := re.image.Bounds()
	if re.row + b.Min.Y >= b.Max.Y-1 {
		return nil
	}
	region := &rowRegion{
		image: re.image,
		row:   re.row,
	}
	re.row++
	return region
}

func RowEnum(im *image.RGBA) RegionEnum {
	return &rowEnum{
		image: im,
		row:   0,
	}
}
