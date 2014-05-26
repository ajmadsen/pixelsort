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
	"runtime"
)

var (
	offsets = []image.Point{
		image.Point{-1, -1},
		image.Point{+0, -1},
		image.Point{+1, -1},
		image.Point{-1, +0},
		image.Point{+1, +0},
		image.Point{-1, +1},
		image.Point{+0, +1},
		image.Point{+1, +1},
	}
)

type hist []int

func newHist(size int) hist {
	h := make(hist, size)
	for i := range h {
		h[i] = 0
	}
	return h
}

type workUnit struct {
	filename string
	image    *image.Gray
	min, max int
}

func thresholdAndSave(work <-chan workUnit, kill <-chan chan<- int) {
	for {
		select {
		case u := <-work:
			log.Printf("Thresholding [%v, %v]", u.min, u.max)

			outfile, err := os.Create(u.filename)
			if err != nil {
				log.Fatal(err)
			}

			thresh := scaleImage(u.image, u.min, u.max)
			invert(thresh)

			log.Printf("Saving %s", u.filename)
			err = png.Encode(outfile, thresh)
			if err != nil {
				log.Fatal(err)
			}

			err = outfile.Close()
			if err != nil {
				log.Fatal(err)
			}
		case quit := <-kill:
			quit <- 1
			return
		}
	}
}

func findThreshold(h hist) (int, int) {
	min, max := int(-1), int(-1)
	total := uint64(0)
	for i := range h {
		total += uint64(h[i])
	}
	cnt := uint64(0)
	for i := range h {
		cnt += uint64(h[i])
		if h[i] > 0 && min < 0 {
			min = i
		} else if min >= 0 && h[i] == 0 && max < 0 {
			max = i-1
		}
	}
	if min > max {
		return max, min
	}
	return min, max
}

func cDiff(c1, c2 color.Color) uint8 {
	g1 := c1.(color.Gray)
	g2 := c2.(color.Gray)
	if g1.Y > g2.Y {
		return g1.Y - g2.Y
	}
	return g2.Y - g1.Y
}

func imageGrad(m *image.Gray) (*image.Gray, hist) {
	b := m.Bounds()
	imrgb := image.NewGray(b)

	h := newHist(math.MaxUint8 + 1)
	for p := b.Min; p.Y < b.Max.Y; p.Y = p.Y + 1 {
		for p.X = b.Min.X; p.X < b.Max.X; p.X += 1 {
			diff := int32(0)
			c := int32(0)
			for _, off := range offsets {
				pt := p.Add(off)
				if pt.In(b) {
					c++
					diff += int32(cDiff(m.At(p.X, p.Y), m.At(pt.X, pt.Y)))
				}
			}
			diff /= c

			h[diff] += 1

			// diff should be in the range of [0.0, 1.0], indicating
			// the degree of difference between the neighboring pixels
			// and the current one
			//log.Print(diff)
			imrgb.Set(p.X, p.Y, color.Gray{uint8(diff)})
		}
	}

	return imrgb, h
}

func threshold(x, min, max int) uint8 {
	x = (x - min) * math.MaxUint8 / max
	if x < 0 {
		x = 0
	}
	return uint8(x)
}

func bw(x, min int) uint8 {
	if x < min {
		return uint8(0)
	}
	return uint8(255)
}

func scaleImage(m *image.Gray, min, max int) *image.Gray {
	b := m.Bounds()
	im := image.NewGray(b)
	max -= min

	for j := b.Min.Y; j < b.Max.Y; j++ {
		for i := b.Min.X; i < b.Max.X; i++ {
			c := m.At(i, j).(color.Gray)
			c.Y = bw(int(c.Y), min)
			im.SetGray(i, j, c)
		}
	}

	return im
}

func invert(m *image.Gray) {
	b := m.Bounds()
	for j := b.Min.Y; j < b.Max.Y; j++ {
		for i := b.Min.X; i < b.Max.X; i++ {
			c := m.At(i, j).(color.Gray)
			c.Y = math.MaxUint8 - c.Y
			m.SetGray(i, j, c)
		}
	}
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

	log.Println("Computing pixel gradient")
	grad, h := imageGrad(gray)
	min, max := findThreshold(h)

	runtime.GOMAXPROCS(runtime.NumCPU())

	work := make(chan workUnit)
	kill := make(chan chan<- int, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		go thresholdAndSave(work, kill)
	}

	for i := 0; i < max-min; i++ {
		work <- workUnit{
			filename: fmt.Sprintf(flag.Arg(1), i),
			image:    grad,
			min:      min + i,
			max:      max,
		}
	}

	for i := 0; i < cap(kill); i++ {
		quit := make(chan int)
		kill <- quit
		<-quit
	}
}
