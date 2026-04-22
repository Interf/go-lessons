package geometry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointDistanceTo(t *testing.T) {

	p := Point{
		X: 1.0,
		Y: 2.0,
	}

	tests := []struct {
		name  string
		Point Point
		want  float64
	}{
		{
			name: "DistanceTo Point positive",
			Point: Point{
				X: 1.3,
				Y: 2.0,
			},
			want: 0.3,
		},
		{
			name: "DistanceTo Point negative",
			Point: Point{
				X: -1,
				Y: -2.2,
			},
			want: 4.6,
		},
		{
			name: "DistanceTo Point null",
			Point: Point{
				X: 0,
				Y: 0,
			},
			want: 2.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.DistanceTo(tt.Point)

			assert.InDelta(t, tt.want, got, 0.1)
		})
	}
}

func TestPointIsInRadius(t *testing.T) {
	p := Point{
		X: 1.0,
		Y: 2.0,
	}

	tests := []struct {
		name   string
		Point  Point
		radius float64
		want   bool
	}{
		{
			name: "Point in zero radius",
			Point: Point{
				X: 1.0,
				Y: 2.0,
			},
			radius: 0,
			want:   true,
		},
		{
			name: "Point in radius",
			Point: Point{
				X: 1.0,
				Y: 2.0,
			},
			radius: 5,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.IsInRadius(tt.Point, tt.radius)

			assert.Equal(t, tt.want, got)
		})
	}
}
