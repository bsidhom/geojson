package wire

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Test that unmarshal works on explicit JSON representations.
// NOTE: We use reflect.DeepEqual for comparisons. Types stored in interfaces
// must have matching types as the deserialized JSON values. In particular,
// numbers must be represented as float64 values.
func TestObject_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		s        string
		expected Object
	}{
		{
			s:        `{"type":"Point","coordinates":[0,1]}`,
			expected: &Point{Coordinates: []float64{0, 1}},
		},
		{
			s: `{"type":"MultiPoint","coordinates":[[100,0],[101,1]]}`,
			expected: &MultiPoint{
				Coordinates: [][]float64{
					{100, 0},
					{101, 1},
				},
			},
		},
		{
			s: `{"type":"LineString","coordinates":[[100,0],[101,1]]}`,
			expected: &LineString{
				Coordinates: [][]float64{
					{100, 0},
					{101, 1},
				},
			},
		},
		{
			s: `{"type":"MultiLineString","coordinates":[[[100,0],[101,1]],[[102,2],[103,3]]]}`,
			expected: &MultiLineString{
				Coordinates: [][][]float64{
					{
						{100, 0},
						{101, 1},
					},
					{
						{102, 2},
						{103, 3},
					},
				},
			},
		},
		{
			s: `{"type":"Polygon","coordinates":[[[100,0],[101,0],[101,1],[100,1],[100,0]],[[100.8,0.8],[100.8,0.2],[100.2,0.2],[100.2,0.8],[100.8,0.8]]]}`,
			expected: &Polygon{
				Coordinates: [][][]float64{
					{
						{100, 0},
						{101, 0},
						{101, 1},
						{100, 1},
						{100, 0},
					},
					{
						{100.8, 0.8},
						{100.8, 0.2},
						{100.2, 0.2},
						{100.2, 0.8},
						{100.8, 0.8},
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
				Coordinates: [][][][]float64{
					{
						{
							{102, 2},
							{103, 2},
							{103, 3},
							{102, 3},
							{102, 2},
						},
					},
					{
						{
							{100, 0},
							{101, 0},
							{101, 1},
							{100, 1},
							{100, 0},
						},
						{
							{100.2, 0.2},
							{100.2, 0.8},
							{100.8, 0.8},
							{100.8, 0.2},
							{100.2, 0.2},
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
					&Point{Coordinates: []float64{100, 0}},
					&LineString{
						Coordinates: [][]float64{
							{101, 0},
							{102, 1},
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
  }
}`,
			expected: &Feature{
				BBox: []float64{-10, -10, 10, 10},
				Geometry: &Polygon{
					Coordinates: [][][]float64{
						{
							{-10, -10},
							{10, -10},
							{10, 10},
							{-10, -10},
						},
					},
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
						Geometry: &Point{Coordinates: []float64{102, 0.5}},
						Properties: map[string]interface{}{
							"prop0": "value0",
						},
					},
					{
						Geometry: &LineString{
							Coordinates: [][]float64{
								{102, 0},
								{103, 1},
								{104, 0},
								{105, 1},
							},
						},
						Properties: map[string]interface{}{
							"prop0": "value0",
							"prop1": float64(0),
						},
					},
					{
						Geometry: &Polygon{
							Coordinates: [][][]float64{
								{
									{100, 0},
									{101, 0},
									{101, 1},
									{100, 1},
									{100, 0},
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
		var obj Wrapper
		err := json.Unmarshal([]byte(c.s), &obj)
		if err != nil {
			t.Errorf("error unmarshaling case %d: %v", i, err)
			continue
		}
		value := obj.Value
		if !reflect.DeepEqual(value, c.expected) {
			t.Errorf("for case %d expected %#v, got %#v", i, c.expected, value)
		}
	}
}
