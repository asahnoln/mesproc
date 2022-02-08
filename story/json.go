package story

import (
	"encoding/json"
	"io"
)

type JSONExpectGeo struct {
	Lat, Lon  float64
	Precision float64
}

type JSONStep struct {
	Command   bool
	Expect    *string
	Response  *string
	Fail      string
	Responses []string
	ExpectGeo *JSONExpectGeo
}

// Load loads story steps from given JSON file. Structure should be as follows:
//   [
//     {
//       "command": true,
//       "expect": "start",
//       "response": "let's start"
//     },
//     {
//       "expect": "go to step 2",
//       "response": "now at step 2",
//       "fail": "still at step 1"
//     },
//     {
//       "expectGeo": {
//         "lat": 43.257169,
//         "lon": 76.924515,
//         "precision": 50
//       },
//       "response": "proper geo",
//       "fail": "still waiting for geo"
//     },
//     {
//       "expect": "finish",
//       "response": "now finished",
//       "fail": "still at step 2"
//     }
//   ]
func Load(r io.Reader) (*Story, error) {
	s := New()
	steps := make([]JSONStep, 0)
	err := json.NewDecoder(r).Decode(&steps)
	if err != nil {
		return s, err
	}

	for _, ss := range steps {
		step := NewStep().Fail(ss.Fail)

		switch {
		case ss.Response != nil:
			step.Respond(*ss.Response)
		case ss.Responses != nil:
			step.Respond(ss.Responses...)
		}

		if ss.Expect != nil {
			step = step.Expect(*ss.Expect)
		} else {
			step = step.ExpectGeo(ss.ExpectGeo.Lat, ss.ExpectGeo.Lon, ss.ExpectGeo.Precision)
		}

		if ss.Command {
			s.AddCommand(step)
		} else {
			s.Add(step)
		}
	}

	return s, nil
}
