package pixelsort

type intensity struct {
	Region
}

func IntensitySorter(r Region) RegionSorter {
	return &intensity{r}
}

func (h *intensity) Less(i, j int) bool {
	r1, g1, b1, _ := h.At(i).RGBA()
	r2, g2, b2, _ := h.At(i).RGBA()
	return r1*r1+g1*g1+b1*b1 < r2*r2+g2*g2+b2*b2
}

func (h *intensity) Swap(i, j int) {
	bytes := h.Bytes()
	idx1 := h.Idx(i)
	idx2 := h.Idx(j)
	bytes[idx1], bytes[idx2] = bytes[idx2], bytes[idx1]
	bytes[idx1 + 1], bytes[idx2 + 1] = bytes[idx2 + 1], bytes[idx1 + 1]
	bytes[idx1 + 2], bytes[idx2 + 2] = bytes[idx2 + 2], bytes[idx1 + 2]
	bytes[idx1 + 3], bytes[idx2 + 3] = bytes[idx2 + 3], bytes[idx1 + 3]
}

func (h *intensity) Len() int {
	return h.Size()
}

