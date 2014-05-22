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
	l := h.Size()
	h.swapInternal(i, j)

	// Also swap pixels to the left and right
	if i-1 > 0 && j-1 > 0 {
		h.swapInternal(i-1, j-1)
	}
	if i+1 < l && j+1 < l {
		h.swapInternal(i+1, j+1)
	}
}

func (h *intensity) swapInternal(i, j int) {
	bytes := h.Bytes()
	idx1 := h.Idx(i)
	idx2 := h.Idx(j)
	bytes[idx1], bytes[idx2] = bytes[idx2], bytes[idx1]
	bytes[idx1+1], bytes[idx2+1] = bytes[idx2+1], bytes[idx1+1]
	bytes[idx1+2], bytes[idx2+2] = bytes[idx2+2], bytes[idx1+2]
	bytes[idx1+3], bytes[idx2+3] = bytes[idx2+3], bytes[idx1+3]
}

func (h *intensity) Len() int {
	return h.Size()
}
