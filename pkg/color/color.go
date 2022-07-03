package color

import "math"

func RGB(red, green, blue int) uint8 {
	return uint8(16 + 36*valueToIndex(red) + 6*valueToIndex(green) + valueToIndex(blue))
}

func RelativeBrightness(red, green, blue int) float64 {
	return 0.2126*relativeBrightnessComponent(float64(red)/255) + 0.7152*relativeBrightnessComponent(float64(green)/255) + 0.0722*relativeBrightnessComponent(float64(blue)/255)
}

func relativeBrightnessComponent(value float64) float64 {
	if value <= 0.03928 {
		return value / 12.92
	}
	return math.Pow(((value + 0.055) / 1.055), 2.4)
}

func valueToIndex(value int) int {
	switch {
	case value < 95/2:
		return 0
	case value < (95+135)/2:
		return 1
	case value < (135+175)/2:
		return 2
	case value < (175+215)/2:
		return 3
	case value < (215+255)/2:
		return 4
	default:
		return 5
	}
}
