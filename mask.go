package spriter

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
