package main

import (
	"github.com/ajmadsen/pixelsort"
  "image"
	"os"
	"fmt"

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

	fmt.Println("Sorting image")
	sorter := pixelsort.Hue(im)
	sorter.Sort()

	fmt.Println("Saving image")
	os.Remove("sorted.png")
	out, err := os.Create("sorted.png")
	if err != nil {
		fmt.Printf("could not create file: %v\n", err.Error())
		os.Exit(-1)
	}

	png.Encode(out, sorter.Image())
}

