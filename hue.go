package pixelsort

import (
	"math"
)

type hue struct {
	Region
}

func HueSorter(r Region) RegionSorter {
	return &hue{r}
}

func (h *hue) Less(i, j int) bool {
	r1i, g1i, b1i, _ := h.At(i).RGBA()
	r2i, g2i, b2i, _ := h.At(j).RGBA()

	// Normalize
	r1 := float64(r1i) / 255.0
	g1 := float64(g1i) / 255.0
	b1 := float64(b1i) / 255.0
	r2 := float64(r2i) / 255.0
	g2 := float64(g2i) / 255.0
	b2 := float64(b2i) / 255.0

	alpha1 := 0.5 * (2*r1 - g1 - b1)
	beta1 := math.Sqrt(3) / 2.0 * (g1 - b1)
	hue1 := math.Atan2(beta1, alpha1)

	alpha2 := 0.5 * (2*r2 - g2 - b2)
	beta2 := math.Sqrt(3) / 2.0 * (g2 - b2)
	hue2 := math.Atan2(beta2, alpha2)

	//c1 := math.Sqrt(alpha1*alpha1 + beta1*beta1)
	//c2 := math.Sqrt(alpha2*alpha2 + beta2*beta2)

	return hue1 < hue2
}

func (h *hue) Swap(i, j int) {
	tmp := h.At(i)
	h.Set(i, h.At(j))
	h.Set(j, tmp)
}

func (h *hue) Len() int {
	return h.Size()
}

