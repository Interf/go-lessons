package geometry

import "math"

const AreaMinPoints = 3

type Polygon struct {
	Points []Point
}

func (poly *Polygon) Contains(p Point) bool {
	n := len(poly.Points)
	inside := false

	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		pi := poly.Points[i]
		pj := poly.Points[j]

		intersect := ((pi.Y > p.Y) != (pj.Y > p.Y)) &&
			(p.X < (pj.X-pi.X)*(p.Y-pi.Y)/(pj.Y-pi.Y)+pi.X)

		if intersect {
			inside = !inside
		}
	}

	return inside
}

func (poly *Polygon) Area() float64 {
	n := len(poly.Points)
	if n < AreaMinPoints {
		return 0
	}

	sum := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		sum += poly.Points[i].X*poly.Points[j].Y - poly.Points[j].X*poly.Points[i].Y
	}

	return math.Abs(sum) / 2
}
