package pixelsort

type intensity struct {
	Region
}

func IntensitySorter(r Region) RegionSorter {
	return &intensity{r}
}

func (h *intensity) Less(i, j int) bool {
	r1, g1, b1, _ := h.At(i).RGBA()
	r2, g2, b2, _ := h.At(j).RGBA()
	return r1*r1+g1*g1+b1*b1 < r2*r2+g2*g2+b2*b2
}

func (h *intensity) Swap(i, j int) {
	tmp := h.At(i)
	h.Set(i, h.At(j))
	h.Set(j, tmp)
}

func (h *intensity) Len() int {
	return h.Size()
}

