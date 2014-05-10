package pixelsort

import (
	"image"
)

type Sorter interface {
	Sort()
	Image() image.Image
}

