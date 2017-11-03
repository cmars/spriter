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
	return NewMask([]string{
		"      ",
		"    ..",
		"    .|",
		"   ..|",
		"   ..|",
		"  ...|",
		" ...//",
		" ...//",
		" ...//",
		" ....|",
		"   ...",
		"      ",
	}, true, false)
}
