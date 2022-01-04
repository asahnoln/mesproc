package story

import (
	"fmt"
	"math"
)

type Step struct {
	expectation string
	response    string
	failMessage string
	isGeo       bool
	geoExp      [3]float64
}

func NewStep() *Step {
	return &Step{}
}

func (s *Step) Expectation() string {
	return s.expectation
}

func (s *Step) Response() string {
	return s.response
}

func (s *Step) FailMessage() string {
	return s.failMessage
}

func (s *Step) Expect(e string) *Step {
	s.expectation = e
	return s
}

func (s *Step) Respond(r string) *Step {
	s.response = r
	return s
}

func (s *Step) Fail(e string) *Step {
	s.failMessage = e
	return s
}

func (s *Step) ExpectGeo(lat, lon float64, precision float64) *Step {
	s.isGeo = true
	s.geoExp = [3]float64{lat, lon, precision}
	return s
}

func (s *Step) checkGeo(m string) bool {
	var lat, lon float64
	fmt.Sscanf(m, "%f,%f", &lat, &lon)
	return distance(s.geoExp[0], s.geoExp[1], lat, lon) <= s.geoExp[2]
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	p1, p2 := degToRad(lat1), degToRad(lat2)
	dp := p2 - p1
	dl := degToRad(lon2) - degToRad(lon1)
	r := 6371000.0

	a := math.Pow(math.Sin(dp/2), 2) + math.Cos(p1)*math.Cos(p2)*math.Pow(math.Sin(dl/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := r * c
	return d
}

func degToRad(x float64) float64 {
	return x * (math.Pi / 180.0)
}
