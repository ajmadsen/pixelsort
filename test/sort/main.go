package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"sort"
)

type pixels []uint8

func (p pixels) Len() int { return len(p) / 4 * 3 }
func (p pixels) Less(i, j int) bool {
	iadj := i/3*4 + i%3
	jadj := j/3*4 + j%3
	return p[iadj] < p[jadj]
}
func (p pixels) Swap(i, j int) {
	iadj := i/3*4 + i%3
	jadj := j/3*4 + j%3
	p[iadj], p[jadj] = p[jadj], p[iadj]
	if counter += 1; counter%10000 == 0 {
		saveImage(sorting, fmt.Sprintf("%06d.png", counter/10000))
	}
}

var (
	sorting *image.RGBA
	counter = int(0)
)

func saveImage(m image.Image, name string) {
	log.Print("Opening output")
	outfile, err := os.Create(name)
	defer outfile.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Writing output")
	png.Encode(outfile, m)
}

func makeRGB(m image.Image) *image.RGBA {
	if rgb, ok := m.(*image.RGBA); ok {
		return rgb
	}
	b := m.Bounds()
	rgb := image.NewRGBA(b)
	draw.Draw(rgb, b, m, b.Min, draw.Src)
	return rgb
}

func sortimage(m image.Image) image.Image {
	sorting = makeRGB(m)
	sort.Sort(pixels(sorting.Pix))
	return sorting
}

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Printf("Usage: %s [flags] <input image> <output image>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	log.Print("Opening file")
	infile, err := os.Open(flag.Arg(0))
	defer infile.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Decoding image")
	inimg, _, err := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Sorting image")
	im := sortimage(inimg)

	saveImage(im, flag.Arg(1))
}
