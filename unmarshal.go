package geojson

import (
	"encoding/json"
	"fmt"

	"github.com/bsidhom/geojson/wire"
)

var _ json.Unmarshaler = (*Wrapper)(nil)
var _ json.Unmarshaler = (*FeatureCollection)(nil)
var _ json.Unmarshaler = (*Feature)(nil)
var _ json.Unmarshaler = (*GeometryCollection)(nil)
var _ json.Unmarshaler = (*MultiPolygon)(nil)
var _ json.Unmarshaler = (*Polygon)(nil)
var _ json.Unmarshaler = (*MultiLineString)(nil)
var _ json.Unmarshaler = (*LineString)(nil)
var _ json.Unmarshaler = (*MultiPoint)(nil)
var _ json.Unmarshaler = (*Point)(nil)

func (w *Wrapper) UnmarshalJSON(b []byte) error {
	var wireWrapper wire.Wrapper
	err := json.Unmarshal(b, &wireWrapper)
	if err != nil {
		return err
	}
	switch t := wireWrapper.Value.(type) {
	case *wire.FeatureCollection:
		f := &FeatureCollection{}
		err := f.unmarshalFrom(t)
		if err != nil {
			return err
		}
		w.Value = f
	case *wire.Feature:
		f := &Feature{}
		err := f.unmarshalFrom(t)
		if err != nil {
			return err
		}
		w.Value = f
	case wire.Geometry:
		result, err := unmarshalGeometry(t)
		if err != nil {
			return err
		}
		w.Value = result
	default:
		return fmt.Errorf("invalid wire object type: %T", t)
	}
	return nil
}

func (f *FeatureCollection) UnmarshalJSON(b []byte) error {
	var w wire.FeatureCollection
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	return f.unmarshalFrom(&w)
}

func (f *FeatureCollection) unmarshalFrom(w *wire.FeatureCollection) error {
	*f = FeatureCollection{}
	if len(w.Features) == 0 {
		return nil
	}
	features := make([]Feature, len(w.Features))
	for i, wireFeature := range w.Features {
		err := features[i].unmarshalFrom(&wireFeature)
		if err != nil {
			return fmt.Errorf("unmarshal FeatureCollection: feature %d: %v", i, err)
		}
	}

	f.Features = features
	return nil
}

func (f *Feature) UnmarshalJSON(b []byte) error {
	var w wire.Feature
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	return f.unmarshalFrom(&w)
}

func (f *Feature) unmarshalFrom(w *wire.Feature) error {
	*f = Feature{}
	g, err := unmarshalGeometry(w.Geometry)
	if err != nil {
		return err
	}

	f.Geometry = g
	f.Properties = w.Properties
	f.ID = w.ID
	return nil
}

func (g *GeometryCollection) UnmarshalJSON(b []byte) error {
	var w wire.GeometryCollection
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	return g.unmarshalFrom(w.Geometries)
}

func (g *GeometryCollection) unmarshalFrom(geometries []wire.Geometry) error {
	*g = GeometryCollection{}
	if len(geometries) == 0 {
		return nil
	}
	gs := make([]Geometry, len(geometries))
	for i, geometry := range geometries {
		g, err := unmarshalGeometry(geometry)
		if err != nil {
			return fmt.Errorf("unmarshal GeometryCollection: geometry %d: %v", i, err)
		}
		gs[i] = g
	}

	g.Geometries = gs
	return nil
}

func (m *MultiPolygon) UnmarshalJSON(b []byte) error {
	var w wire.MultiPolygon
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	return m.unmarshalFrom(w.Coordinates)
}

func (m *MultiPolygon) unmarshalFrom(coords [][][][]float64) error {
	*m = MultiPolygon{}
	if len(coords) == 0 {
		return nil
	}
	polygons := make([]Polygon, len(coords))
	for i, polygonCoords := range coords {
		err := polygons[i].unmarshalFrom(polygonCoords)
		if err != nil {
			return fmt.Errorf("unmarshal MultiPolygon: %v", err)
		}
	}
	m.Polygons = polygons
	return nil
}

func (p *Polygon) UnmarshalJSON(b []byte) error {
	var w wire.Polygon
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	return p.unmarshalFrom(w.Coordinates)
}

func (p *Polygon) unmarshalFrom(coords [][][]float64) error {
	numRings := len(coords)
	if numRings < 1 {
		return fmt.Errorf("unmarshal Polygon: must have at least 1 linear ring")
	}
	rings := make([]LineString, numRings)
	for i, ringCoords := range coords {
		// Resist the urge to reach in and inspect coordinates first. Unmarshal
		// and then do verification on parsed lines.
		err := rings[i].unmarshalFrom(ringCoords)
		if err != nil {
			return fmt.Errorf("unmarshal Polygon: %v", err)
		}

		numPoints := len(rings[i].Points)
		if numPoints < 4 {
			return fmt.Errorf("unmarshal Polygon: each linear ring requires at least 4 points, ring %d has %d", i, numPoints)
		}

		firstPoint := rings[i].Points[0]
		lastPoint := rings[i].Points[numPoints-1]
		if firstPoint != lastPoint {
			return fmt.Errorf("unmarshal Polygon: linear ring %d first point (%v) does not match last point (%v)", i, firstPoint, lastPoint)
		}
	}

	p.Rings = rings
	return nil
}

