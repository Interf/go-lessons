package figure

import (
	"errors"
	"flag"
	"fmt"
	"pureProject/internal/figure/geometry"
	"pureProject/internal/figure/handlers"
)

var ErrCircleModeEmptyCenterRaw = errors.New("for circle mode --center is required")
var ErrCircleModeEmptyRadius = errors.New("for circle mode --radius must be > 0")
var ErrCircleModeParseCenter = errors.New("error while parsing center")
var ErrCircleModeParseCheckInside = errors.New("error while parsing check inside")

var ErrPolygonModeEmpty = errors.New("polygon doesn't contain enough coordinates")
var ErrPolygonModeMinPointCount = 3
var ErrPolygonModeParseRawPoints = errors.New("error while parsing raw points")
var ErrPolygonModeParseCheckInside = errors.New("error while parsing check inside")

type App struct{}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() error {

	var polygonPoints geometry.PointsFlag
	var isCircle bool
	var centerRaw string
	var radius float64
	var checkInsideRaw string

	flag.Var(&polygonPoints, "points", "polygon points in format x,y")
	flag.BoolVar(&isCircle, "circle", false, "use circle polygon")
	flag.StringVar(&centerRaw, "center", "", "center to use in circle polygon")
	flag.Float64Var(&radius, "radius", 0, "circle radius")
	flag.StringVar(&checkInsideRaw, "check-inside", "", "point to check in format x,y")

	flag.Parse()

	if isCircle {
		return runCircleMode(centerRaw, radius, checkInsideRaw)
	}

	return runPolygonMode(polygonPoints.RawPoints, checkInsideRaw)
}

func runCircleMode(centerRaw string, radius float64, checkInsideRaw string) error {
	if centerRaw == "" {
		return ErrCircleModeEmptyCenterRaw
	}
	if radius <= 0 {
		return ErrCircleModeEmptyRadius
	}

	center, err := handlers.ParsePoint(centerRaw)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCircleModeParseCenter, err)
	}

	circle := geometry.Circle{
		Center: center,
		Radius: radius,
	}

	if checkInsideRaw != "" {
		p, err := handlers.ParsePoint(checkInsideRaw)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrCircleModeParseCheckInside, err)
		}

		if circle.Contains(p) {
			fmt.Printf(
				"Point %s inside circle with center %s and radius %.2f\n",
				handlers.FormatPoint(p),
				handlers.FormatPoint(center),
				radius,
			)
		} else {
			fmt.Printf(
				"Point %s outside circle with center %s and radius %.2f\n",
				handlers.FormatPoint(p),
				handlers.FormatPoint(center),
				radius,
			)
		}

		return nil
	}

	fmt.Printf("Area of circle with radius %.2f is %.2f\n", radius, circle.Area())

	return nil
}

func runPolygonMode(rawPoints []string, checkInsideRaw string) error {
	if len(rawPoints) < ErrPolygonModeMinPointCount {
		return ErrPolygonModeEmpty
	}

	points := make([]geometry.Point, 0, len(rawPoints))
	for i := range rawPoints {
		p, err := handlers.ParsePoint(rawPoints[i])
		if err != nil {
			return fmt.Errorf("%w: %s", ErrPolygonModeParseRawPoints, err)
		}
		points = append(points, p)
	}

	polygon := geometry.Polygon{
		Points: points,
	}

	if checkInsideRaw != "" {
		p, err := handlers.ParsePoint(checkInsideRaw)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrPolygonModeParseCheckInside, err)
		}

		if polygon.Contains(p) {
			fmt.Printf("Point %s inside polygon\n", handlers.FormatPoint(p))
		} else {
			fmt.Printf("Point %s outside polygon\n", handlers.FormatPoint(p))
		}

		return nil
	}

	fmt.Printf("Area of %s is %g\n", handlers.FormatPolygonPoints(points), polygon.Area())

	return nil
}
