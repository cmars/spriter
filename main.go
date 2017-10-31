package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	options := DefaultOptions()
	g := NewGenerator(Spaceship(), options)
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

type Pixel int8

const (
	PixelBorder       = Pixel(-1)
	PixelEmpty        = iota
	PixelEmptyOrBody  = iota
	PixelBorderOrBody = iota
	PixelBody         = iota

	p_ = PixelBorder
	p0 = PixelEmpty
	p1 = PixelEmptyOrBody
	p2 = PixelBorderOrBody
)

type Mask struct {
	Bitmap     []Pixel
	MaskWidth  int
	MaskHeight int
	MirrorX    bool
	MirrorY    bool
}

func (m *Mask) ImageWidth() int {
	if m.MirrorX {
		return m.MaskWidth * 2
	}
	return m.MaskWidth
}

func (m *Mask) ImageHeight() int {
	if m.MirrorY {
		return m.MaskHeight * 2
	}
	return m.MaskHeight
}

func (m *Mask) get(x, y int) Pixel {
	if x >= m.MaskWidth && m.MirrorX {
		x = m.MaskWidth - (x - m.MaskWidth) - 1
	}
	if y >= m.MaskHeight && m.MirrorY {
		y = m.MaskHeight - (y - m.MaskHeight) - 1
	}
	return m.Bitmap[y*m.MaskWidth+x]
}

func (m *Mask) chooseBody(r *rand.Rand) {
	for i := range m.Bitmap {
		switch m.Bitmap[i] {
		case PixelEmptyOrBody:
			if r.Float64() < 0.5 {
				m.Bitmap[i] = PixelEmpty
			} else {
				m.Bitmap[i] = PixelBody
			}
		case PixelBorderOrBody:
			if r.Float64() < 0.5 {
				m.Bitmap[i] = PixelBorder
			} else {
				m.Bitmap[i] = PixelBody
			}
		}
	}
}

func (m *Mask) chooseEdges() {
	for y := 0; y < m.MaskHeight; y++ {
		for x := 0; x < m.MaskWidth; x++ {
			if m.Bitmap[y*m.MaskWidth+x] > PixelEmpty {
				above := (y-1)*m.MaskWidth + x
				if y-1 >= 0 && m.Bitmap[above] == PixelEmpty {
					m.Bitmap[above] = PixelBorder
				}
				if !m.MirrorY {
					below := (y+1)*m.MaskWidth + x
					if y+1 < m.MaskHeight && m.Bitmap[below] == PixelEmpty {
						m.Bitmap[below] = PixelBorder
					}
				}
				left := y*m.MaskWidth + x - 1
				if x-1 >= 0 && m.Bitmap[left] == PixelEmpty {
					m.Bitmap[left] = PixelBorder
				}
				if !m.MirrorX {
					right := y*m.MaskWidth + x + 1
					if x+1 < m.MaskWidth && m.Bitmap[right] == PixelEmpty {
						m.Bitmap[right] = PixelBorder
					}
				}
			}
		}
	}
}

type Options struct {
	Colored         bool
	Edgep_ightness  float64
	ColorVariations float64
	p_ightnessNoise float64
	Saturation      float64
	Random          *rand.Rand
}

func DefaultOptions() *Options {
	return &Options{
		Colored:         false,
		Edgep_ightness:  0.3,
		ColorVariations: 0.2,
		p_ightnessNoise: 0.3,
		Saturation:      0.5,
		Random:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func Spaceship() *Mask {
	return &Mask{
		Bitmap: []Pixel{
			p0, p0, p0, p0, p0, p0,
			p0, p0, p0, p0, p1, p1,
			p0, p0, p0, p0, p1, p_,
			p0, p0, p0, p1, p1, p_,
			p0, p0, p0, p1, p1, p_,
			p0, p0, p1, p1, p1, p_,
			p0, p1, p1, p1, p2, p2,
			p0, p1, p1, p1, p2, p2,
			p0, p1, p1, p1, p2, p2,
			p0, p1, p1, p1, p1, p_,
			p0, p0, p0, p1, p1, p1,
			p0, p0, p0, p0, p0, p0,
		},
		MaskWidth:  6,
		MaskHeight: 12,
		MirrorX:    true,
		MirrorY:    false,
	}
}

type Generator struct {
	mask    *Mask
	options *Options
	r       *rand.Rand
}

func NewGenerator(mask *Mask, options *Options) *Generator {
	if options == nil {
		options = DefaultOptions()
	}
	return &Generator{
		mask:    mask,
		options: options,
	}
}

func (g *Generator) Sprite() *image.RGBA {
	t := *g.mask
	t.chooseBody(g.options.Random)
	t.chooseEdges()
	m := image.NewRGBA(image.Rect(0, 0, g.mask.ImageHeight(), g.mask.ImageWidth()))
	for y := 0; y < g.mask.ImageHeight(); y++ {
		for x := 0; x < g.mask.ImageWidth(); x++ {
			if g.mask.get(x, y) == PixelBorder {
				m.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				m.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}
	return m
}
