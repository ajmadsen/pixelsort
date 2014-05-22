package pixelsort

import (
	"image/color"
	"math"
)

type hue struct {
	Region
	himg []hsv
}

func HueSorter(r Region) RegionSorter {
	l := r.Size()
	hs := &hue{r, make([]hsv, l)}
	for i := 0; i < l; i++ {
		hs.himg[i] = colorToHsv(r.At(i))
	}
	return hs
}

func (h *hue) Less(i, j int) bool {
	hsv1 := h.himg[i]
	hsv2 := h.himg[j]

	return hsv1.h*(hsv1.s+hsv1.v) < hsv2.h*(hsv2.s+hsv2.v)
}

func (h *hue) Swap(i, j int) {
	l := h.Size()
	h.swapInternal(i, j)

	// Also swap pixels to the left and right
	//if i-1 > 0 && j-1 > 0 {
	//	h.swapInternal(i-1, j-1)
	//}
	if i+1 < l && j+1 < l {
		h.swapInternal(i+1, j+1)
	}
	if i+2 < l && j+2 < l {
		h.swapInternal(i+2, j+2)
	}
	if i+3 < l && j+3 < l {
		h.swapInternal(i+3, j+3)
	}
	//if i+4 < l && j+4 < l {
	//	h.swapInternal(i+4, j+4)
	//}
}

func (h *hue) swapInternal(i, j int) {
	h.himg[i], h.himg[j] = h.himg[j], h.himg[i]
	bytes := h.Bytes()
	idx1 := h.Idx(i)
	idx2 := h.Idx(j)
	bytes[idx1], bytes[idx2] = bytes[idx2], bytes[idx1]
	bytes[idx1+1], bytes[idx2+1] = bytes[idx2+1], bytes[idx1+1]
	bytes[idx1+2], bytes[idx2+2] = bytes[idx2+2], bytes[idx1+2]
	bytes[idx1+3], bytes[idx2+3] = bytes[idx2+3], bytes[idx1+3]
}

func (h *hue) Len() int {
	return h.Size()
}

type hsv struct {
	h, s, v float64
}

// Inspired by http://lolengine.net/blog/2013/01/13/fast-rgb-to-hsv
func colorToHsv(c color.Color) hsv {
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

	return hsv{h, s, v}
}
