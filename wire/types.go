package wire

import (
	"encoding/json"
	"fmt"
)

const (
	featureCollectionType  = "FeatureCollection"
	featureType            = "Feature"
	geometryCollectionType = "GeometryCollection"
	multiPolygonType       = "MultiPolygon"
	polygonType            = "Polygon"
	multiLineStringType    = "MultiLineString"
	lineStringType         = "LineString"
	multiPointType         = "MultiPoint"
	pointType              = "Point"
)

var typeNames = map[string]func() Object{
	featureCollectionType:  func() Object { return &FeatureCollection{} },
	featureType:            func() Object { return &Feature{} },
	geometryCollectionType: func() Object { return &GeometryCollection{} },
	multiPolygonType:       func() Object { return &MultiPolygon{} },
	polygonType:            func() Object { return &Polygon{} },
	multiLineStringType:    func() Object { return &MultiLineString{} },
	lineStringType:         func() Object { return &LineString{} },
	multiPointType:         func() Object { return &MultiPoint{} },
	pointType:              func() Object { return &Point{} },
}

// A Wrapper holds a deserialized GeoJSON value. This special type allows for
// GeoJSON objects to be deserialized into their correct types. Use this to
// deserialize values of unknown GeoJSON type. If the type is known ahead of
// time or values are being serialized, the bare types can be used.
//
// Note that if a bare type is used for deserialization, the json "type" field
// will not be verified.
type Wrapper struct {
	Type  string `json:"type"`
	Value Object `json:"-"`
}

// An Object is any GeoJSON object type. This includes FeatureCollection,
// Feature, and all Geometry types.
type Object interface {
	isObject()
}

// A Geometry represents any GeoJSON geometry type: GeometryCollection,
// MultiPolygon, Polygon, MultiLineString, LineString, MultiPoint, or Point.
type Geometry interface {
	// Geometry types must implement json.Marshaler to guarantee their custom
	// marshalers are used.
	json.Marshaler
	isGeometry()
}

type FeatureCollection struct {
	BBox     []float64 `json:"bbox,omitempty"`
	Features []Feature `json:"features"`
}

func (f *FeatureCollection) isObject() {}

func (f *FeatureCollection) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *FeatureCollection) GoString() string {
	return fmt.Sprintf("%#v", *f)
}

type Feature struct {
	BBox       []float64              `json:"bbox,omitempty"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
	ID         string                 `json:"id,omitempty"`
}

func (f *Feature) isObject() {}

func (f *Feature) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *Feature) GoString() string {
	return fmt.Sprintf("%#v", *f)
}

type GeometryCollection struct {
	BBox       []float64  `json:"bbox,omitempty"`
	Geometries []Geometry `json:"geometries"`
}

func (g *GeometryCollection) isObject() {}

func (g *GeometryCollection) isGeometry() {}

func (g *GeometryCollection) String() string {
	return fmt.Sprintf("%v", *g)
}

func (g *GeometryCollection) GoString() string {
	return fmt.Sprintf("%#v", *g)
}

type MultiPolygon struct {
	BBox        []float64       `json:"bbox,omitempty"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

func (m *MultiPolygon) isObject() {}

func (m *MultiPolygon) isGeometry() {}

func (m *MultiPolygon) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MultiPolygon) GoString() string {
	return fmt.Sprintf("%#v", *m)
}

type Polygon struct {
	BBox        []float64     `json:"bbox,omitempty"`
	Coordinates [][][]float64 `json:"coordinates"`
}

func (p *Polygon) isObject() {}

func (p *Polygon) isGeometry() {}

func (p *Polygon) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *Polygon) GoString() string {
	return fmt.Sprintf("%#v", *p)
}

type MultiLineString struct {
	BBox        []float64     `json:"bbox,omitempty"`
	Coordinates [][][]float64 `json:"coordinates"`
}

func (m *MultiLineString) isObject() {}

func (m *MultiLineString) isGeometry() {}

func (m *MultiLineString) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MultiLineString) GoString() string {
	return fmt.Sprintf("%#v", *m)
}

type LineString struct {
	BBox        []float64   `json:"bbox,omitempty"`
	Coordinates [][]float64 `json:"coordinates"`
}

func (ls *LineString) isObject() {}

func (ls *LineString) isGeometry() {}

func (ls *LineString) String() string {
	return fmt.Sprintf("%v", *ls)
}

func (ls *LineString) GoString() string {
	return fmt.Sprintf("%#v", *ls)
}

type MultiPoint struct {
	BBox        []float64   `json:"bbox,omitempty"`
	Coordinates [][]float64 `json:"coordinates"`
}

func (m *MultiPoint) isObject() {}

func (m *MultiPoint) isGeometry() {}

func (m *MultiPoint) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MultiPoint) GoString() string {
	return fmt.Sprintf("%#v", *m)
}

type Point struct {
	BBox        []float64 `json:"bbox,omitempty"`
	Coordinates []float64 `json:"coordinates"`
}

func (p *Point) isObject() {}

func (p *Point) isGeometry() {}

func (p *Point) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *Point) GoString() string {
	return fmt.Sprintf("%#v", *p)
}
