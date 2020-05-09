package wire

import (
	"encoding/json"
	"fmt"
)

var _ json.Unmarshaler = (*Wrapper)(nil)
var _ json.Unmarshaler = (*Feature)(nil)
var _ json.Unmarshaler = (*GeometryCollection)(nil)

func (obj *Wrapper) UnmarshalJSON(b []byte) error {
	type object Wrapper
	o := &object{}
	err := json.Unmarshal(b, o)
	if err != nil {
		return err
	}
	factory, ok := typeNames[o.Type]
	if !ok {
		return fmt.Errorf("invalid type name: %q", o.Type)
	}
	t := factory()
	err = json.Unmarshal(b, t)
	if err != nil {
		return err
	}
	o.Value = t
	*obj = Wrapper(*o)
	return nil
}

func (f *Feature) UnmarshalJSON(b []byte) error {
	type WireType struct {
		BBox       []float64              `json:"bbox"`
		Geometry   json.RawMessage        `json:"geometry"`
		Properties map[string]interface{} `json:"properties"`
		ID         string                 `json:"id"`
	}
	var w WireType
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}
	var obj Wrapper
	err = json.Unmarshal(w.Geometry, &obj)
	if err != nil {
		return err
	}

	geometry, ok := obj.Value.(Geometry)
	if !ok {
		return fmt.Errorf("invalid non-geometry type: %T", geometry)
	}

	f.BBox = w.BBox
	f.Geometry = geometry
	f.Properties = w.Properties
	f.ID = w.ID
	return nil
}

func (g *GeometryCollection) UnmarshalJSON(b []byte) error {
	type WireType struct {
		BBox       []float64         `json:"bbox"`
		Geometries []json.RawMessage `json:"geometries"`
	}
	var w WireType
	err := json.Unmarshal(b, &w)
	if err != nil {
		return err
	}

	geometries := make([]Geometry, len(w.Geometries))
	for i, wireGeometry := range w.Geometries {
		var obj Wrapper
		err := json.Unmarshal(wireGeometry, &obj)
		if err != nil {
			return err
		}
		geometry, ok := obj.Value.(Geometry)
		if !ok {
			return fmt.Errorf("invalid non-geometry type: %T", obj.Value)
		}
		geometries[i] = geometry
	}

	g.BBox = w.BBox
	g.Geometries = geometries
	return nil
}
