package pixelsort

import (
	"image"
	"image/color"
	"sort"
)

type Region interface {
	// At returns the relative pixel in the region
	At(n int) color.Color

	// Size returns the maximum index in the region
	Size() int

	// Set sets the value of a pixel in the region
	Set(n int, c color.Color)

	sort.Interface
}

type rowRegion struct {
	image *image.RGBA
	row   int
	less  func(i, j color.Color) bool
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

func (i *rowRegion) Len() int {
	return i.Size()
}

func (in *rowRegion) Less(i, j int) bool {
	c1 := in.At(i)
	c2 := in.At(j)
	return in.less(c1, c2)
}

func (in *rowRegion) Swap(i, j int) {
	tmp := in.At(i)
	in.Set(i, in.At(j))
	in.Set(j, tmp)
}

type RegionEnum interface {
	Next() Region
}

type rowEnum struct {
	image *image.RGBA
	row   int
	less  func(i, j color.Color) bool
}

func (re *rowEnum) Next() Region {
	b := re.image.Bounds()
	if re.row+b.Min.Y >= b.Max.Y-1 {
		return nil
	}
	region := &rowRegion{
		image: re.image,
		row:   re.row,
		less:  re.less,
	}
	re.row++
	return region
}

func RowEnum(im *image.RGBA, less func(i, j color.Color) bool) RegionEnum {
	return &rowEnum{
		image: im,
		row:   0,
		less:  less,
	}
}
