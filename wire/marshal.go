package wire

import "encoding/json"

var _ json.Marshaler = (*FeatureCollection)(nil)
var _ json.Marshaler = (*Feature)(nil)
var _ json.Marshaler = (*GeometryCollection)(nil)
var _ json.Marshaler = (*MultiPolygon)(nil)
var _ json.Marshaler = (*Polygon)(nil)
var _ json.Marshaler = (*MultiLineString)(nil)
var _ json.Marshaler = (*LineString)(nil)
var _ json.Marshaler = (*MultiPoint)(nil)
var _ json.Marshaler = (*Point)(nil)

func (f *FeatureCollection) MarshalJSON() ([]byte, error) {
	type WireType FeatureCollection
	type t struct {
		Type string `json:"type"`
		*WireType
	}
	v := t{
		Type:     featureCollectionType,
		WireType: (*WireType)(f),
	}
	return json.Marshal(v)
}

func (f *Feature) MarshalJSON() ([]byte, error) {
	type WireType Feature
	type t struct {
		Type string `json:"type"`
		*WireType
	}
	v := t{
		Type:     featureType,
		WireType: (*WireType)(f),
	}
	return json.Marshal(v)
}

// All geometry types use the same logic for marshaling. Unfortunately, without
// generics or resorting to reflection, we have to spell this out for each type.

func (g *GeometryCollection) MarshalJSON() ([]byte, error) {
	type WireType GeometryCollection
	type t struct {
		Type string `json:"type"`
		*WireType
	}
	v := t{
		Type:     geometryCollectionType,
		WireType: (*WireType)(g),
	}
	return json.Marshal(v)
}

func (m *MultiPolygon) MarshalJSON() ([]byte, error) {
	type WireType MultiPolygon
	type t struct {
		Type string `json:"type"`
		*WireType
	}
	v := t{
		Type:     multiPolygonType,
		WireType: (*WireType)(m),
	}
	return json.Marshal(v)
}

func (p *Polygon) MarshalJSON() ([]byte, error) {
	type WireType Polygon
	type t struct {
		Type string `json:"type"`
		*WireType
	}
	v := t{
		Type:     polygonType,
		WireType: (*WireType)(p),
	}
	return json.Marshal(v)
}

func (m *MultiLineString) MarshalJSON() ([]byte, error) {
	type WireType MultiLineString
	type t struct {
		Type string `json:"type"`
		*WireType
	}
	v := t{
		Type:     multiLineStringType,
		WireType: (*WireType)(m),
	}
	return json.Marshal(v)
}

func (ls *LineString) MarshalJSON() ([]byte, error) {
	type WireType MultiPoint
	type t struct {
		Type string `json:"type"`
		*WireType
	}
	v := t{
		Type:     lineStringType,
		WireType: (*WireType)(ls),
	}
	return json.Marshal(v)
}

func (m *MultiPoint) MarshalJSON() ([]byte, error) {
	type WireType MultiPoint
	type t struct {
		Type string `json:"type"`
		*WireType
	}
	v := t{
		Type:     multiPointType,
		WireType: (*WireType)(m),
	}
	return json.Marshal(v)
}

func (p *Point) MarshalJSON() ([]byte, error) {
	type WireType Point
	type t struct {
		Type string `json:"type"`
		*WireType
	}
	v := t{
		Type:     pointType,
		WireType: (*WireType)(p),
	}
	return json.Marshal(v)
}
