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
	"os"
	"sort"
)

type pixels []uint8
type pixim image.RGBA

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

func (p *pixim) Len() int {
	b := (*image.RGBA)(p).Bounds()
	return b.Dx() * b.Dy()
}

func (p *pixim) Less(i, j int) bool {
	m := (*image.RGBA)(p)
	b := m.Bounds()
	iy := i / b.Dx()
	ix := i % b.Dx()
	jy := j / b.Dx()
	jx := j % b.Dx()
	c1, c2 := m.At(ix, iy).(color.RGBA), m.At(jx, jy).(color.RGBA)
	dist1 := uint32(c1.R)*uint32(c1.R) + uint32(c1.G)*uint32(c1.G) + uint32(c1.B)*uint32(c1.B)
	dist2 := uint32(c2.R)*uint32(c2.R) + uint32(c2.G)*uint32(c2.G) + uint32(c2.B)*uint32(c2.B)
	return dist1 < dist2
}

func (p *pixim) Swap(i, j int) {
	m := (*image.RGBA)(p)
	b := m.Bounds()
	iy := i / b.Dx()
	ix := i % b.Dx()
	jy := j / b.Dx()
	jx := j % b.Dx()
	pi := m.PixOffset(ix, iy)
	pj := m.PixOffset(jx, jy)
	m.Pix[pi+0], m.Pix[pj+0] = m.Pix[pj+0], m.Pix[pi+0]
	m.Pix[pi+1], m.Pix[pj+1] = m.Pix[pj+1], m.Pix[pi+1]
	m.Pix[pi+2], m.Pix[pj+2] = m.Pix[pj+2], m.Pix[pi+2]
	m.Pix[pi+3], m.Pix[pj+3] = m.Pix[pj+3], m.Pix[pi+3]
	if counter += 1; counter%1000 == 0 {
		saveImage(sorting, fmt.Sprintf("%06d.png", counter/1000))
	}
}

var (
	sorting *image.RGBA
	counter = int(0)
	method  = flag.Int("m", 0, "the method to use (0=planar, 1=pixel)")
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

func sortimage0(m image.Image) image.Image {
	sorting = makeRGB(m)
	sort.Sort(pixels(sorting.Pix))
	return sorting
}

func sortimage1(m image.Image) image.Image {
	sorting = makeRGB(m)
	sort.Sort((*pixim)(sorting))
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
	var im image.Image
	switch *method {
	case 0:
		im = sortimage0(inimg)
	case 1:
		im = sortimage1(inimg)
	default:
		log.Fatalf("invalid method: %v", *method)
	}

	saveImage(im, flag.Arg(1))
}
