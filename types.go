package geojson

var _ Object = &FeatureCollection{}
var _ Object = &Feature{}
var _ Object = &GeometryCollection{}
var _ Object = &MultiPolygon{}
var _ Object = &Polygon{}
var _ Object = &MultiLineString{}
var _ Object = &LineString{}
var _ Object = &MultiPoint{}
var _ Object = &Point{}

var _ Geometry = &GeometryCollection{}
var _ Geometry = &MultiPolygon{}
var _ Geometry = &Polygon{}
var _ Geometry = &MultiLineString{}
var _ Geometry = &LineString{}
var _ Geometry = &MultiPoint{}
var _ Geometry = &Point{}

// A Wrapper wraps any GeoJSON object for JSON deserialization when the base
// type is not known in advance. If the type is known or if an object is being
// serialized, the bare type itself can be used.
type Wrapper struct {
	Value Object
}

// An Object is any GeoJSON-encoded value. Protocols are advised to wrap outer
// (transmitted) data in FeatureCollections, if doing so makes sense.
type Object interface {
	isObject()
}

// A Geometry is any of the following concrete object definitions:
// - GeometryCollection
// - MultiPolygon
// - Polygon
// - MultiLineString,
// - LineString
// - MultiPoint
// - Point
type Geometry interface {
	Object
	isGeometry()
}

// A FeatureCollection is a collection of features.
type FeatureCollection struct {
	// The features contained in this collection.
	Features []Feature
}

func (*FeatureCollection) isObject() {}

// A Feature is any spatially-bounded entity with optional metadata.
type Feature struct {
	// The geometric definition of this feature.
	Geometry Geometry
	// Properties associated with this feature. For example, this might include
	// a feature name along with other standard metadata.
	Properties map[string]interface{}
	// ID associated with this feature. Optional. For best compatibility, this
	// should go under properties.
	ID string
}

func (*Feature) isObject() {}

// A GeometryCollection is a collection of geometric structures. Note that a
// GeometryCollection may technically contain nested GeometryCollections, but
// that applications should avoid doing so.
type GeometryCollection struct {
	Geometries []Geometry
}

func (*GeometryCollection) isObject() {}

func (*GeometryCollection) isGeometry() {}

// A MultiPolygon is a collection of Polygons.
type MultiPolygon struct {
	// Polygons in this MultiPolygon.
	Polygons []Polygon
}

func (*MultiPolygon) isObject() {}

func (*MultiPolygon) isGeometry() {}

// A Polygon defines a shape in WGS84.
//
// The shape is defined by a sequence of linear rings. Each linear ring is a
// LineString whose first position is equivalent to its last (i.e., forms a
// closed area). Linear rings enclose the area to the left of the LineString as
// defined. This means that "clockwise" LineStrings (when viewed in a planar
// projection and not crossing the poles or anti-meridian) actually form a
// negative space and define the area outside of the enclosed area.
//
// The first linear ring in a Polygon must enclose its defined area (i.e., it
// should be counter-clockwise when not cutting the anti-meridian or enclosing
// the poles). The following linear rings that make up a Polygon define holes
// in its otherwise contiguous enclosure.
//
// This library does not make any attempt to be backward-compatible with
// GeoJSON 2008, which did not enforce handedness in enclosing linear rings.
type Polygon struct {
	// Linear rings that constitute this Polygon. Each LineString must consist
	// of at least 4 positions.
	Rings []LineString
}

func (*Polygon) isObject() {}

func (*Polygon) isGeometry() {}

// A MultiLineString is a collection of LineStrings. While this has the same
// structure as a Polygon, there is no semantic beyond that of the individual
// LineStrings it contains. Additionally, the contained LineStrings need not
// form linear rings.
type MultiLineString struct {
	// LineStrings within this MultiLineString.
	Lines []LineString
}

func (*MultiLineString) isObject() {}

func (*MultiLineString) isGeometry() {}

// A LineString is a sequenced collections of point positions forming a
// contiguous path.
type LineString struct {
	// Positions that make up this LineString. Must contain at least 2 points.
	Points []Point
}

func (*LineString) isObject() {}

func (*LineString) isGeometry() {}

// A MultiPoint is a collection of multiple point positions.
type MultiPoint struct {
	// Individual points that make up this MultiPoint.
	Points []Point
}

func (*MultiPoint) isObject() {}

func (*MultiPoint) isGeometry() {}

// A Point is the most basic type and captures a single, precise position in
// WGS84.
type Point struct {
	// X position (usually longitude, sometimes easting)
	X float64
	// Y position (usually latitude, sometimes northing)
	Y float64
	// Elevation in opaque units.
	Elevation float64
	// Whether the associated elevation is valid.
	HasElevation bool
}

func (*Point) isObject() {}

func (*Point) isGeometry() {}
