package shapes

import (
	"math"
)

// Circle represents a circle.
type Circle struct {
	Radius float64
}

// Area calculates the area of the circle.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}