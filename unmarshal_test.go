package geojson

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAll_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		s        string
		expected Object
	}{
		{
			s:        `{"type":"Point","coordinates":[3,4]}`,
			expected: &Point{X: 3, Y: 4},
		},
		{
			s:        `{"type":"Point","coordinates":[3,4, 5]}`,
			expected: &Point{X: 3, Y: 4, Elevation: 5, HasElevation: true},
		},
		{
			s: `{"type":"MultiPoint","coordinates":[[100,0],[101,1]]}`,
			expected: &MultiPoint{
				Points: []Point{
					{X: 100, Y: 0},
					{X: 101, Y: 1},
				},
			},
		},
		{
			s: `{"type":"LineString","coordinates":[[100,0],[101,1]]}`,
			expected: &LineString{
				Points: []Point{
					{X: 100, Y: 0},
					{X: 101, Y: 1},
				},
			},
		},
		{
			s: `{"type":"MultiLineString","coordinates":[[[100,0],[101,1]],[[102,2],[103,3]]]}`,
			expected: &MultiLineString{
				Lines: []LineString{
					{
						Points: []Point{
							{X: 100, Y: 0},
							{X: 101, Y: 1},
						},
					},
					{
						Points: []Point{
							{X: 102, Y: 2},
							{X: 103, Y: 3},
						},
					},
				},
			},
		},
		{
			s: `{"type":"Polygon","coordinates":[[[100,0],[101,0],[101,1],[100,1],[100,0]],[[100.8,0.8],[100.8,0.2],[100.2,0.2],[100.2,0.8],[100.8,0.8]]]}`,
			expected: &Polygon{
				Rings: []LineString{
					{
						Points: []Point{
							{X: 100, Y: 0},
							{X: 101, Y: 0},
							{X: 101, Y: 1},
							{X: 100, Y: 1},
							{X: 100, Y: 0},
						},
					},
					{
						Points: []Point{
							{X: 100.8, Y: 0.8},
							{X: 100.8, Y: 0.2},
							{X: 100.2, Y: 0.2},
							{X: 100.2, Y: 0.8},
							{X: 100.8, Y: 0.8},
						},
					},
				},
			},
		},
		{
			s: `{
  "type": "MultiPolygon",
  "coordinates": [
    [
      [
        [102.0, 2.0],
        [103.0, 2.0],
        [103.0, 3.0],
        [102.0, 3.0],
        [102.0, 2.0]
      ]
    ],
    [
      [
        [100.0, 0.0],
        [101.0, 0.0],
        [101.0, 1.0],
        [100.0, 1.0],
        [100.0, 0.0]
      ],
      [
        [100.2, 0.2],
        [100.2, 0.8],
        [100.8, 0.8],
        [100.8, 0.2],
        [100.2, 0.2]
      ]
    ]
  ]
}`,
			expected: &MultiPolygon{
				Polygons: []Polygon{
					{
						Rings: []LineString{
							{
								Points: []Point{
									{X: 102.0, Y: 2.0},
									{X: 103.0, Y: 2.0},
									{X: 103.0, Y: 3.0},
									{X: 102.0, Y: 3.0},
									{X: 102.0, Y: 2.0},
								},
							},
						},
					},
					{
						Rings: []LineString{
							{
								Points: []Point{
									{X: 100.0, Y: 0.0},
									{X: 101.0, Y: 0.0},
									{X: 101.0, Y: 1.0},
									{X: 100.0, Y: 1.0},
									{X: 100.0, Y: 0.0},
								},
							},
							{
								Points: []Point{
									{X: 100.2, Y: 0.2},
									{X: 100.2, Y: 0.8},
									{X: 100.8, Y: 0.8},
									{X: 100.8, Y: 0.2},
									{X: 100.2, Y: 0.2},
								},
							},
						},
					},
				},
			},
		},
		{
			s: `{
  "type": "GeometryCollection",
  "geometries": [
    {
      "type": "Point",
      "coordinates": [100.0, 0.0]
    },
    {
      "type": "LineString",
      "coordinates": [
        [101.0, 0.0],
        [102.0, 1.0]
      ]
    }
  ]
}`,
			expected: &GeometryCollection{
				Geometries: []Geometry{
					&Point{X: 100, Y: 0},
					&LineString{
						Points: []Point{
							{X: 101, Y: 0},
							{X: 102, Y: 1},
						},
					},
				},
			},
		},
		{
			s: `{
  "type": "Feature",
  "bbox": [-10.0, -10.0, 10.0, 10.0],
  "geometry": {
    "type": "Polygon",
    "coordinates": [
      [
        [-10.0, -10.0],
        [10.0, -10.0],
        [10.0, 10.0],
        [-10.0, -10.0]
      ]
    ]
  },
  "properties": {
    "prop0": "value0"
  }
}`,
			expected: &Feature{
				Geometry: &Polygon{
					Rings: []LineString{
						{
							Points: []Point{
								{X: -10, Y: -10},
								{X: 10, Y: -10},
								{X: 10, Y: 10},
								{X: -10, Y: -10},
							},
						},
					},
				},
				Properties: map[string]interface{}{
					"prop0": "value0",
				},
			},
		},
		{
			s: `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [102.0, 0.5]
      },
      "properties": {
        "prop0": "value0"
      }
    },
    {
      "type": "Feature",
      "geometry": {
        "type": "LineString",
        "coordinates": [
          [102.0, 0.0],
          [103.0, 1.0],
          [104.0, 0.0],
          [105.0, 1.0]
        ]
      },
      "properties": {
        "prop0": "value0",
        "prop1": 0.0
      }
    },
    {
      "type": "Feature",
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [100.0, 0.0],
            [101.0, 0.0],
            [101.0, 1.0],
            [100.0, 1.0],
            [100.0, 0.0]
          ]
        ]
      },
      "properties": {
        "prop0": "value0",
        "prop1": {
          "this": "that"
        }
      }
    }
  ]
}`,
			expected: &FeatureCollection{
				Features: []Feature{
					{
						Geometry: &Point{X: 102, Y: 0.5},
						Properties: map[string]interface{}{
							"prop0": "value0",
						},
					},
					{
						Geometry: &LineString{
							Points: []Point{
								{X: 102.0, Y: 0.0},
								{X: 103.0, Y: 1.0},
								{X: 104.0, Y: 0.0},
								{X: 105.0, Y: 1.0},
							},
						},
						Properties: map[string]interface{}{
							"prop0": "value0",
							"prop1": float64(0),
						},
					},
					{
						Geometry: &Polygon{
							Rings: []LineString{
								{
									Points: []Point{
										{X: 100.0, Y: 0.0},
										{X: 101.0, Y: 0.0},
										{X: 101.0, Y: 1.0},
										{X: 100.0, Y: 1.0},
										{X: 100.0, Y: 0.0},
									},
								},
							},
						},
						Properties: map[string]interface{}{
							"prop0": "value0",
							"prop1": map[string]interface{}{
								"this": "that",
							},
						},
					},
				},
			},
		},
	}
	for i, c := range cases {
		var w Wrapper
		err := json.Unmarshal([]byte(c.s), &w)
		if err != nil {
			t.Errorf("error unmarshaling case %d: %v", i, err)
			continue
		}
		if !reflect.DeepEqual(w.Value, c.expected) {
			t.Errorf("case %d failed: expected %#v, got %#v", i, c.expected, w.Value)
		}
	}
}
