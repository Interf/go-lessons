package handlers

import (
	"fmt"
	"pureProject/internal/figure/geometry"
	"strconv"
	"strings"
)

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func FormatPoint(p geometry.Point) string {
	return "[" + formatFloat(p.X) + "," + formatFloat(p.Y) + "]"
}

func FormatPolygonPoints(points []geometry.Point) string {
	parts := make([]string, 0, len(points))
	for _, p := range points {
		parts = append(parts, fmt.Sprintf("{%g,%g}", p.X, p.Y))
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
