package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	options := DefaultOptions()
	options.Colored = true
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
	EdgeBrightness  float64
	ColorVariations float64
	BrightnessNoise float64
	Saturation      float64
	Random          *rand.Rand
}

func DefaultOptions() *Options {
	return &Options{
		Colored:         false,
		EdgeBrightness:  0.3,
		ColorVariations: 0.2,
		BrightnessNoise: 0.3,
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
	/*
		for y := 0; y < g.mask.ImageHeight(); y++ {
			for x := 0; x < g.mask.ImageWidth(); x++ {
				if g.mask.get(x, y) == PixelBorder {
					m.Set(x, y, color.RGBA{0, 0, 0, 255})
				} else {
					m.Set(x, y, color.RGBA{255, 255, 255, 255})
				}
			}
		}
	*/

	var (
		borderColor        color.RGBA
		isVerticalGradient bool
		ulen, vlen         int
		hue                float64 = g.options.Random.Float64()
		saturation         float64 = math.Max(
			math.Min(g.options.Random.Float64()*g.options.Saturation, 1.0),
			0)
	)
	if g.options.Colored {
		edgeColor := uint8(255.0 * g.options.EdgeBrightness)
		borderColor.R = edgeColor
		borderColor.G = edgeColor
		borderColor.B = edgeColor
		borderColor.A = 255
	} else {
		borderColor = color.RGBA{0, 0, 0, 255}
	}
	if g.options.Random.Float64() > 0.5 {
		isVerticalGradient = true
		ulen = g.mask.ImageHeight()
		vlen = g.mask.ImageWidth()
	} else {
		ulen = g.mask.ImageWidth()
		vlen = g.mask.ImageHeight()
	}

	for u := 0; u < ulen; u++ {
		isNewColor := math.Abs(((g.options.Random.Float64()*2.0 - 1.0) +
			(g.options.Random.Float64()*2.0 - 1.0) +
			(g.options.Random.Float64()*2.0 - 1.0)) / 3.0)
		if isNewColor > (1.0 - g.options.ColorVariations) {
			hue = g.options.Random.Float64()
		}

		for v := 0; v < vlen; v++ {
			var (
				pixel Pixel
				x, y  int
				c     color.RGBA
			)
			if isVerticalGradient {
				y, x = u, v
			} else {
				x, y = u, v
			}
			pixel = t.get(x, y)

			if g.options.Colored {
				if pixel == PixelBorder {
					c = borderColor
				} else if pixel == PixelBody {
					brightness := math.Sin(
						(float64(u)/float64(ulen))*math.Pi) * (1.0 - g.options.BrightnessNoise)
					c = hsl2rgb(hue, saturation, brightness)
				}
			} else if pixel == PixelBorder {
				c = borderColor
			}
			log.Println(c)

			m.Set(x, y, c)
		}
	}
	return m
}

func hsl2rgb(h, s, l float64) color.RGBA {
	var (
		i          int
		f, p, q, t float64
	)
	i = int(math.Floor(h * 6.0))
	f = h*6.0 - float64(i)
	p = l * (1.0 - s)
	q = l * (1.0 - f*s)
	t = l * (1.0 - (1.0-f)*s)

	var c color.RGBA
	c.A = 255
	switch i % 6 {
	case 0:
		c.R = uint8(l * 255.0)
		c.G = uint8(t * 255.0)
		c.B = uint8(p * 255.0)
	case 1:
		c.R = uint8(q * 255.0)
		c.G = uint8(l * 255.0)
		c.B = uint8(p * 255.0)
	case 2:
		c.R = uint8(p * 255.0)
		c.G = uint8(l * 255.0)
		c.B = uint8(t * 255.0)
	case 3:
		c.R = uint8(p * 255.0)
		c.G = uint8(q * 255.0)
		c.B = uint8(l * 255.0)
	case 4:
		c.R = uint8(t * 255.0)
		c.G = uint8(p * 255.0)
		c.B = uint8(l * 255.0)
	case 5:
		c.R = uint8(l * 255.0)
		c.G = uint8(p * 255.0)
		c.B = uint8(q * 255.0)
	}
	return c
}
