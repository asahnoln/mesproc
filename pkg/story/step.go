package story

import (
	"fmt"
	"math"

	"github.com/asahnoln/mesproc/pkg/store"
)

// Step is a building block of a story.
// It holds information on what message it expects from the user to advance the story
// and how it would respond to proper or a wrong message.
type Step struct {
	expectation string
	responses   []string
	failMessage string
	isGeo       bool
	geoExp      [3]float64
	store       store.Step
	additional  map[int]map[string]interface{}
}

// NewStep returns a new Step
func NewStep() *Step {
	return &Step{}
}

// Expectation returns expected message to advance the story.
func (s *Step) Expectation() string {
	return s.expectation
}

// Response returns response of the Step
func (s *Step) Response() string {
	return s.responses[0]
}

// Responses returns response of the Step
func (s *Step) Responses() []string {
	return s.responses
}

// FailMessage returns a fail message which is given when given message is not what expected.
func (s *Step) FailMessage() string {
	return s.failMessage
}

// Expect sets expected message for the Step
func (s *Step) Expect(e string) *Step {
	s.expectation = e
	return s
}

// Respond sets response for the right message
func (s *Step) Respond(r ...string) *Step {
	s.responses = r
	return s
}

// Fail sets a fail message for the step which is returned when given input is not what expected in the step.
func (s *Step) Fail(e string) *Step {
	s.failMessage = e
	return s
}

// ExpectGeo sets expectation for the step to be a geo location instead of plain text
func (s *Step) ExpectGeo(lat, lon float64, precision float64) *Step {
	s.isGeo = true
	s.geoExp = [3]float64{lat, lon, precision}
	return s
}

// ExpectSave prepares the step to save incoming message
func (s *Step) ExpectSave(store store.Step) *Step {
	s.store = store
	return s
}

func (s *Step) Additional(step int, field string, value interface{}) *Step {
	if s.additional == nil {
		s.additional = make(map[int]map[string]interface{})
	}
	if s.additional[step] == nil {
		s.additional[step] = make(map[string]interface{})
	}

	s.additional[step][field] = value
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
