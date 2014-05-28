package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
)

var (
	kernX = [][]int{
		[]int{-1, +0, +1},
		[]int{-2, +0, +2},
		[]int{-1, +0, +1},
	}
	kernY = [][]int{
		[]int{+1, +2, +1},
		[]int{+0, +0, +0},
		[]int{-1, -2, -1},
	}
	kernR = image.Rect(-1, -1, 2, 2)
)

func convolve(m *image.Gray, kernX, kernY [][]int, r image.Rectangle) uint8 {
	d := r.Sub(kernR.Min)
	b := m.Bounds()

	magx := int(0)
	magy := int(0)
	for y := d.Min.Y; y < d.Max.Y; y++ {
		for x := d.Min.X; x < d.Max.X; x++ {
			if (kernX[y][x] == 0) && (kernY[y][x] == 0) {
				continue
			}
			v := m.At(b.Min.X+x, b.Min.Y+y).(color.Gray)
			magx += int(v.Y) * kernX[y][x]
			magy += int(v.Y) * kernY[y][x]
		}
	}
	mag := math.Sqrt(float64(magx*magx) + float64(magy*magy))

	if mag < 0 {
		return 0
	}
	if mag > 255 {
		return 255
	}
	return uint8(mag)
}

func sobel(m *image.Gray) *image.Gray {
	b := m.Bounds()
	im := image.NewGray(b)

	for p := b.Min; p.Y < b.Max.Y; p.Y++ {
		for p.X = b.Min.X; p.X < b.Max.X; p.X++ {
			window := kernR.Add(p).Intersect(b)
			mag := math.MaxUint8 - convolve(m.SubImage(window).(*image.Gray), kernX, kernY, window.Sub(p))
			im.SetGray(p.X, p.Y, color.Gray{mag})
		}
	}

	return im
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <input image> <output image>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	log.Println("Opening image in")
	infile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	log.Println("Decoding image in")
	inimg, _, err := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}

	b := inimg.Bounds()
	gray := image.NewGray(b)
	draw.Draw(gray, b, inimg, b.Min, draw.Src)

	//g, err := os.Create("gray.png")
	//png.Encode(g, gray)

	log.Println("Computing pixel gradient")
	s := sobel(gray)

	log.Println("Writing output")
	outfile, err := os.Create(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(outfile, s)
	if err != nil {
		log.Fatal(err)
	}
}
