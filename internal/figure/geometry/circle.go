package geometry

import "math"

type Circle struct {
	Center Point
	Radius float64
}

func (c *Circle) Contains(p Point) bool {
	return p.DistanceTo(c.Center) <= c.Radius
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
