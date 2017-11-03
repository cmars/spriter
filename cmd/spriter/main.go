package main

import (
	"image/png"
	"log"
	"os"

	"github.com/cmars/spriter"
)

func main() {
	options := spriter.DefaultOptions()
	options.Colored = true
	mask := spriter.Spaceship()
	log.Printf("mask can represent %d bits", mask.BitLen())
	g := spriter.NewGenerator(mask, options)
	m := g.Sprite()
	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, m)
	if err != nil {
		panic(err)
	}
}