func (m *MultiLineString) UnmarshalJSON(b []byte) error {
	var w wire.MultiLineString
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	return m.unmarshalFrom(w.Coordinates)
}

func (m *MultiLineString) unmarshalFrom(coords [][][]float64) error {
	*m = MultiLineString{}
	if len(coords) == 0 {
		return nil
	}
	lines := make([]LineString, len(coords))
	for i, lineCoords := range coords {
		err := lines[i].unmarshalFrom(lineCoords)
		if err != nil {
			return fmt.Errorf("unmarshal MultiLineString: %v", err)
		}
	}

	m.Lines = lines
	return nil
}

func (ls *LineString) UnmarshalJSON(b []byte) error {
	var w wire.LineString
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	return ls.unmarshalFrom(w.Coordinates)
}

func (ls *LineString) unmarshalFrom(coords [][]float64) error {
	*ls = LineString{}
	numPoints := len(coords)
	if numPoints < 2 {
		return fmt.Errorf("unmarshal LineString: must have at least 2 points, got %d", numPoints)
	}
	points := make([]Point, len(coords))
	for i, pointCoords := range coords {
		err := points[i].unmarshalFrom(pointCoords)
		if err != nil {
			return fmt.Errorf("unmarshal LineString: %v", err)
		}
	}
	ls.Points = points
	return nil
}

func (m *MultiPoint) UnmarshalJSON(b []byte) error {
	var w wire.MultiPoint
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	return m.unmarshalFrom(w.Coordinates)
}

func (m *MultiPoint) unmarshalFrom(coords [][]float64) error {
	*m = MultiPoint{}
	if len(coords) == 0 {
		return nil
	}
	points := make([]Point, len(coords))
	for i, pointCoords := range coords {
		err := points[i].unmarshalFrom(pointCoords)
		if err != nil {
			return fmt.Errorf("unmarshal MultiPoint: %v", err)
		}
	}

	m.Points = points
	return nil
}

func (p *Point) UnmarshalJSON(b []byte) error {
	var w wire.Point
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	return p.unmarshalFrom(w.Coordinates)
}

func (p *Point) unmarshalFrom(coords []float64) error {
	// Reset Point.
	*p = Point{}
	numCoords := len(coords)
	if numCoords < 2 || numCoords > 3 {
		return fmt.Errorf("unmarshal Point: must have 2-3 coordinates, got %d", numCoords)
	}

	p.X = coords[0]
	p.Y = coords[1]
	if numCoords == 3 {
		p.Elevation = coords[2]
		p.HasElevation = true
	}
	return nil
}

func unmarshalGeometry(g wire.Geometry) (Geometry, error) {
	var result Geometry
	switch t := g.(type) {
	case *wire.GeometryCollection:
		g := &GeometryCollection{}
		// NOTE: This ultimately calls unmarshalGeometry recursively. The
		// implementation should tolerate this anyway because
		// GeometryCollection objects are allowed to contain other
		// GeometryCollections per the spec, even though this is advised
		// against.
		err := g.unmarshalFrom(t.Geometries)
		if err != nil {
			return nil, err
		}
		result = g
	case *wire.MultiPolygon:
		m := &MultiPolygon{}
		err := m.unmarshalFrom(t.Coordinates)
		if err != nil {
			return nil, err
		}
		result = m
	case *wire.Polygon:
		p := &Polygon{}
		err := p.unmarshalFrom(t.Coordinates)
		if err != nil {
			return nil, err
		}
		result = p
	case *wire.MultiLineString:
		m := &MultiLineString{}
		err := m.unmarshalFrom(t.Coordinates)
		if err != nil {
			return nil, err
		}
		result = m
	case *wire.LineString:
		ls := &LineString{}
		err := ls.unmarshalFrom(t.Coordinates)
		if err != nil {
			return nil, err
		}
		result = ls
	case *wire.MultiPoint:
		m := &MultiPoint{}
		err := m.unmarshalFrom(t.Coordinates)
		if err != nil {
			return nil, err
		}
		result = m
	case *wire.Point:
		p := &Point{}
		err := p.unmarshalFrom(t.Coordinates)
		if err != nil {
			return nil, err
		}
		result = p
	default:
		return nil, fmt.Errorf("invalid wire geometry type: %T", t)
	}
	return result, nil
}
