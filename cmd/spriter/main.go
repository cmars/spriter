package main

import (
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"

	"github.com/cmars/spriter"
)

func main() {
	options := spriter.DefaultOptions()
	options.Colored = true
	options.ColorVariations = 0.5
	options.Saturation = 0.8
	mask := spriter.Spaceship()
	log.Printf("mask can represent %d bits", mask.BitLen())
	g := spriter.NewGenerator(mask, options)
	m := g.Sprite()
	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	resized := resize.Resize(128, 128, m, resize.NearestNeighbor)
	err = png.Encode(f, resized)
	if err != nil {
		panic(err)
	}
}
