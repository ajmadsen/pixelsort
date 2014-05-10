package main

import (
	"github.com/ajmadsen/pixelsort"
	"image"
	"image/draw"
	"os"
	"fmt"
	"math/rand"
	"time"

  // Import these packages for their side effects, namely decoding
  // images of different formats.
  "image/png"
  _ "image/jpeg"
  _ "image/gif"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %v <image>\n", os.Args[0])
		os.Exit(-1)
	}

	rand.Seed(time.Now().UnixNano())

	r, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("failed to open file %v: %v\n", os.Args[1], err.Error())
		os.Exit(-1)
	}

	fmt.Println("Loading image")
	im, _, err := image.Decode(r)
	if err != nil {
		fmt.Printf("failed to decode image %v: %v\n", os.Args[1], err.Error())
		os.Exit(-1)
	}

	// Convert to usable format
	b := im.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)

	fmt.Println("Sorting image")
	re := pixelsort.RowEnum(m)
	sorter := pixelsort.PixelSort(re, pixelsort.HueSorter)
	sorter.Sort()

	fmt.Println("Saving image")
	os.Remove("sorted.png")
	out, err := os.Create("sorted.png")
	if err != nil {
		fmt.Printf("could not create file: %v\n", err.Error())
		os.Exit(-1)
	}

	png.Encode(out, m)
}

