package shapes

// Rectangle represents a rectangle.
type Rectangle struct {
	Length float64
	Width  float64
}

// Area calculates the area of the rectangle.
func (r Rectangle) Area() float64 {
	return r.Length * r.Width
}