package main

import (
	"github.com/ajmadsen/pixelsort"
	"image"
	"image/draw"
	"os"
	"log"
	"math/rand"
	"time"
	"flag"
	"runtime/pprof"

  // Import these packages for their side effects, namely decoding
  // images of different formats.
  "image/png"
  _ "image/jpeg"
  _ "image/gif"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write mem profile to file")
	infile = flag.String("i", "", "file to read")
	outfile = flag.String("o", "sorted.png", "file to write")
)


func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	rand.Seed(time.Now().UnixNano())

	r, err := os.Open(*infile)
	if err != nil {
		log.Fatalf("failed to open file %v: %v\n", *infile, err)
	}

	log.Println("Loading image")
	im, _, err := image.Decode(r)
	if err != nil {
		log.Fatalf("failed to decode image %v: %v\n", *infile, err)
	}
	err = r.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Convert to usable format
	b := im.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)

	log.Println("Sorting image")
	re := pixelsort.ColEnum(m)
	sorter := pixelsort.PixelSort(re, pixelsort.HueSorter)
	sorter.Sort()

	log.Println("Saving image")
	os.Remove("sorted.png")
	out, err := os.Create(*outfile)
	if err != nil {
		log.Fatalf("could not create file: %v\n", err.Error())
	}

	err = png.Encode(out, m)
	if err != nil {
		log.Fatal(err)
	}
	
	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}

}

