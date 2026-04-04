package geometry

import "strings"

type PointsFlag struct {
	RawPoints []string
}

func (p *PointsFlag) String() string {
	return strings.Join(p.RawPoints, " ")
}

func (p *PointsFlag) Set(value string) error {
	p.RawPoints = append(p.RawPoints, value)
	return nil
}
