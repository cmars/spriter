package main

import (
	"image/png"
	"log"
	"os"

	"spritegen"
)

func main() {
	options := spritegen.DefaultOptions()
	options.Colored = true
	mask := spritegen.Spaceship()
	log.Printf("mask can represent %d bits", mask.BitLen())
	g := spritegen.NewGenerator(mask, options)
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
