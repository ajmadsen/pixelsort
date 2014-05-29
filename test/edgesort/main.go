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
	"sort"
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

	thresh  = flag.Int("t", 0, "threshold for the sobel filter")
	compute = flag.Bool("compute", false, "compute best threshold and exit")
	psobel  = flag.Bool("sobel", false, "only dump sobel and exit")
	phist   = flag.Bool("hist", false, "print histogram")
)

type pixim struct {
	*image.RGBA
	region []int
}

func (p *pixim) Len() int {
	b := p.Bounds()
	return b.Dx() * b.Dy()
}

func (p *pixim) Less(i, j int) bool {
	ri, rj := region(p.region, i), region(p.region, j)
	if ri != rj {
		return ri < rj
	}

	b := p.Bounds()
	c1, c2 := p.At(b.Min.X+i, b.Min.Y).(color.RGBA), p.At(b.Min.X+j, b.Min.Y).(color.RGBA)
	dist1 := uint32(c1.R)*uint32(c1.R) + uint32(c1.G)*uint32(c1.G) + uint32(c1.B)*uint32(c1.B)
	dist2 := uint32(c2.R)*uint32(c2.R) + uint32(c2.G)*uint32(c2.G) + uint32(c2.B)*uint32(c2.B)
	return dist1 < dist2
}

func (p *pixim) Swap(i, j int) {
	b := p.Bounds()
	pi := p.PixOffset(b.Min.X+i, b.Min.Y)
	pj := p.PixOffset(b.Min.X+j, b.Min.Y)
	p.Pix[pi+0], p.Pix[pj+0] = p.Pix[pj+0], p.Pix[pi+0]
	p.Pix[pi+1], p.Pix[pj+1] = p.Pix[pj+1], p.Pix[pi+1]
	p.Pix[pi+2], p.Pix[pj+2] = p.Pix[pj+2], p.Pix[pi+2]
	p.Pix[pi+3], p.Pix[pj+3] = p.Pix[pj+3], p.Pix[pi+3]
}

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

func sobel(m *image.Gray) (*image.Gray, []int) {
	b := m.Bounds()
	im := image.NewGray(b)
	hist := make([]int, 256)

	for p := b.Min; p.Y < b.Max.Y; p.Y++ {
		for p.X = b.Min.X; p.X < b.Max.X; p.X++ {
			window := kernR.Add(p).Intersect(b)
			mag := convolve(m.SubImage(window).(*image.Gray), kernX, kernY, window.Sub(p))
			im.SetGray(p.X, p.Y, color.Gray{mag})
			hist[mag]++
		}
	}

	return im, hist
}

func threshold(m *image.Gray, t uint8) {
	b := m.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := m.At(x, y).(color.Gray)
			if c.Y > t {
				c.Y = 255
			} else {
				c.Y = 0
			}
			m.Set(x, y, c)
		}
	}
}

// otsu's method
func computeThresh(hist []int, total int) int {
	sum := float64(0)
	for i := 0; i < 256; i++ {
		sum += float64(i * hist[i])
	}

	rsum := float64(0)
	bgWeight := float64(0)
	fgWeight := float64(0)
	bgMean := float64(0)
	fgMean := float64(0)

	bcMax := float64(0)
	bcBin := int(0)

	for i := 0; i < 256; i++ {
		bgWeight += float64(hist[i])
		if bgWeight == 0 {
			continue
		}
		if bgWeight == sum {
			break
		}
		rsum += float64(i) * float64(hist[i])
		fgWeight = float64(total) - bgWeight
		bgMean = rsum / bgWeight
		fgMean = (sum - rsum) / fgWeight
		bcVar := fgWeight * bgWeight * math.Pow(bgMean-fgMean, 2)
		if bcVar > bcMax {
			bcMax, bcBin = bcVar, i
		}
	}

	return bcBin
}

func edgeToRegion(m *image.Gray) [][]int {
	b := m.Bounds()
	regions := make([][]int, b.Dy())
	for i := range regions {
		for x := b.Min.X; x < b.Max.X; x++ {
			if m.Pix[m.Stride*i+x] > 0 {
				regions[i] = append(regions[i], x)
			}
		}
	}
	return regions
}

func region(regions []int, idx int) int {
	for i := range regions {
		if idx <= regions[i] {
			return i
		}
	}
	return len(regions)
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

func sortimage(m image.Image, regions [][]int) image.Image {
	b := m.Bounds()
	rgb := makeRGB(m)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		sub := rgb.SubImage(image.Rect(b.Min.X, b.Min.Y+y, b.Max.X, b.Min.Y+y+1))
		r := &pixim{
			sub.(*image.RGBA),
			regions[y],
		}
		sort.Sort(r)
	}
	return rgb
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
	s, hist := sobel(gray)

	if *phist {
		log.Print(hist)
	}

	if *compute {
		t := computeThresh(hist, b.Dx()*b.Dy())
		log.Printf("Best threshold = %v", t)

		threshold(s, uint8(t))

		sout, err := os.Create(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		png.Encode(sout, s)
		sout.Close()
		os.Exit(0)
	}

	threshold(s, uint8(*thresh))

	if *psobel {
		sout, err := os.Create(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		png.Encode(sout, s)
		sout.Close()
		os.Exit(0)
	}

	log.Println("Computing regions")
	regions := edgeToRegion(s)

	log.Println("Sorting image")
	sorted := sortimage(inimg, regions)

	log.Println("Writing output")
	outfile, err := os.Create(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(outfile, sorted)
	if err != nil {
		log.Fatal(err)
	}
}
