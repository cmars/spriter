package spriter

import (
	"image"
	"image/color"
	"math"
	mathrand "math/rand"
)

type Options struct {
	Colored         bool
	EdgeBrightness  float64
	ColorVariations float64
	BrightnessNoise float64
	Saturation      float64
	Flipper         Flipper
}

type Generator struct {
	mask    *Mask
	options *Options
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
	t.chooseBody(g.options.Flipper)
	t.chooseEdges()
	m := image.NewRGBA(image.Rect(0, 0, g.mask.ImageHeight(), g.mask.ImageWidth()))

	r := mathrand.New(mathrand.NewSource(g.options.Flipper.Int64()))

	var (
		borderColor        color.RGBA
		isVerticalGradient bool
		ulen, vlen         int
		hue                float64 = r.Float64()
		saturation         float64 = math.Max(
			math.Min(r.Float64()*g.options.Saturation, 1.0),
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
	if r.Float64() > 0.5 {
		isVerticalGradient = true
		ulen = g.mask.ImageHeight()
		vlen = g.mask.ImageWidth()
	} else {
		ulen = g.mask.ImageWidth()
		vlen = g.mask.ImageHeight()
	}

	for u := 0; u < ulen; u++ {
		isNewColor := math.Abs(((r.Float64()*2.0 - 1.0) +
			(r.Float64()*2.0 - 1.0) +
			(r.Float64()*2.0 - 1.0)) / 3.0)
		if isNewColor > (1.0 - g.options.ColorVariations) {
			hue = r.Float64()
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
