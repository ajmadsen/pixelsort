package pixelsort

import (
	"image"
	"image/draw"
	"image/color"
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

func (in *intensity) Image() image.Image {
	return in.image
}

func Intensity(im image.Image) Sorter {
	b := im.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)
	return &intensity{
		image: m,
		re:    RowEnum(m, intensitySort),
	}
}

// Implement Sort iterface on intensity
func intensitySort(i, j color.Color) bool {
	r1, g1, b1, _ := i.RGBA()
	r2, g2, b2, _ := j.RGBA()
	return r1*r1+g1*g1+b1*b1 < r2*r2+g2*g2+b2*b2
}

