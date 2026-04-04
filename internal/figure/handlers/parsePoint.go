package handlers

import (
	"errors"
	"fmt"
	"pureProject/internal/figure/geometry"
	"strconv"
	"strings"
)

var ErrParsePointFormat = errors.New("invalid point format")
var ErrParsePointParseXFloat = errors.New("invalid x in point")
var ErrParsePointParseYFloat = errors.New("invalid y in point")

func ParsePoint(s string) (geometry.Point, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return geometry.Point{}, fmt.Errorf("%w, expected x,y: %s", ErrParsePointFormat, s)
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return geometry.Point{}, fmt.Errorf("%w %q: %s", ErrParsePointParseXFloat, parts[0], err)
	}

	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return geometry.Point{}, fmt.Errorf("%w %q: %s", ErrParsePointParseYFloat, s, err)
	}

	return geometry.Point{
		X: x,
		Y: y,
	}, nil
}
