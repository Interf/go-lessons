package geometry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolygonContains(t *testing.T) {
	p := Polygon{
		Points: []Point{
			{0, 0},
			{1, 3},
			{2, 4},
		},
	}

	tests := []struct {
		name string
		p    Point
		want bool
	}{
		{
			name: "Point contains",
			p:    Point{1, 3},
			want: true,
		},
		{
			name: "Point not contains",
			p:    Point{100, 100},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.Contains(tt.p)

			assert.Equal(t, tt.want, got)
		})
	}

}

func TestPolygonArea(t *testing.T) {

	tests := map[string]struct {
		p    Polygon
		want float64
	}{
		"Three points": {
			p: Polygon{
				Points: []Point{
					{0, 0},
					{0, 1},
					{1, 1},
				},
			},
			want: 0.5,
		},

		"Four points": {
			p: Polygon{
				Points: []Point{
					{0, 0},
					{0, 1},
					{1, 1},
					{1, 0},
				},
			},
			want: 1,
		},

		"Four points with zero value": {
			p: Polygon{
				Points: []Point{
					{0, 0},
					{0, 0},
					{0, 0},
					{0, 0},
				},
			},
			want: 0,
		},

		"Four points with negative value": {
			p: Polygon{
				Points: []Point{
					{-12, -5},
					{-55, -2},
					{-22, -3},
					{-33, -2},
				},
			},
			want: 22,
		},

		"Two points": {
			p: Polygon{
				Points: []Point{
					{0, 0},
					{0, 1},
				},
			},
			want: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.p.Area()

			assert.InDelta(t, tt.want, got, 0.1)
		})
	}
}
