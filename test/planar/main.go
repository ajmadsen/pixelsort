package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
	"os"
)

var (
	inFileName  string
	outFileName string
)

func init() {
	flag.StringVar(&inFileName, "in", "", "input file name")
	flag.StringVar(&outFileName, "out", "", "output file name")
}

func main() {
	infile, err := os.Open(inFileName)
	if err != nil {
		panic(err)
	}

	input, _, err := image.Decode(infile)
	if err != nil {
		panic(err)
	}

	log.Println("converting image")
	inputRGBA := image.NewNRGBA(input.Bounds())
	draw.Draw(inputRGBA, inputRGBA.Bounds(), input, image.Pt(0, 0), draw.Src)
}

func sortPlane(im *image.NRGBA, plane int, maxCmp int) {
	b := im.Bounds()
	cmp := 0
	pt := b.Min
	for ; pt.X >= 0; pt = nextPx(pt, b) {
		for npt := nextPx(pt, b); npt.X >= 0; npt = nextPx(npt, b) {
			cmp++
			if cmpPlanePx(pt, npt, im, plane) {
				swpPlanePx(pt, npt, im, plane)
			}
			if cmp > maxCmp {
				return
			}
		}
	}
}

func getPlanePx(pt image.Point, im *image.NRGBA, plane int) uint8 {
	switch plane {
	case 0:
		return im.NRGBAAt(pt.X, pt.Y).R
	case 1:
		return im.NRGBAAt(pt.X, pt.Y).G
	case 2:
		return im.NRGBAAt(pt.X, pt.Y).B
	case 3:
		return im.NRGBAAt(pt.X, pt.Y).A
	}
	return 0
}

func cmpPlanePx(p1, p2 image.Point, im *image.NRGBA, plane int) bool {
	return getPlanePx(p1, im, plane) < getPlanePx(p2, im, plane)
}

func swpPlanePx(pt1, pt2 image.Point, im *image.NRGBA, plane int) {
	c1 := im.NRGBAAt(pt1.X, pt1.Y)
	c2 := im.NRGBAAt(pt2.X, pt2.Y)
	switch plane {
	case 0:
		tmp := c1.R
		c1.R = c2.R
		c2.R = tmp
	case 1:
		tmp := c1.G
		c1.G = c2.G
		c2.G = tmp
	case 2:
		tmp := c1.B
		c1.B = c2.B
		c2.B = tmp
	case 3:
		tmp := c1.A
		c1.A = c2.A
		c2.A = tmp
	}
	im.SetNRGBA(pt1.X, pt1.Y, c1)
	im.SetNRGBA(pt2.X, pt2.Y, c2)
}

func nextPx(pt image.Point, b image.Rectangle) (npt image.Point) {
	// column first, then row
	npt = pt

	npt.Y++
	if npt.Y > b.Max.Y {
		npt.Y = b.Min.Y
		npt.X++
		if npt.X > b.Max.X {
			npt.Y = -1
			npt.X = -1
		}
	}

	return
}
