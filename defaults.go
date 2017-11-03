package spriter

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
