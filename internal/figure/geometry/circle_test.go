package geometry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCircleContains(t *testing.T) {
	c := Circle{
		Center: Point{
			X: 5,
			Y: 5,
		},
		Radius: 5,
	}

	tests := []struct {
		name string
		p    Point
		want bool
	}{
		{
			name: "inside",
			p: Point{
				X: 0,
				Y: 0,
			},
			want: false,
		},
		{
			name: "inside",
			p: Point{
				X: 5,
				Y: 5,
			},
			want: true,
		},
		{
			name: "on border",
			p: Point{
				X: 10,
				Y: 5,
			},
			want: true,
		},
		{
			name: "on border",
			p: Point{
				X: 10,
				Y: 10,
			},
			want: false,
		},
		{
			name: "outside",
			p: Point{
				X: 100,
				Y: 100,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.Contains(tt.p)

			if got != tt.want {
				t.Errorf("Contains() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircleArea(t *testing.T) {
	tests := []struct {
		name string
		c    Circle
		want float64
	}{
		{
			name: "correct result",
			c: Circle{
				Center: Point{
					X: 5,
					Y: 5,
				},
				Radius: 5,
			},
			want: 78.5,
		},
		{
			name: "null radius",
			c: Circle{
				Center: Point{
					X: -5,
					Y: -1,
				},
				Radius: 0,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.Area()

			assert.InDelta(t, tt.want, got, 0.1)
		})
	}
}
