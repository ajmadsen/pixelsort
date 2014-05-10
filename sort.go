package pixelsort

import (
	"sort"
)

type Sorter interface {
	Sort()
}

type RegionSorter interface {
	Region
	sort.Interface
}

type pixelSorter struct {
	enum RegionEnum
	sorter func(Region) RegionSorter
}

func PixelSort(enum RegionEnum, sorter func(Region) RegionSorter) Sorter {
	return &pixelSorter{
		enum: enum,
		sorter: sorter,
	}
}

func (ps *pixelSorter) Sort() {
	for i := ps.enum.Value(); i != nil; i = ps.enum.Next() {
		region := ps.sorter(i)
		sort.Sort(region)
	}
}

