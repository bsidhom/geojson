package wire

import (
	"encoding/json"
	"reflect"
	"testing"
)

// The internal json encoding may not be stable. Instead of testing string
// representation, verify that round trips preserve information. The main
// downside is that we cannot test for omitted fields in JSON this way.
func TestRoundTrip_MarshalJSON(t *testing.T) {
	cases := []Object{
		&Point{Coordinates: []float64{0, 1}},
		&Point{Coordinates: []float64{3, 1}, BBox: []float64{0, 10, 0, 10}},
		&MultiPoint{
			Coordinates: [][]float64{
				{0, 1},
				{1, 0},
			},
		},
		&MultiPoint{
			BBox: []float64{0, 0, 1, 1},
			Coordinates: [][]float64{
				{0, 1},
				{1, 0},
			},
		},
		&LineString{
			Coordinates: [][]float64{
				{0, 1},
				{33, 1},
			},
		},
		&MultiLineString{
			Coordinates: [][][]float64{
				{
					{0, 1},
					{1, 2},
				},
				{
					{3, 5},
					{4, 6},
					{4, 7},
				},
			},
		},
		&Polygon{
			Coordinates: [][][]float64{
				{
					{1, 0},
					{0, 1},
					{-1, 0},
					{1, 0},
				},
			},
		},
		&MultiPolygon{
			Coordinates: [][][][]float64{
				{
					{
						{1, 0},
						{0, 1},
						{-1, 0},
						{1, 0},
					},
				},
			},
		},
		&GeometryCollection{
			Geometries: []Geometry{
				&Point{Coordinates: []float64{3, 4}},
				&Polygon{
					Coordinates: [][][]float64{
						{
							{1, 0},
							{0, 1},
							{-1, 0},
							{1, 0},
						},
					},
				},
			},
		},
		&Feature{
			Geometry:   &Point{Coordinates: []float64{4, 5}},
			Properties: map[string]interface{}{"name": "My Point"},
		},
		&FeatureCollection{
			Features: []Feature{
				{
					Geometry: &MultiPoint{
						Coordinates: [][]float64{
							{0, 5},
							{6, 7},
						},
					},
					Properties: map[string]interface{}{"number": float64(3)},
				},
			},
		},
	}

	for i, c := range cases {
		b, err := json.Marshal(c)
		if err != nil {
			t.Errorf("failed to serialize case %d (%T): %v", i, c, err)
			continue
		}
		var obj Wrapper
		err = json.Unmarshal(b, &obj)
		if err != nil {
			t.Errorf("failed to deserialize case %d (%T): %v", i, c, err)
			continue
		}
		v := obj.Value
		if !reflect.DeepEqual(c, v) {
			t.Errorf("round trip %d (%T) failed: expected %#v, got %#v", i, c, c, v)
		}
	}
}

// TODO: Marshal then unmarshal into map[string]interface{} to verify that
// omitted fields are not serialized.
