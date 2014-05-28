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
)

func less(c1, c2 color.RGBA) bool {
	return uint32(c1.R)*uint32(c1.G)*uint32(c1.B) < uint32(c2.R)*uint32(c2.G)*uint32(c2.B)
}

func swap(m *image.RGBA, p1, p2 image.Point) {
	i, j := m.PixOffset(p1.X, p1.Y), m.PixOffset(p2.X, p2.Y)
	m.Pix[i+0], m.Pix[j+0] = m.Pix[j+0], m.Pix[i+0]
	m.Pix[i+1], m.Pix[j+1] = m.Pix[j+1], m.Pix[i+1]
	m.Pix[i+2], m.Pix[j+2] = m.Pix[j+2], m.Pix[i+2]
	m.Pix[i+3], m.Pix[j+3] = m.Pix[j+3], m.Pix[i+3]
}

func minPix(m *image.RGBA, center image.Point) image.Point {
	b := m.Bounds()
	minP := center
	minC := m.At(minP.X, minP.Y).(color.RGBA)
	for p := b.Min; p.Y < b.Max.Y; p.Y++ {
		for p.X = b.Min.X; p.X < b.Max.X; p.X++ {
			if p == center {
				continue
			}
			test := m.At(p.X, p.Y).(color.RGBA)
			if less(test, minC) {
				minC, minP = test, p
			}
		}
	}
	return minP
}

func sort(m image.Image) image.Image {
	b := m.Bounds()
	rgb, ok := m.(*image.RGBA)
	if !ok {
		rgb = image.NewRGBA(b)
		draw.Draw(rgb, b, m, b.Min, draw.Src)
	}

	window := image.Rect(-1, -1, 2, 2)
	for i := 0; i < 10; i++ {
		for p := b.Min; p.Y < b.Max.Y; p.Y++ {
			for p.X = b.Min.X; p.X < b.Max.X; p.X++ {
				sub := rgb.SubImage(window.Add(p)).(*image.RGBA)
				min := minPix(sub, p)
				if min != p {
					swap(rgb, min, p)
				}
			}
		}
	}

	return rgb
}

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Errorf("Usage: %s [flags] <input image> <output image>\n", os.Args[0])
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
	im := sort(inimg)

	log.Print("Opening output")
	outfile, err := os.Create(flag.Arg(1))
	defer outfile.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Writing output")
	png.Encode(outfile, im)
}
