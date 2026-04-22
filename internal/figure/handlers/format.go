package handlers

import (
	"fmt"
	"pureProject/internal/figure/geometry"
	"strings"
)

func FormatPolygonPoints(points []geometry.Point) string {
	parts := make([]string, 0, len(points))
	for _, p := range points {
		parts = append(parts, fmt.Sprintf("{%g,%g}", p.X, p.Y))
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
