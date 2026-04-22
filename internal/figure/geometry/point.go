package geometry

import (
	"math"
	"strconv"
)

type Point struct {
	X, Y float64
}

const eps = 1e-9

func (p *Point) DistanceTo(other Point) float64 {
	dx := other.X - p.X
	dy := other.Y - p.Y

	return math.Sqrt(dx*dx + dy*dy)
}

func (p *Point) IsInRadius(center Point, radius float64) bool {
	return (p.DistanceTo(center) - radius) <= eps
}

func (p Point) String() string {
	x := strconv.FormatFloat(p.X, 'f', -1, 64)
	y := strconv.FormatFloat(p.Y, 'f', -1, 64)

	return "[" + x + "," + y + "]"
}
