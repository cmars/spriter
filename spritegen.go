package spritegen

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"image"
	"image/color"
	"math"
	"math/big"
	mathrand "math/rand"
)

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

func (m *Mask) BitLen() int {
	n := 0
	for i := range m.Bitmap {
		switch m.Bitmap[i] {
		case PixelEmptyOrBody, PixelBorderOrBody:
			n++
		}
	}
	return n
}

func (m *Mask) chooseBody(f Flipper) {
	for i := range m.Bitmap {
		switch m.Bitmap[i] {
		case PixelEmptyOrBody:
			if f.Next() {
				m.Bitmap[i] = PixelEmpty
			} else {
				m.Bitmap[i] = PixelBody
			}
		case PixelBorderOrBody:
			if f.Next() {
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

type Flipper interface {
	Next() bool
	Int64() int64
}

type SeededFlipper struct {
	*mathrand.Rand
}

func NewSeededFlipper(seed int64) Flipper {
	return &SeededFlipper{mathrand.New(mathrand.NewSource(seed))}
}

func (f *SeededFlipper) Next() bool {
	return f.Float64() > 0.5
}

func (f *SeededFlipper) Int64() int64 {
	return f.Int63()
}

type BytesFlipper struct {
	b big.Int
	i int
}

func NewBytesFlipper(b []byte) Flipper {
	f := &BytesFlipper{}
	f.b.SetBytes(b)
	f.i = 0
	return f
}

func (f *BytesFlipper) Next() bool {
	if f.i >= f.b.BitLen() {
		f.i = 0
	}
	result := f.b.Bit(f.i) > 0
	f.i++
	return result
}

func (f *BytesFlipper) Int64() int64 {
	return int64(binary.BigEndian.Uint64(f.b.Bytes()[:8]))
}

type RandomFlipper struct {
	*BytesFlipper
}

func NewRandomFlipper() Flipper {
	f := NewBytesFlipper(randomBytes(32)).(*BytesFlipper)
	return &RandomFlipper{f}
}

func (f *RandomFlipper) Next() bool {
	if f.BytesFlipper.i >= f.b.BitLen() {
		f.BytesFlipper.b.SetBytes(randomBytes(32))
	}
	return f.BytesFlipper.Next()
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := cryptorand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

type Options struct {
	Colored         bool
	EdgeBrightness  float64
	ColorVariations float64
	BrightnessNoise float64
	Saturation      float64
	Flipper         Flipper
}

func DefaultOptions() *Options {
	return &Options{
		Colored:         false,
		EdgeBrightness:  0.3,
		ColorVariations: 0.2,
		BrightnessNoise: 0.3,
		Saturation:      0.5,
		Flipper:         NewRandomFlipper(),
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
	r       *mathrand.Rand
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
