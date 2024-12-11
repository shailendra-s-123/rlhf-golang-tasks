package shapes

import (
	"testing"
)

func TestCircleArea(t *testing.T) {
	tests := []struct {
		name      string
		radius    float64
		expected  float64
		tolerance float64
	}{
		{"Unit circle", 1, math.Pi, 0.01},
		{"Circle with radius 2", 2, 4 * math.Pi, 0.01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Circle{Radius: tt.radius}
			area := c.Area()
			if math.Abs(area-tt.expected) > tt.tolerance {
				t.Errorf("expected %f, got %f", tt.expected, area)
			}
		})
	}
}