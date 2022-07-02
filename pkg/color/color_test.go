package color

import (
	"testing"
)

func TestRGB(t *testing.T) {
	tests := []struct {
		name  string
		red   int
		green int
		blue  int
		want  int
	}{
		{
			name:  "red",
			red:   255,
			green: 0,
			blue:  0,
			want:  196,
		},
		{
			name:  "green",
			red:   0,
			green: 255,
			blue:  0,
			want:  46,
		},
		{
			name:  "blue",
			red:   0,
			green: 0,
			blue:  255,
			want:  21,
		},
		{
			name:  "white",
			red:   255,
			green: 255,
			blue:  255,
			want:  231,
		},
		{
			name:  "black",
			red:   0,
			green: 0,
			blue:  0,
			want:  16,
		},
		{
			name:  "violet",
			red:   215,
			green: 135,
			blue:  255,
			want:  177,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RGB(tt.red, tt.green, tt.blue); got != tt.want {
				t.Errorf("RGB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelativeBrightness(t *testing.T) {
	tests := []struct {
		name  string
		red   int
		green int
		blue  int
		want  float64
	}{
		{
			name:  "red",
			red:   255,
			green: 0,
			blue:  0,
			want:  0.2126,
		},
		{
			name:  "green",
			red:   0,
			green: 255,
			blue:  0,
			want:  0.7152,
		},
		{
			name:  "blue",
			red:   0,
			green: 0,
			blue:  255,
			want:  0.0722,
		},
		{
			name:  "white",
			red:   255,
			green: 255,
			blue:  255,
			want:  1,
		},
		{
			name:  "black",
			red:   0,
			green: 0,
			blue:  0,
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RelativeBrightness(tt.red, tt.green, tt.blue); got != tt.want {
				t.Errorf("RelativeBrightness() = %v, want %v", got, tt.want)
			}
		})
	}
}
