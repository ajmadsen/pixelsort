package pixelsort

import (
	"image"
	"image/draw"
	"sort"
)

type intensity struct {
	image *image.RGBA
	re    RegionEnum
}

func (in *intensity) Sort() {
	for i := in.re.Next(); i != nil; i = in.re.Next() {
		sort.Sort(i.(*rowRegion))
	}
}

func Intensity(im image.Image) Sorter {
	b := im.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)
	return &intensity{
		image: m,
		re:    RowEnum(m),
	}
}

// Implement Sort iterface on intensity
func (i *rowRegion) Len() int {
	return i.Size()
}

func (in *rowRegion) Less(i, j int) bool {
	c1 := in.At(i)
	c2 := in.At(j)

	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return a1*r1*r1+a1*g1*g1+a1*b1*b1 < a2*r2*r2+a2*g2*g2+a2*b2*b2
}

func (in *rowRegion) Swap(i, j int) {
	tmp := in.At(i)
	in.Set(i, in.At(j))
	in.Set(j, tmp)
}
