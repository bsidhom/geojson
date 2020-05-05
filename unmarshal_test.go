package geojson

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalPoint(t *testing.T) {
	type Case struct {
		s        string
		expected Point
	}
	cases := []Case{
		{
			s: `{"type": "Point", "coordinates": [100, 3]}`,
			expected: Point{
				X:              100,
				Y:              3,
				Elevation:      0,
				HasElevation:   false,
				AdditionalData: map[string]interface{}{},
			},
		},
		{
			s: `{"type": "Point", "coordinates": [35, 5.6, 10]}`,
			expected: Point{
				X:              35,
				Y:              5.6,
				Elevation:      10,
				HasElevation:   true,
				AdditionalData: map[string]interface{}{},
			},
		},
		{
			// Ensure we preserve additional data.
			s: `{"type": "Point", "coordinates": [3, 3], "mydata": {"nested": "object"}}`,
			expected: Point{
				X: 3,
				Y: 3,
				AdditionalData: map[string]interface{}{
					"mydata": map[string]interface{}{
						"nested": "object",
					},
				},
			},
		},
	}
	for _, c := range cases {
		var p Point
		err := json.Unmarshal([]byte(c.s), &p)
		if err != nil {
			t.Errorf("error with unmarshal(%q): %v", c.s, err)
		}
		if !reflect.DeepEqual(p, c.expected) {
			t.Errorf("unmarshal(%q) should have been %#v, got %#v", c.s, c.expected, p)
		}
	}
}
