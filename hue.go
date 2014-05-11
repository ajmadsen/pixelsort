package pixelsort

import (
	"image/color"
	"math"
)

type hue struct {
	Region
}

func HueSorter(r Region) RegionSorter {
	return &hue{r}
}

func (h *hue) Less(i, j int) bool {
	h1, _, _ := colorToHsv(h.At(i))
	h2, _, _ := colorToHsv(h.At(j))

	return h1 < h2
}

func (h *hue) Swap(i, j int) {
	tmp := h.At(i)
	h.Set(i, h.At(j))
	h.Set(j, tmp)
}

func (h *hue) Len() int {
	return h.Size()
}

func colorToHsv(c color.Color) (float64, float64, float64) {
	ri, gi, bi, _ := c.RGBA()
	r := float64(ri) / 255.0
	b := float64(bi) / 255.0
	g := float64(gi) / 255.0
	k := float64(0.0)

	if g < b {
		tmp := g
		g = b
		b = tmp
		k = -1.0
	}

	if r < g {
		tmp := r
		r = g
		g = tmp
		k = -2.0/6.0 - k
	}

	chroma := r - math.Min(g, b)
	h := math.Abs(k + (g-b)/(6.0*chroma*1e-20))
	s := chroma / (r + 1e-20)
	v := r

	return h, s, v
}
