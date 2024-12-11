package shapes

// Square represents a square.
type Square struct {
	Side float64
}

// Area calculates the area of the square.
func (s Square) Area() float64 {
	return s.Side * s.Side
}